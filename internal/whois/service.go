package whois

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"
)

// WHOISResult represents domain or IP registration info.
type WHOISResult struct {
	Query       string   `json:"query"`
	Type        string   `json:"type"` // "domain" or "ip"
	Domain      string   `json:"domain,omitempty"`
	Registrar   string   `json:"registrar,omitempty"`
	Created     string   `json:"created,omitempty"`
	Updated     string   `json:"updated,omitempty"`
	Expires     string   `json:"expires,omitempty"`
	Status      []string `json:"status,omitempty"`
	NameServers []string `json:"name_servers,omitempty"`
	DNSSEC      string   `json:"dnssec,omitempty"`

	// IP specific
	Network      string `json:"network,omitempty"`
	CIDR         string `json:"cidr,omitempty"`
	Organization string `json:"organization,omitempty"`
	Country      string `json:"country,omitempty"`
	ASN          string `json:"asn,omitempty"`

	RawText string `json:"raw_text,omitempty"`
	Source  string `json:"source"`
}

// WHOISService performs RDAP and WHOIS lookups.
type WHOISService struct {
	httpClient *http.Client
}

// NewWHOISService creates a new WHOISService.
func NewWHOISService() *WHOISService {
	return &WHOISService{
		httpClient: ssrf.NewSafeHTTPClient(config.DefaultWHOISTimeout),
	}
}

// Query looks up WHOIS/RDAP information for a domain or IP.
func (s *WHOISService) Query(ctx context.Context, target string) (*WHOISResult, error) {
	target = strings.TrimSpace(target)
	if strings.Contains(target, "://") {
		parts := strings.Split(target, "://")
		if len(parts) > 1 {
			target = parts[1]
		}
	}
	if slashIdx := strings.Index(target, "/"); slashIdx != -1 {
		target = target[:slashIdx]
	}
	if h, _, err := net.SplitHostPort(target); err == nil {
		target = h
	}

	if err := ssrf.ValidateTargetHost(target); err != nil {
		return nil, err
	}

	// Case 1: IP address
	if addr, err := netip.ParseAddr(target); err == nil {
		return s.queryIP(ctx, addr)
	}

	// Case 2: Domain name
	return s.queryDomain(ctx, target)
}

func (s *WHOISService) queryDomain(ctx context.Context, domain string) (*WHOISResult, error) {
	// Try RDAP via rdap.org first
	rdapURL := fmt.Sprintf("https://rdap.org/domain/%s", domain)
	req, err := http.NewRequestWithContext(ctx, "GET", rdapURL, nil)
	if err == nil {
		req.Header.Set("Accept", "application/rdap+json, application/json")
		resp, err := s.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var rdapData rdapDomainResponse
			if err := json.NewDecoder(resp.Body).Decode(&rdapData); err == nil {
				return parseRDAPDomain(domain, rdapData), nil
			}
		}
	}

	// Fallback to traditional port 43 WHOIS
	return s.queryWHOISTCP(ctx, domain, "whois.iana.org")
}

func (s *WHOISService) queryIP(ctx context.Context, addr netip.Addr) (*WHOISResult, error) {
	rdapURL := fmt.Sprintf("https://rdap.org/ip/%s", addr.String())
	req, err := http.NewRequestWithContext(ctx, "GET", rdapURL, nil)
	if err == nil {
		req.Header.Set("Accept", "application/rdap+json, application/json")
		resp, err := s.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var rdapData rdapIPResponse
			if err := json.NewDecoder(resp.Body).Decode(&rdapData); err == nil {
				return parseRDAPIP(addr.String(), rdapData), nil
			}
		}
	}

	// Fallback to traditional port 43 WHOIS
	return s.queryWHOISTCP(ctx, addr.String(), "whois.arin.net")
}

func (s *WHOISService) queryWHOISTCP(ctx context.Context, query, server string) (*WHOISResult, error) {
	var dialer net.Dialer
	dialer.Timeout = config.DefaultWHOISTimeout

	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(server, "43"))
	if err != nil {
		return &WHOISResult{
			Query:   query,
			Type:    "domain",
			RawText: fmt.Sprintf("WHOIS query error: %v", err),
			Source:  "error",
		}, nil
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(config.DefaultWHOISTimeout))
	fmt.Fprintf(conn, "%s\r\n", query)

	var rawLines []string
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		rawLines = append(rawLines, scanner.Text())
	}

	raw := strings.Join(rawLines, "\n")
	res := parseRawWHOIS(query, raw)
	res.Source = "whois-port43"
	return res, nil
}

type rdapDomainResponse struct {
	LdhName     string `json:"ldhName"`
	Handle      string `json:"handle"`
	Status      []string `json:"status"`
	Entities    []struct {
		VcardArray []interface{} `json:"vcardArray"`
		Roles      []string      `json:"roles"`
	} `json:"entities"`
	Events []struct {
		EventAction string `json:"eventAction"`
		EventDate   string `json:"eventDate"`
	} `json:"events"`
	Nameservers []struct {
		LdhName string `json:"ldhName"`
	} `json:"nameservers"`
	SecureDNS struct {
		DelegationSigned bool `json:"delegationSigned"`
	} `json:"secureDNS"`
}

type rdapIPResponse struct {
	Handle       string `json:"handle"`
	StartAddress string `json:"startAddress"`
	EndAddress   string `json:"endAddress"`
	Country      string `json:"country"`
	Name         string `json:"name"`
	Entities     []struct {
		VcardArray []interface{} `json:"vcardArray"`
		Roles      []string      `json:"roles"`
	} `json:"entities"`
}

func parseRDAPDomain(domain string, data rdapDomainResponse) *WHOISResult {
	res := &WHOISResult{
		Query:   domain,
		Type:    "domain",
		Domain:  data.LdhName,
		Status:  data.Status,
		Source:  "rdap",
	}

	if data.SecureDNS.DelegationSigned {
		res.DNSSEC = "signedDelegation"
	} else {
		res.DNSSEC = "unsigned"
	}

	for _, ns := range data.Nameservers {
		if ns.LdhName != "" {
			res.NameServers = append(res.NameServers, ns.LdhName)
		}
	}

	for _, ev := range data.Events {
		switch ev.EventAction {
		case "registration":
			res.Created = ev.EventDate
		case "last changed", "last update":
			res.Updated = ev.EventDate
		case "expiration":
			res.Expires = ev.EventDate
		}
	}

	for _, entity := range data.Entities {
		for _, role := range entity.Roles {
			if role == "registrar" {
				if len(entity.VcardArray) > 1 {
					if cards, ok := entity.VcardArray[1].([]interface{}); ok {
						for _, item := range cards {
							if arr, ok := item.([]interface{}); ok && len(arr) > 3 {
								if arr[0] == "fn" {
									res.Registrar = fmt.Sprintf("%v", arr[3])
								}
							}
						}
					}
				}
			}
		}
	}

	return res
}

func parseRDAPIP(ipStr string, data rdapIPResponse) *WHOISResult {
	return &WHOISResult{
		Query:        ipStr,
		Type:         "ip",
		Network:      fmt.Sprintf("%s - %s", data.StartAddress, data.EndAddress),
		Organization: data.Name,
		Country:      data.Country,
		Source:       "rdap",
	}
}

func parseRawWHOIS(query, raw string) *WHOISResult {
	res := &WHOISResult{
		Query:   query,
		Type:    "domain",
		RawText: raw,
	}

	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToLower(line), "registrar:") {
			res.Registrar = strings.TrimSpace(line[len("registrar:"):])
		} else if strings.HasPrefix(strings.ToLower(line), "creation date:") {
			res.Created = strings.TrimSpace(line[len("creation date:"):])
		} else if strings.HasPrefix(strings.ToLower(line), "updated date:") {
			res.Updated = strings.TrimSpace(line[len("updated date:"):])
		} else if strings.HasPrefix(strings.ToLower(line), "registry expiry date:") || strings.HasPrefix(strings.ToLower(line), "exp date:") {
			res.Expires = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.HasPrefix(strings.ToLower(line), "name server:") {
			ns := strings.ToLower(strings.TrimSpace(line[len("name server:"):]))
			res.NameServers = append(res.NameServers, ns)
		}
	}

	return res
}
