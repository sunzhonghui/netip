package webcheck

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"sync"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"

	"golang.org/x/net/idna"
)

// EndpointStatus represents connection state for a protocol/IP version.
type EndpointStatus struct {
	Supported  bool    `json:"supported"`
	StatusCode int     `json:"status_code,omitempty"`
	LatencyMs  int64   `json:"latency_ms,omitempty"`
	Error      string  `json:"error,omitempty"`
}

// ProtocolCheckResult contains DNS and HTTP/HTTPS status for an IP version.
type ProtocolCheckResult struct {
	DNS       bool           `json:"dns"`
	Addresses []string       `json:"addresses"`
	HTTP      EndpointStatus `json:"http"`
	HTTPS     EndpointStatus `json:"https"`
}

// IPv6CheckResponse represents the full diagnostic summary for IPv6 readiness.
type IPv6CheckResponse struct {
	Domain     string              `json:"domain"`
	IPv4       ProtocolCheckResult `json:"ipv4"`
	IPv6       ProtocolCheckResult `json:"ipv6"`
	Supported  bool                `json:"supported"`
	Conclusion string              `json:"conclusion"`
}

// IPv6Checker checks dual-stack IPv4/IPv6 capabilities of a website.
type IPv6Checker struct{}

// NewIPv6Checker creates an IPv6Checker.
func NewIPv6Checker() *IPv6Checker {
	return &IPv6Checker{}
}

// Check tests DNS A/AAAA and HTTP/HTTPS connectivity over IPv4 and IPv6.
func (c *IPv6Checker) Check(ctx context.Context, domain string) (*IPv6CheckResponse, error) {
	domain = strings.TrimSpace(domain)
	// Strip scheme if user provided http:// or https://
	if strings.Contains(domain, "://") {
		parts := strings.Split(domain, "://")
		if len(parts) > 1 {
			domain = parts[1]
		}
	}
	if slashIdx := strings.Index(domain, "/"); slashIdx != -1 {
		domain = domain[:slashIdx]
	}
	if h, _, err := net.SplitHostPort(domain); err == nil {
		domain = h
	}

	puny, err := idna.ToASCII(domain)
	if err == nil && puny != "" {
		domain = puny
	}

	if err := ssrf.ValidateTargetHost(domain); err != nil {
		return nil, err
	}

	checkCtx, cancel := context.WithTimeout(ctx, config.DefaultIPv6CheckTimeout)
	defer cancel()

	var resolver net.Resolver
	var ipv4Addrs, ipv6Addrs []netip.Addr

	// 1. Resolve IPv4 (A)
	if ips, err := resolver.LookupNetIP(checkCtx, "ip4", domain); err == nil {
		for _, ip := range ips {
			if ssrf.IsPublicIP(ip) {
				ipv4Addrs = append(ipv4Addrs, ip)
			}
		}
	}

	// 2. Resolve IPv6 (AAAA)
	if ips, err := resolver.LookupNetIP(checkCtx, "ip6", domain); err == nil {
		for _, ip := range ips {
			if ssrf.IsPublicIP(ip) {
				ipv6Addrs = append(ipv6Addrs, ip)
			}
		}
	}

	resp := &IPv6CheckResponse{
		Domain: domain,
		IPv4: ProtocolCheckResult{
			DNS:       len(ipv4Addrs) > 0,
			Addresses: formatIPList(ipv4Addrs),
		},
		IPv6: ProtocolCheckResult{
			DNS:       len(ipv6Addrs) > 0,
			Addresses: formatIPList(ipv6Addrs),
		},
	}

	var wg sync.WaitGroup

	// Test IPv4 HTTP / HTTPS
	if len(ipv4Addrs) > 0 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			resp.IPv4.HTTP = checkHTTP(checkCtx, domain, ipv4Addrs[0], 80, false)
		}()
		go func() {
			defer wg.Done()
			resp.IPv4.HTTPS = checkHTTP(checkCtx, domain, ipv4Addrs[0], 443, true)
		}()
	}

	// Test IPv6 HTTP / HTTPS
	if len(ipv6Addrs) > 0 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			resp.IPv6.HTTP = checkHTTP(checkCtx, domain, ipv6Addrs[0], 80, false)
		}()
		go func() {
			defer wg.Done()
			resp.IPv6.HTTPS = checkHTTP(checkCtx, domain, ipv6Addrs[0], 443, true)
		}()
	}

	wg.Wait()

	// Compute conclusion
	if resp.IPv6.DNS && (resp.IPv6.HTTPS.Supported || resp.IPv6.HTTP.Supported) {
		resp.Supported = true
		resp.Conclusion = "该网站完整支持 IPv6"
	} else if resp.IPv6.DNS {
		resp.Supported = false
		resp.Conclusion = "网站已配置 IPv6 DNS (AAAA 记录)，但 IPv6 服务无法连接或超时"
	} else {
		resp.Supported = false
		resp.Conclusion = "该网站未配置 IPv6 (无 AAAA 记录)"
	}

	return resp, nil
}

func checkHTTP(ctx context.Context, hostname string, targetIP netip.Addr, port int, useTLS bool) EndpointStatus {
	start := time.Now()
	targetAddr := net.JoinHostPort(targetIP.String(), fmt.Sprintf("%d", port))

	var dialer net.Dialer
	dialer.Timeout = 3 * time.Second

	conn, err := dialer.DialContext(ctx, "tcp", targetAddr)
	if err != nil {
		return EndpointStatus{
			Supported: false,
			Error:     err.Error(),
		}
	}
	defer conn.Close()

	if useTLS {
		tlsConn := tls.Client(conn, &tls.Config{
			ServerName:         hostname,
			InsecureSkipVerify: false,
		})
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			return EndpointStatus{
				Supported: false,
				Error:     "TLS handshake failed: " + err.Error(),
			}
		}
		conn = tlsConn
	}

	// Send minimal HTTP HEAD request
	reqStr := fmt.Sprintf("HEAD / HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nConnection: close\r\n\r\n", hostname, ssrf.UserAgent)
	_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := conn.Write([]byte(reqStr)); err != nil {
		return EndpointStatus{
			Supported: false,
			Error:     err.Error(),
		}
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	latency := time.Since(start).Milliseconds()

	if err != nil && n == 0 {
		return EndpointStatus{
			Supported: false,
			Error:     err.Error(),
		}
	}

	respStr := string(buf[:n])
	statusCode := 200
	if strings.HasPrefix(respStr, "HTTP/") {
		parts := strings.Split(respStr, " ")
		if len(parts) > 1 {
			var code int
			if _, err := fmt.Sscanf(parts[1], "%d", &code); err == nil {
				statusCode = code
			}
		}
	}

	return EndpointStatus{
		Supported:  true,
		StatusCode: statusCode,
		LatencyMs:  latency,
	}
}

func formatIPList(addrs []netip.Addr) []string {
	var list []string
	for _, a := range addrs {
		list = append(list, a.String())
	}
	return list
}
