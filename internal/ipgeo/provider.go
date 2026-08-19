package ipgeo

import (
	"net/netip"
)

// GeoResult holds normalized geographic and ISP details.
type GeoResult struct {
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Province    string `json:"province"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	Source      string `json:"source"`
}

// GeoProvider interface for IP geolocation lookups.
type GeoProvider interface {
	Lookup(addr netip.Addr) (*GeoResult, error)
	Name() string
	Available() bool
	Close() error
}
