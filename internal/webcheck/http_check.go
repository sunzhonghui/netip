package webcheck

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"
)

// HTTPCheckResult represents the full HTTP site diagnostic breakdown.
type HTTPCheckResult struct {
	URL           string            `json:"url"`
	DNSMs         int64             `json:"dns_ms"`
	TCPMs         int64             `json:"tcp_ms"`
	TLSMs         int64             `json:"tls_ms"`
	TTFBMs        int64             `json:"ttfb_ms"`
	TotalMs       int64             `json:"total_ms"`
	Protocol      string            `json:"protocol"`
	StatusCode    int               `json:"status_code"`
	StatusText    string            `json:"status_text"`
	ResolvedIP    string            `json:"resolved_ip"`
	Server        string            `json:"server,omitempty"`
	ContentType   string            `json:"content_type,omitempty"`
	ContentLength int64             `json:"content_length,omitempty"`
	Headers       map[string]string `json:"headers,omitempty"`
}

// HTTPChecker handles detailed website connectivity diagnostics.
type HTTPChecker struct {
	transport *http.Transport
}

// NewHTTPChecker creates an HTTPChecker.
func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		transport: ssrf.NewSafeTransport(config.DefaultHTTPTimeout),
	}
}

// Check performs a timed diagnostic request against the target URL.
func (c *HTTPChecker) Check(ctx context.Context, targetURL string) (*HTTPCheckResult, error) {
	targetURL = strings.TrimSpace(targetURL)
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	u, err := ssrf.ParseAndValidateURL(targetURL)
	if err != nil {
		return nil, err
	}

	checkCtx, cancel := context.WithTimeout(ctx, config.DefaultHTTPTimeout)
	defer cancel()

	var (
		t0, tDNS, tTCP, tTLS, tTTFB, tDone time.Time
		dnsStart, tcpStart, tlsStart       time.Time
		resolvedIP                         string
	)

	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now()
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			tDNS = time.Now()
			if len(info.Addrs) > 0 {
				resolvedIP = info.Addrs[0].String()
			}
		},
		ConnectStart: func(network, addr string) {
			tcpStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			tTCP = time.Now()
		},
		TLSHandshakeStart: func() {
			tlsStart = time.Now()
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			tTLS = time.Now()
		},
		GotFirstResponseByte: func() {
			tTTFB = time.Now()
		},
	}

	req, err := http.NewRequestWithContext(httptrace.WithClientTrace(checkCtx, trace), "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", ssrf.UserAgent)
	req.Header.Set("Accept", "*/*")

	client := &http.Client{
		Transport: c.transport,
		Timeout:   config.DefaultHTTPTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= config.DefaultMaxHTTPRedirects {
				return ssrf.ErrTooManyRedirects
			}
			return ssrf.ValidateTargetHost(req.URL.Host)
		},
	}

	t0 = time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP check failed: %w", err)
	}
	defer resp.Body.Close()

	// Read small response body (up to 64KB)
	_, _ = io.CopyN(io.Discard, resp.Body, 64*1024)
	tDone = time.Now()

	// Compute intervals
	var dnsMs, tcpMs, tlsMs, ttfbMs int64
	if !tDNS.IsZero() && !dnsStart.IsZero() {
		dnsMs = tDNS.Sub(dnsStart).Milliseconds()
	}
	if !tTCP.IsZero() && !tcpStart.IsZero() {
		tcpMs = tTCP.Sub(tcpStart).Milliseconds()
	}
	if !tTLS.IsZero() && !tlsStart.IsZero() {
		tlsMs = tTLS.Sub(tlsStart).Milliseconds()
	}
	if !tTTFB.IsZero() {
		ttfbMs = tTTFB.Sub(t0).Milliseconds()
	}

	totalMs := tDone.Sub(t0).Milliseconds()

	// Filter key headers
	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &HTTPCheckResult{
		URL:           u.String(),
		DNSMs:         dnsMs,
		TCPMs:         tcpMs,
		TLSMs:         tlsMs,
		TTFBMs:        ttfbMs,
		TotalMs:       totalMs,
		Protocol:      resp.Proto,
		StatusCode:    resp.StatusCode,
		StatusText:    resp.Status,
		ResolvedIP:    resolvedIP,
		Server:        resp.Header.Get("Server"),
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
		Headers:       headers,
	}, nil
}
