package dnsx

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"

	"github.com/miekg/dns"
	"golang.org/x/net/idna"
)

// DNSAnswer represents a single DNS record answer.
type DNSAnswer struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	TTL      uint32 `json:"ttl"`
	Priority uint16 `json:"priority,omitempty"`
}

// ResolverResult holds responses from a specific resolver.
type ResolverResult struct {
	NodeId    string      `json:"node_id,omitempty"`
	NodeName  string      `json:"node_name,omitempty"`
	ISP       string      `json:"isp,omitempty"`
	Resolver  string      `json:"resolver"`
	LatencyMs int64       `json:"latency_ms"`
	Answers   []DNSAnswer `json:"answers"`
	Error     string      `json:"error,omitempty"`
}

// DNSQueryResult represents the full response from multi-resolver query.
type DNSQueryResult struct {
	Name    string           `json:"name"`
	Type    string           `json:"type"`
	Results []ResolverResult `json:"results"`
}

// DNSService performs multi-resolver DNS lookups.
type DNSService struct {
	resolvers []string
}

// NewDNSService creates a DNSService.
func NewDNSService(resolvers []string) *DNSService {
	if len(resolvers) == 0 {
		resolvers = []string{"system", "1.1.1.1", "8.8.8.8", "223.5.5.5", "119.29.29.29"}
	}
	return &DNSService{
		resolvers: resolvers,
	}
}

// SupportedRecordTypes maps string record names to miekg/dns uint16 types.
var SupportedRecordTypes = map[string]uint16{
	"A":     dns.TypeA,
	"AAAA":  dns.TypeAAAA,
	"CNAME": dns.TypeCNAME,
	"MX":    dns.TypeMX,
	"TXT":   dns.TypeTXT,
	"NS":    dns.TypeNS,
	"CAA":   dns.TypeCAA,
	"SRV":   dns.TypeSRV,
	"PTR":   dns.TypePTR,
	"SOA":   dns.TypeSOA,
}

// Query performs concurrent DNS lookups across all configured resolvers.
func (s *DNSService) Query(ctx context.Context, domain string, recordType string) (*DNSQueryResult, error) {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return nil, fmt.Errorf("domain name cannot be empty")
	}

	// IDN Punycode conversion
	punyDomain, err := idna.ToASCII(domain)
	if err == nil && punyDomain != "" {
		domain = punyDomain
	}

	if err := ssrf.ValidateTargetHost(domain); err != nil {
		return nil, err
	}

	recordType = strings.ToUpper(strings.TrimSpace(recordType))
	qType, ok := SupportedRecordTypes[recordType]
	if !ok {
		return nil, fmt.Errorf("unsupported DNS record type: %s", recordType)
	}

	// Ensure FQDN trailing dot
	fqdn := dns.Fqdn(domain)

	// Enforce overall timeout
	queryCtx, cancel := context.WithTimeout(ctx, config.DefaultDNSOverallTimeout)
	defer cancel()

	var wg sync.WaitGroup
	results := make([]ResolverResult, len(s.resolvers))

	for i, r := range s.resolvers {
		wg.Add(1)
		go func(idx int, resolver string) {
			defer wg.Done()
			results[idx] = s.queryResolver(queryCtx, fqdn, qType, resolver)
		}(i, r)
	}

	wg.Wait()

	return &DNSQueryResult{
		Name:    strings.TrimSuffix(fqdn, "."),
		Type:    recordType,
		Results: results,
	}, nil
}

func (s *DNSService) queryResolver(ctx context.Context, fqdn string, qType uint16, resolver string) ResolverResult {
	start := time.Now()
	res := ResolverResult{
		Resolver: resolver,
		Answers:  []DNSAnswer{},
	}

	if resolver == "system" {
		return s.querySystem(ctx, fqdn, qType, start)
	}

	// Address for miekg/dns client
	resolverAddr := resolver
	if !strings.Contains(resolverAddr, ":") {
		resolverAddr = net.JoinHostPort(resolverAddr, "53")
	}

	client := &dns.Client{
		Net:     "udp",
		Timeout: config.DefaultDNSSingleResolverTimeout,
	}

	msg := new(dns.Msg)
	msg.SetQuestion(fqdn, qType)
	msg.RecursionDesired = true

	in, rtt, err := client.ExchangeContext(ctx, msg, resolverAddr)
	if err != nil {
		res.LatencyMs = time.Since(start).Milliseconds()
		res.Error = err.Error()
		return res
	}

	res.LatencyMs = rtt.Milliseconds()
	if in.Rcode != dns.RcodeSuccess && in.Rcode != dns.RcodeNameError {
		res.Error = dns.RcodeToString[in.Rcode]
		return res
	}

	for _, rr := range in.Answer {
		if ans := parseRR(rr); ans != nil {
			res.Answers = append(res.Answers, *ans)
		}
	}

	return res
}

func (s *DNSService) querySystem(ctx context.Context, fqdn string, qType uint16, start time.Time) ResolverResult {
	cleanName := strings.TrimSuffix(fqdn, ".")
	res := ResolverResult{
		Resolver: "System",
		Answers:  []DNSAnswer{},
	}

	var r net.Resolver
	switch qType {
	case dns.TypeA:
		ips, err := r.LookupIP(ctx, "ip4", cleanName)
		res.LatencyMs = time.Since(start).Milliseconds()
		if err != nil {
			res.Error = err.Error()
			return res
		}
		for _, ip := range ips {
			res.Answers = append(res.Answers, DNSAnswer{
				Type:  "A",
				Value: ip.String(),
				TTL:   300,
			})
		}
	case dns.TypeAAAA:
		ips, err := r.LookupIP(ctx, "ip6", cleanName)
		res.LatencyMs = time.Since(start).Milliseconds()
		if err != nil {
			res.Error = err.Error()
			return res
		}
		for _, ip := range ips {
			res.Answers = append(res.Answers, DNSAnswer{
				Type:  "AAAA",
				Value: ip.String(),
				TTL:   300,
			})
		}
	case dns.TypeCNAME:
		cname, err := r.LookupCNAME(ctx, cleanName)
		res.LatencyMs = time.Since(start).Milliseconds()
		if err != nil {
			res.Error = err.Error()
			return res
		}
		res.Answers = append(res.Answers, DNSAnswer{
			Type:  "CNAME",
			Value: strings.TrimSuffix(cname, "."),
			TTL:   300,
		})
	case dns.TypeMX:
		mxs, err := r.LookupMX(ctx, cleanName)
		res.LatencyMs = time.Since(start).Milliseconds()
		if err != nil {
			res.Error = err.Error()
			return res
		}
		for _, mx := range mxs {
			res.Answers = append(res.Answers, DNSAnswer{
				Type:     "MX",
				Value:    strings.TrimSuffix(mx.Host, "."),
				Priority: mx.Pref,
				TTL:      300,
			})
		}
	case dns.TypeTXT:
		txts, err := r.LookupTXT(ctx, cleanName)
		res.LatencyMs = time.Since(start).Milliseconds()
		if err != nil {
			res.Error = err.Error()
			return res
		}
		for _, txt := range txts {
			res.Answers = append(res.Answers, DNSAnswer{
				Type:  "TXT",
				Value: txt,
				TTL:   300,
			})
		}
	case dns.TypeNS:
		nss, err := r.LookupNS(ctx, cleanName)
		res.LatencyMs = time.Since(start).Milliseconds()
		if err != nil {
			res.Error = err.Error()
			return res
		}
		for _, ns := range nss {
			res.Answers = append(res.Answers, DNSAnswer{
				Type:  "NS",
				Value: strings.TrimSuffix(ns.Host, "."),
				TTL:   300,
			})
		}
	default:
		// Fallback to 1.1.1.1 for complex types if system resolver has no direct Go stdlib method
		return s.queryResolver(ctx, fqdn, qType, "1.1.1.1")
	}

	return res
}

func parseRR(rr dns.RR) *DNSAnswer {
	if rr == nil {
		return nil
	}

	ans := &DNSAnswer{
		TTL: rr.Header().Ttl,
	}

	switch v := rr.(type) {
	case *dns.A:
		ans.Type = "A"
		ans.Value = v.A.String()
	case *dns.AAAA:
		ans.Type = "AAAA"
		ans.Value = v.AAAA.String()
	case *dns.CNAME:
		ans.Type = "CNAME"
		ans.Value = strings.TrimSuffix(v.Target, ".")
	case *dns.MX:
		ans.Type = "MX"
		ans.Value = strings.TrimSuffix(v.Mx, ".")
		ans.Priority = v.Preference
	case *dns.TXT:
		ans.Type = "TXT"
		ans.Value = strings.Join(v.Txt, " ")
	case *dns.NS:
		ans.Type = "NS"
		ans.Value = strings.TrimSuffix(v.Ns, ".")
	case *dns.CAA:
		ans.Type = "CAA"
		ans.Value = fmt.Sprintf("%d %s %q", v.Flag, v.Tag, v.Value)
	case *dns.SRV:
		ans.Type = "SRV"
		ans.Value = fmt.Sprintf("%d %d %d %s", v.Priority, v.Weight, v.Port, strings.TrimSuffix(v.Target, "."))
		ans.Priority = v.Priority
	case *dns.PTR:
		ans.Type = "PTR"
		ans.Value = strings.TrimSuffix(v.Ptr, ".")
	case *dns.SOA:
		ans.Type = "SOA"
		ans.Value = fmt.Sprintf("%s %s %d %d %d %d %d",
			strings.TrimSuffix(v.Ns, "."),
			strings.TrimSuffix(v.Mbox, "."),
			v.Serial, v.Refresh, v.Retry, v.Expire, v.Minttl)
	default:
		ans.Type = dns.TypeToString[rr.Header().Rrtype]
		ans.Value = rr.String()
	}

	return ans
}
