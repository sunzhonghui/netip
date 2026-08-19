package clientip

import (
	"net/http"
	"net/netip"
	"testing"
)

func TestClientIP(t *testing.T) {
	trusted := []netip.Prefix{
		netip.MustParsePrefix("10.0.0.0/8"),
		netip.MustParsePrefix("127.0.0.1/32"),
		netip.MustParsePrefix("172.16.0.0/12"),
		netip.MustParsePrefix("fd00::/8"),
	}

	tests := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		trusted    []netip.Prefix
		expected   string
	}{
		{
			name:       "Direct IPv4 connection without headers",
			remoteAddr: "203.0.113.195:45231",
			headers:    nil,
			trusted:    trusted,
			expected:   "203.0.113.195",
		},
		{
			name:       "Direct IPv6 connection without headers",
			remoteAddr: "[2001:db8::1]:54321",
			headers:    nil,
			trusted:    trusted,
			expected:   "2001:db8::1",
		},
		{
			name:       "Direct IPv4-mapped IPv6 connection",
			remoteAddr: "[::ffff:198.51.100.42]:12345",
			headers:    nil,
			trusted:    trusted,
			expected:   "198.51.100.42",
		},
		{
			name:       "Untrusted client attempts XFF spoofing",
			remoteAddr: "203.0.113.50:50123",
			headers: map[string]string{
				"X-Forwarded-For": "1.1.1.1, 8.8.8.8",
				"X-Real-IP":       "1.1.1.1",
			},
			trusted:  trusted,
			expected: "203.0.113.50", // spoofed header ignored
		},
		{
			name:       "Trusted proxy forwarding single client IP via XFF",
			remoteAddr: "10.0.0.2:8080",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.100",
			},
			trusted:  trusted,
			expected: "203.0.113.100",
		},
		{
			name:       "Trusted multi-hop proxies: client -> trusted proxy 1 -> trusted proxy 2 -> app",
			remoteAddr: "172.16.0.5:8080",
			headers: map[string]string{
				"X-Forwarded-For": "203.0.113.88, 10.0.0.10",
			},
			trusted:  trusted,
			expected: "203.0.113.88",
		},
		{
			name:       "Trusted multi-hop with spoofed prefix: spoofed -> real client -> trusted proxy",
			remoteAddr: "10.0.0.5:8080",
			headers: map[string]string{
				"X-Forwarded-For": "9.9.9.9, 203.0.113.77, 10.0.0.1",
			},
			trusted:  trusted,
			expected: "203.0.113.77",
		},
		{
			name:       "Trusted proxy with IPv6 client via XFF",
			remoteAddr: "10.0.0.2:8080",
			headers: map[string]string{
				"X-Forwarded-For": "240e:390:1234::1",
			},
			trusted:  trusted,
			expected: "240e:390:1234::1",
		},
		{
			name:       "Trusted proxy with X-Real-IP only",
			remoteAddr: "10.0.0.2:8080",
			headers: map[string]string{
				"X-Real-IP": "198.51.100.22",
			},
			trusted:  trusted,
			expected: "198.51.100.22",
		},
		{
			name:       "Trusted proxy but empty headers fallback to remote addr",
			remoteAddr: "10.0.0.2:8080",
			headers:    nil,
			trusted:    trusted,
			expected:   "10.0.0.2",
		},
		{
			name:       "No trusted proxies configured ignores all headers",
			remoteAddr: "10.0.0.2:8080",
			headers: map[string]string{
				"X-Forwarded-For": "1.1.1.1",
			},
			trusted:  nil,
			expected: "10.0.0.2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "http://example.com/", nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			req.RemoteAddr = tt.remoteAddr
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			ip := ClientIP(req, tt.trusted)
			if !ip.IsValid() {
				t.Fatalf("expected valid IP, got invalid")
			}
			if ip.String() != tt.expected {
				t.Errorf("expected IP %s, got %s", tt.expected, ip.String())
			}
		})
	}
}
