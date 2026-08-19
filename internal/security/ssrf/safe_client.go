package ssrf

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"strconv"
	"time"

	"netip/internal/config"
)

const UserAgent = "NetIP-Network-Diagnostics/1.0"

// SafeTransport creates an http.Transport with SSRF & DNS rebinding protection.
func NewSafeTransport(dialTimeout time.Duration) *http.Transport {
	if dialTimeout <= 0 {
		dialTimeout = 5 * time.Second
	}

	return &http.Transport{
		Proxy: nil, // Do not use system environment proxies for internal security checks
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, portStr, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, fmt.Errorf("invalid address %s: %w", addr, err)
			}

			// Validate and resolve to public IP addresses only
			safeIPs, err := ResolveSafeIPs(ctx, host)
			if err != nil {
				return nil, err
			}

			var dialer net.Dialer
			dialer.Timeout = dialTimeout

			// Try the first available safe IP
			var lastErr error
			for _, ip := range safeIPs {
				targetAddr := net.JoinHostPort(ip.String(), portStr)
				conn, err := dialer.DialContext(ctx, network, targetAddr)
				if err == nil {
					return conn, nil
				}
				lastErr = err
			}

			return nil, fmt.Errorf("dial failed: %w", lastErr)
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   config.DefaultTLSDialTimeout,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     true, // Diagnostics should not reuse stale connections
	}
}

// NewSafeHTTPClient creates an http.Client with redirect validation and SSRF defenses.
func NewSafeHTTPClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = config.DefaultHTTPTimeout
	}

	transport := NewSafeTransport(5 * time.Second)

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= config.DefaultMaxHTTPRedirects {
				return ErrTooManyRedirects
			}

			// Validate redirect destination
			if err := ValidateTargetHost(req.URL.Host); err != nil {
				return fmt.Errorf("redirect target blocked: %w", err)
			}

			// Ensure safe scheme
			if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
				return fmt.Errorf("invalid redirect scheme: %s", req.URL.Scheme)
			}

			// Set user agent
			req.Header.Set("User-Agent", UserAgent)
			return nil
		},
	}
}

// SafeDialTCP dials a TCP address after validating SSRF and port.
func SafeDialTCP(ctx context.Context, host string, port int, allowedPorts map[int]bool, timeout time.Duration) (net.Conn, netip.Addr, error) {
	if err := ValidatePort(port, allowedPorts); err != nil {
		return nil, netip.Addr{}, err
	}

	safeIPs, err := ResolveSafeIPs(ctx, host)
	if err != nil {
		return nil, netip.Addr{}, err
	}

	var dialer net.Dialer
	dialer.Timeout = timeout

	var lastErr error
	for _, ip := range safeIPs {
		targetAddr := net.JoinHostPort(ip.String(), strconv.Itoa(port))
		conn, err := dialer.DialContext(ctx, "tcp", targetAddr)
		if err == nil {
			return conn, ip, nil
		}
		lastErr = err
	}

	return nil, netip.Addr{}, fmt.Errorf("connection failed: %w", lastErr)
}

// SafeDialTLS connects securely via TLS after validating SSRF.
func SafeDialTLS(ctx context.Context, host string, port int, allowedPorts map[int]bool, tlsConfig *tls.Config, timeout time.Duration) (*tls.Conn, netip.Addr, error) {
	rawConn, ip, err := SafeDialTCP(ctx, host, port, allowedPorts, timeout)
	if err != nil {
		return nil, netip.Addr{}, err
	}

	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
		}
	} else if tlsConfig.ServerName == "" {
		tlsConfig = tlsConfig.Clone()
		tlsConfig.ServerName = host
	}

	tlsConn := tls.Client(rawConn, tlsConfig)
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		rawConn.Close()
		return nil, ip, fmt.Errorf("tls handshake failed: %w", err)
	}

	return tlsConn, ip, nil
}

// ParseAndValidateURL parses raw URL string and checks for safe scheme and host.
func ParseAndValidateURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidTarget, err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("%w: scheme must be http or https", ErrInvalidTarget)
	}

	if u.Hostname() == "" {
		return nil, fmt.Errorf("%w: missing host", ErrInvalidTarget)
	}

	if err := ValidateTargetHost(u.Hostname()); err != nil {
		return nil, err
	}

	return u, nil
}
