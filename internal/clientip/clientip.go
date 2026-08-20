package clientip

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

// ClientIP extracts the real client IP from an HTTP request.
// It only trusts X-Forwarded-For / X-Real-IP if the direct RemoteAddr
// matches one of the trusted proxy CIDR prefixes.
func ClientIP(r *http.Request, trusted []netip.Prefix) netip.Addr {
	if r == nil {
		return netip.Addr{}
	}

	remoteAddrStr := r.RemoteAddr
	remoteHost, _, err := net.SplitHostPort(remoteAddrStr)
	if err != nil {
		remoteHost = remoteAddrStr
	}

	remoteHost = strings.TrimPrefix(remoteHost, "[")
	remoteHost = strings.TrimSuffix(remoteHost, "]")

	remoteIP, err := netip.ParseAddr(remoteHost)
	if err != nil {
		return netip.Addr{}
	}

	remoteIP = normalizeIP(remoteIP)

	// If remote IP is not from a trusted proxy, return direct remote IP
	if !isTrusted(remoteIP, trusted) {
		return remoteIP
	}

	// Remote IP is trusted. Check CF-Connecting-IP first (Cloudflare)
	if cf := r.Header.Get("CF-Connecting-IP"); cf != "" {
		if ip, err := netip.ParseAddr(strings.TrimSpace(cf)); err == nil {
			return normalizeIP(ip)
		}
	}

	// Check X-Forwarded-For next
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip := extractClientFromXFF(xff, trusted); ip.IsValid() {
			return ip
		}
	}

	// Check X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		xri = strings.TrimSpace(xri)
		if ip, err := netip.ParseAddr(xri); err == nil {
			return normalizeIP(ip)
		}
	}

	return remoteIP
}

// extractClientFromXFF traverses the XFF list from right to left,
// skipping trusted proxies to find the real client IP.
func extractClientFromXFF(xff string, trusted []netip.Prefix) netip.Addr {
	parts := strings.Split(xff, ",")
	var ips []netip.Addr

	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.TrimPrefix(p, "[")
		p = strings.TrimSuffix(p, "]")
		if ip, err := netip.ParseAddr(p); err == nil {
			ips = append(ips, normalizeIP(ip))
		}
	}

	if len(ips) == 0 {
		return netip.Addr{}
	}

	// Traverse from right to left (last hop backwards)
	for i := len(ips) - 1; i >= 0; i-- {
		if !isTrusted(ips[i], trusted) {
			return ips[i]
		}
	}

	// If all are within trusted ranges, return the leftmost IP
	return ips[0]
}

func isTrusted(ip netip.Addr, trusted []netip.Prefix) bool {
	if len(trusted) == 0 {
		return false
	}
	for _, prefix := range trusted {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

func normalizeIP(ip netip.Addr) netip.Addr {
	if ip.Is4In6() {
		return ip.Unmap()
	}
	return ip
}
