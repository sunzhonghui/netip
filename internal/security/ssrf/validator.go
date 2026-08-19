package ssrf

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strings"
)

var (
	ErrBlockedPrivateIP   = errors.New("target resolves to a private or reserved network address")
	ErrBlockedHostname    = errors.New("target hostname is blocked")
	ErrNoPublicIPResolved = errors.New("no safe public IP address found for target")
	ErrInvalidTarget      = errors.New("target is invalid")
	ErrPortNotAllowed     = errors.New("target port is not permitted")
	ErrTooManyRedirects   = errors.New("too many redirects")
)

// List of blocked IP prefixes for IPv4 and IPv6
var blockedPrefixes = []netip.Prefix{
	// IPv4 Loopback
	netip.MustParsePrefix("127.0.0.0/8"),
	// IPv4 Private networks
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.168.0.0/16"),
	// IPv4 Link-local
	netip.MustParsePrefix("169.254.0.0/16"),
	// IPv4 Current / Unspecified
	netip.MustParsePrefix("0.0.0.0/8"),
	// IPv4 Carrier-grade NAT
	netip.MustParsePrefix("100.64.0.0/10"),
	// IPv4 Benchmark & Reserved
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("255.255.255.255/32"),

	// IPv6 Loopback & Unspecified
	netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("::/128"),
	// IPv6 Unique Local Address (ULA)
	netip.MustParsePrefix("fc00::/7"),
	// IPv6 Link-Local
	netip.MustParsePrefix("fe80::/10"),
	// IPv6 Multicast
	netip.MustParsePrefix("ff00::/8"),
	// IPv6 Discard prefix
	netip.MustParsePrefix("100::/64"),
	// IPv6 Documentation prefix
	netip.MustParsePrefix("2001:db8::/32"),
}

// Blocked hostnames / domain suffixes
var blockedHostnames = []string{
	"localhost",
	"ip6-localhost",
	"ip6-loopback",
	"metadata.google.internal",
	"169.254.169.254",
	"instance-data",
}

// IsPublicIP verifies if an IP address is a globally routable, non-private IP.
func IsPublicIP(addr netip.Addr) bool {
	if !addr.IsValid() {
		return false
	}

	// Normalize IPv4-mapped IPv6 (::ffff:127.0.0.1 -> 127.0.0.1)
	if addr.Is4In6() {
		addr = addr.Unmap()
	}

	if addr.IsLoopback() || addr.IsPrivate() || addr.IsLinkLocalUnicast() ||
		addr.IsLinkLocalMulticast() || addr.IsMulticast() || addr.IsUnspecified() {
		return false
	}

	for _, prefix := range blockedPrefixes {
		if prefix.Contains(addr) {
			return false
		}
	}

	return true
}

// ValidateTargetHost checks if a target hostname or raw IP is safe.
func ValidateTargetHost(host string) error {
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "" {
		return ErrInvalidTarget
	}

	// Strip port if present
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	host = strings.TrimPrefix(host, "[")
	host = strings.TrimSuffix(host, "]")

	for _, blocked := range blockedHostnames {
		if host == blocked || strings.HasSuffix(host, "."+blocked) {
			return fmt.Errorf("%w: %s", ErrBlockedHostname, host)
		}
	}

	if strings.HasSuffix(host, ".local") ||
		strings.HasSuffix(host, ".internal") ||
		strings.HasSuffix(host, ".lan") ||
		strings.HasSuffix(host, ".home.arpa") {
		return fmt.Errorf("%w: internal domain %s", ErrBlockedHostname, host)
	}

	// If it's a literal IP
	if ip, err := netip.ParseAddr(host); err == nil {
		if !IsPublicIP(ip) {
			return fmt.Errorf("%w: %s", ErrBlockedPrivateIP, host)
		}
	}

	return nil
}

// ResolveSafeIPs resolves a host to its public IP addresses.
// It filters out any internal/private/loopback addresses.
func ResolveSafeIPs(ctx context.Context, host string) ([]netip.Addr, error) {
	if err := ValidateTargetHost(host); err != nil {
		return nil, err
	}

	cleanHost := host
	if h, _, err := net.SplitHostPort(host); err == nil {
		cleanHost = h
	}
	cleanHost = strings.TrimPrefix(cleanHost, "[")
	cleanHost = strings.TrimSuffix(cleanHost, "]")

	// If already a literal IP
	if ip, err := netip.ParseAddr(cleanHost); err == nil {
		if !IsPublicIP(ip) {
			return nil, fmt.Errorf("%w: %s", ErrBlockedPrivateIP, cleanHost)
		}
		return []netip.Addr{ip}, nil
	}

	var resolver net.Resolver
	addrs, err := resolver.LookupNetIP(ctx, "ip", cleanHost)
	if err != nil {
		return nil, fmt.Errorf("dns resolution failed: %w", err)
	}

	var safeIPs []netip.Addr
	for _, addr := range addrs {
		if IsPublicIP(addr) {
			safeIPs = append(safeIPs, addr)
		}
	}

	if len(safeIPs) == 0 {
		return nil, fmt.Errorf("%w for host %s", ErrNoPublicIPResolved, host)
	}

	return safeIPs, nil
}

// ValidatePort checks whether a target port is in the allowed port map.
func ValidatePort(port int, allowedPorts map[int]bool) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("%w: port %d out of range", ErrInvalidTarget, port)
	}
	if len(allowedPorts) > 0 && !allowedPorts[port] {
		return fmt.Errorf("%w: port %d is not in allowed list", ErrPortNotAllowed, port)
	}
	return nil
}
