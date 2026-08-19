package asn

import (
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/oschwald/geoip2-golang"
)

// MaxMindASNProvider queries GeoLite2-ASN / GeoIP2-ASN MMDB.
type MaxMindASNProvider struct {
	db       *geoip2.Reader
	filePath string
}

// NewMaxMindASNProvider opens MaxMind ASN MMDB.
func NewMaxMindASNProvider(dbPath string) (*MaxMindASNProvider, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("maxmind asn database file not found at %s", dbPath)
	}

	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open maxmind asn mmdb: %w", err)
	}

	return &MaxMindASNProvider{
		db:       db,
		filePath: dbPath,
	}, nil
}

func (p *MaxMindASNProvider) Name() string {
	return "maxmind-asn"
}

func (p *MaxMindASNProvider) Available() bool {
	return p != nil && p.db != nil
}

func (p *MaxMindASNProvider) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

func (p *MaxMindASNProvider) Lookup(addr netip.Addr) (*ASNResult, error) {
	if !p.Available() {
		return nil, fmt.Errorf("maxmind asn not available")
	}

	if addr.Is4In6() {
		addr = addr.Unmap()
	}

	stdIP := net.ParseIP(addr.String())
	if stdIP == nil {
		return nil, fmt.Errorf("invalid ip %s", addr.String())
	}

	record, err := p.db.ASN(stdIP)
	if err != nil {
		return nil, err
	}

	return &ASNResult{
		ASN:     int(record.AutonomousSystemNumber),
		ASName:  record.AutonomousSystemOrganization,
		Network: "",
		Source:  "maxmind",
	}, nil
}

func (p *MaxMindASNProvider) LookupASN(asn int) (*ASNResult, error) {
	// MaxMind ASN mmdb is indexed by IP, not ASN number directly
	return nil, fmt.Errorf("lookup by ASN number not supported directly by mmdb")
}
