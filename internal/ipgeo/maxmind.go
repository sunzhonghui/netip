package ipgeo

import (
	"fmt"
	"net"
	"net/netip"
	"os"

	"github.com/oschwald/geoip2-golang"
)

// MaxMindCityProvider provides geolocation from GeoLite2-City / GeoIP2-City MMDB.
type MaxMindCityProvider struct {
	db       *geoip2.Reader
	filePath string
}

// NewMaxMindCityProvider opens a MaxMind City MMDB file.
func NewMaxMindCityProvider(dbPath string) (*MaxMindCityProvider, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("maxmind city database file not found at %s", dbPath)
	}

	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open maxmind mmdb: %w", err)
	}

	return &MaxMindCityProvider{
		db:       db,
		filePath: dbPath,
	}, nil
}

func (p *MaxMindCityProvider) Name() string {
	return "maxmind"
}

func (p *MaxMindCityProvider) Available() bool {
	return p != nil && p.db != nil
}

func (p *MaxMindCityProvider) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}

// Lookup queries MaxMind City mmdb.
func (p *MaxMindCityProvider) Lookup(addr netip.Addr) (*GeoResult, error) {
	if !p.Available() {
		return nil, fmt.Errorf("maxmind provider not available")
	}

	if addr.Is4In6() {
		addr = addr.Unmap()
	}

	stdIP := net.ParseIP(addr.String())
	if stdIP == nil {
		return nil, fmt.Errorf("invalid ip %s", addr.String())
	}

	record, err := p.db.City(stdIP)
	if err != nil {
		return nil, err
	}

	res := &GeoResult{
		Source: "maxmind",
	}

	// Country
	if record.Country.Names != nil {
		if zh, ok := record.Country.Names["zh-CN"]; ok && zh != "" {
			res.Country = zh
		} else if en, ok := record.Country.Names["en"]; ok {
			res.Country = en
		}
	}
	res.CountryCode = record.Country.IsoCode

	// Province / Subdivision
	if len(record.Subdivisions) > 0 {
		sub := record.Subdivisions[0]
		if zh, ok := sub.Names["zh-CN"]; ok && zh != "" {
			res.Province = zh
		} else if en, ok := sub.Names["en"]; ok {
			res.Province = en
		}
	}

	// City
	if record.City.Names != nil {
		if zh, ok := record.City.Names["zh-CN"]; ok && zh != "" {
			res.City = zh
		} else if en, ok := record.City.Names["en"]; ok {
			res.City = en
		}
	}

	return res, nil
}
