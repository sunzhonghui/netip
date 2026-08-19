package sslcheck

import (
	"context"
	"crypto/tls"
	"crypto/x509/pkix"
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"

	"golang.org/x/net/idna"
)

// SSLRequest holds host and port parameters for certificate inspection.
type SSLRequest struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

// SSLCertificateInfo details certificate attributes.
type SSLCertificateInfo struct {
	Subject       string    `json:"subject"`
	Issuer        string    `json:"issuer"`
	SerialNumber  string    `json:"serial_number"`
	DNSNames      []string  `json:"dns_names"`
	NotBefore     time.Time `json:"not_before"`
	NotAfter      time.Time `json:"not_after"`
	DaysRemaining int       `json:"days_remaining"`
}

// SSLCheckResult represents the full SSL/TLS diagnostic result.
type SSLCheckResult struct {
	Hostname      string               `json:"hostname"`
	Port          int                  `json:"port"`
	ResolvedIP    string               `json:"resolved_ip"`
	Valid         bool                 `json:"valid"`
	DaysRemaining int                  `json:"days_remaining"`
	Issuer        string               `json:"issuer"`
	Subject       string               `json:"subject"`
	DNSNames      []string             `json:"dns_names"`
	NotBefore     time.Time            `json:"not_before"`
	NotAfter      time.Time            `json:"not_after"`
	TLSVersion    string               `json:"tls_version"`
	CipherSuite   string               `json:"cipher_suite"`
	Certificates  []SSLCertificateInfo `json:"certificates"`
	Error         string               `json:"error,omitempty"`
}

// SSLService handles TLS certificate retrieval.
type SSLService struct {
	allowedPorts map[int]bool
}

// NewSSLService creates an SSLService.
func NewSSLService(allowedPorts map[int]bool) *SSLService {
	return &SSLService{
		allowedPorts: allowedPorts,
	}
}

// Inspect examines the SSL certificate of the given target.
func (s *SSLService) Inspect(ctx context.Context, req SSLRequest) (*SSLCheckResult, error) {
	hostname := strings.TrimSpace(req.Hostname)
	if strings.Contains(hostname, "://") {
		parts := strings.Split(hostname, "://")
		if len(parts) > 1 {
			hostname = parts[1]
		}
	}
	if slashIdx := strings.Index(hostname, "/"); slashIdx != -1 {
		hostname = hostname[:slashIdx]
	}
	if h, pStr, err := net.SplitHostPort(hostname); err == nil {
		hostname = h
		if req.Port <= 0 {
			var p int
			if _, err := fmt.Sscanf(pStr, "%d", &p); err == nil {
				req.Port = p
			}
		}
	}

	if req.Port <= 0 {
		req.Port = 443
	}

	puny, err := idna.ToASCII(hostname)
	if err == nil && puny != "" {
		hostname = puny
	}

	if err := ssrf.ValidateTargetHost(hostname); err != nil {
		return nil, err
	}

	if err := ssrf.ValidatePort(req.Port, s.allowedPorts); err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		ServerName:         hostname,
		InsecureSkipVerify: false,
	}

	tlsConn, resolvedIP, err := ssrf.SafeDialTLS(ctx, hostname, req.Port, s.allowedPorts, tlsConfig, config.DefaultTLSDialTimeout)
	if err != nil {
		// Try again with InsecureSkipVerify: true to still read expired or untrusted certs for diagnostics
		insecureConfig := &tls.Config{
			ServerName:         hostname,
			InsecureSkipVerify: true,
		}
		insecureConn, ip2, err2 := ssrf.SafeDialTLS(ctx, hostname, req.Port, s.allowedPorts, insecureConfig, config.DefaultTLSDialTimeout)
		if err2 != nil {
			return nil, fmt.Errorf("TLS connection failed: %w", err)
		}
		defer insecureConn.Close()
		return s.parseConnectionState(hostname, req.Port, ip2.String(), insecureConn.ConnectionState(), false, err.Error())
	}
	defer tlsConn.Close()

	return s.parseConnectionState(hostname, req.Port, resolvedIP.String(), tlsConn.ConnectionState(), true, "")
}

func (s *SSLService) parseConnectionState(hostname string, port int, resolvedIP string, state tls.ConnectionState, valid bool, errStr string) (*SSLCheckResult, error) {
	certs := state.PeerCertificates
	if len(certs) == 0 {
		return nil, fmt.Errorf("no certificates returned by peer")
	}

	leaf := certs[0]
	now := time.Now()
	daysRemaining := int(math.Floor(leaf.NotAfter.Sub(now).Hours() / 24.0))

	var certChain []SSLCertificateInfo
	for _, c := range certs {
		certChain = append(certChain, SSLCertificateInfo{
			Subject:       formatSubject(c.Subject),
			Issuer:        formatSubject(c.Issuer),
			SerialNumber:  c.SerialNumber.String(),
			DNSNames:      c.DNSNames,
			NotBefore:     c.NotBefore,
			NotAfter:      c.NotAfter,
			DaysRemaining: int(math.Floor(c.NotAfter.Sub(now).Hours() / 24.0)),
		})
	}

	return &SSLCheckResult{
		Hostname:      hostname,
		Port:          port,
		ResolvedIP:    resolvedIP,
		Valid:         valid && daysRemaining > 0,
		DaysRemaining: daysRemaining,
		Issuer:        formatSubject(leaf.Issuer),
		Subject:       formatSubject(leaf.Subject),
		DNSNames:      leaf.DNSNames,
		NotBefore:     leaf.NotBefore,
		NotAfter:      leaf.NotAfter,
		TLSVersion:    tlsVersionToString(state.Version),
		CipherSuite:   tls.CipherSuiteName(state.CipherSuite),
		Certificates:  certChain,
		Error:         errStr,
	}, nil
}

func formatSubject(pkix pkix.Name) string {
	var parts []string
	if pkix.CommonName != "" {
		parts = append(parts, "CN="+pkix.CommonName)
	}
	if len(pkix.Organization) > 0 {
		parts = append(parts, "O="+strings.Join(pkix.Organization, ", "))
	}
	if len(pkix.Country) > 0 {
		parts = append(parts, "C="+strings.Join(pkix.Country, ", "))
	}
	if len(parts) == 0 {
		return "Unknown"
	}
	return strings.Join(parts, ", ")
}

func tlsVersionToString(v uint16) string {
	switch v {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("0x%04x", v)
	}
}
