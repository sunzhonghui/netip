package ssrf

import (
	"net/netip"
	"testing"
)

func TestIsPublicIP(t *testing.T) {
	blocked := []string{
		"127.0.0.1",
		"127.0.0.2",
		"10.0.0.1",
		"10.255.255.255",
		"172.16.0.1",
		"172.31.255.255",
		"192.168.1.1",
		"169.254.169.254",
		"0.0.0.0",
		"100.64.0.1",
		"255.255.255.255",
		"::1",
		"::",
		"fe80::1",
		"fc00::1",
		"fd12:3456:789a::1",
		"ff02::1",
		"::ffff:127.0.0.1",
		"::ffff:10.0.0.1",
		"::ffff:192.168.1.1",
	}

	for _, ipStr := range blocked {
		addr, err := netip.ParseAddr(ipStr)
		if err != nil {
			t.Fatalf("failed to parse IP %s: %v", ipStr, err)
		}
		if IsPublicIP(addr) {
			t.Errorf("expected %s to be BLOCKED as non-public, but was allowed", ipStr)
		}
	}

	allowed := []string{
		"1.1.1.1",
		"8.8.8.8",
		"114.114.114.114",
		"223.5.5.5",
		"2606:4700:4700::1111",
		"2001:4860:4860::8888",
		"240e:390:1234::1",
	}

	for _, ipStr := range allowed {
		addr, err := netip.ParseAddr(ipStr)
		if err != nil {
			t.Fatalf("failed to parse IP %s: %v", ipStr, err)
		}
		if !IsPublicIP(addr) {
			t.Errorf("expected %s to be ALLOWED as public IP, but was blocked", ipStr)
		}
	}
}

func TestValidateTargetHost(t *testing.T) {
	blockedHosts := []string{
		"localhost",
		"sub.localhost",
		"server.local",
		"service.internal",
		"myrouter.lan",
		"metadata.google.internal",
		"127.0.0.1",
		"10.0.0.1",
		"192.168.1.1",
		"169.254.169.254",
		"[::1]",
	}

	for _, host := range blockedHosts {
		err := ValidateTargetHost(host)
		if err == nil {
			t.Errorf("expected host %s to be rejected by ValidateTargetHost, but got nil error", host)
		}
	}

	allowedHosts := []string{
		"1.1.1.1",
		"8.8.8.8",
		"example.com",
		"google.com",
		"cloudflare.com",
	}

	for _, host := range allowedHosts {
		err := ValidateTargetHost(host)
		if err != nil {
			t.Errorf("expected host %s to be allowed, got error: %v", host, err)
		}
	}
}

func TestValidatePort(t *testing.T) {
	allowedPorts := map[int]bool{
		80:  true,
		443: true,
		22:  true,
	}

	if err := ValidatePort(80, allowedPorts); err != nil {
		t.Errorf("expected port 80 to be valid, got: %v", err)
	}

	if err := ValidatePort(8080, allowedPorts); err == nil {
		t.Errorf("expected port 8080 to be rejected")
	}

	if err := ValidatePort(0, allowedPorts); err == nil {
		t.Errorf("expected port 0 to be rejected")
	}

	if err := ValidatePort(70000, allowedPorts); err == nil {
		t.Errorf("expected port 70000 to be rejected")
	}
}

func TestParseAndValidateURL(t *testing.T) {
	_, err := ParseAndValidateURL("http://127.0.0.1/admin")
	if err == nil {
		t.Errorf("expected loopback URL to be rejected")
	}

	_, err = ParseAndValidateURL("ftp://example.com/file")
	if err == nil {
		t.Errorf("expected ftp URL to be rejected")
	}

	u, err := ParseAndValidateURL("https://example.com/test?query=1")
	if err != nil {
		t.Errorf("expected valid URL to pass, got error: %v", err)
	}
	if u.Hostname() != "example.com" {
		t.Errorf("expected hostname example.com, got %s", u.Hostname())
	}
}
