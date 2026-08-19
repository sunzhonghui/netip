package speedtest

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"
)

// SpeedTestResult holds speed test diagnostics.
type SpeedTestResult struct {
	Target        string  `json:"target"`
	DNSMs         int64   `json:"dns_ms"`
	ConnectMs     int64   `json:"connect_ms"`
	TLSMs         int64   `json:"tls_ms"`
	TTFBMs        int64   `json:"ttfb_ms"`
	DownloadBytes int64   `json:"download_bytes"`
	DownloadMs    int64   `json:"download_ms"`
	SpeedMbps     float64 `json:"speed_mbps"`
	ResolvedIP    string  `json:"resolved_ip"`
}

// SpeedService performs controlled website HTTP speed testing.
type SpeedService struct {
	transport *http.Transport
}

// NewSpeedService creates a SpeedService.
func NewSpeedService() *SpeedService {
	return &SpeedService{
		transport: ssrf.NewSafeTransport(config.DefaultSpeedTestOverallTimeout),
	}
}

// Test performs an HTTP download speed test with a 5MB and 10s ceiling.
func (s *SpeedService) Test(ctx context.Context, targetURL string) (*SpeedTestResult, error) {
	targetURL = strings.TrimSpace(targetURL)
	if !strings.HasPrefix(targetURL, "http://") && !strings.HasPrefix(targetURL, "https://") {
		targetURL = "https://" + targetURL
	}

	u, err := ssrf.ParseAndValidateURL(targetURL)
	if err != nil {
		return nil, err
	}

	testCtx, cancel := context.WithTimeout(ctx, config.DefaultSpeedTestOverallTimeout)
	defer cancel()

	var (
		t0, tDNS, tConnect, tTLS, tTTFB time.Time
		dnsStart, connectStart, tlsStart time.Time
		resolvedIP                       string
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
			connectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			tConnect = time.Now()
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

	req, err := http.NewRequestWithContext(httptrace.WithClientTrace(testCtx, trace), "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", ssrf.UserAgent)
	req.Header.Set("Accept", "*/*")

	client := &http.Client{
		Transport: s.transport,
		Timeout:   config.DefaultSpeedTestOverallTimeout,
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
		return nil, fmt.Errorf("speed test failed: %w", err)
	}
	defer resp.Body.Close()

	// Measure download timing capped at 5 MB
	downloadStart := time.Now()
	limitedReader := io.LimitReader(resp.Body, config.DefaultSpeedTestMaxBytes)
	buf := make([]byte, 32*1024)
	var totalRead int64

	for {
		n, err := limitedReader.Read(buf)
		if n > 0 {
			totalRead += int64(n)
		}
		if err != nil {
			break
		}
	}

	downloadDuration := time.Since(downloadStart)
	downloadMs := downloadDuration.Milliseconds()
	if downloadMs == 0 {
		downloadMs = 1
	}

	// Calculate Mbps: (bytes * 8) / (seconds * 1,000,000)
	seconds := downloadDuration.Seconds()
	var speedMbps float64
	if seconds > 0 && totalRead > 0 {
		speedMbps = float64(totalRead*8) / (seconds * 1000000.0)
		speedMbps = math.Round(speedMbps*10.0) / 10.0
	}

	var dnsMs, connectMs, tlsMs, ttfbMs int64
	if !tDNS.IsZero() && !dnsStart.IsZero() {
		dnsMs = tDNS.Sub(dnsStart).Milliseconds()
	}
	if !tConnect.IsZero() && !connectStart.IsZero() {
		connectMs = tConnect.Sub(connectStart).Milliseconds()
	}
	if !tTLS.IsZero() && !tlsStart.IsZero() {
		tlsMs = tTLS.Sub(tlsStart).Milliseconds()
	}
	if !tTTFB.IsZero() {
		ttfbMs = tTTFB.Sub(t0).Milliseconds()
	}

	return &SpeedTestResult{
		Target:        u.String(),
		DNSMs:         dnsMs,
		ConnectMs:     connectMs,
		TLSMs:         tlsMs,
		TTFBMs:        ttfbMs,
		DownloadBytes: totalRead,
		DownloadMs:    downloadMs,
		SpeedMbps:     speedMbps,
		ResolvedIP:    resolvedIP,
	}, nil
}
