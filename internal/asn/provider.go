package asn

import (
	"net/netip"
)

// ASNResult holds normalized autonomous system details.
type ASNResult struct {
	ASN      int    `json:"asn"`
	ASName   string `json:"as_name"`
	Country  string `json:"country"`
	Registry string `json:"registry"`
	Network  string `json:"network"`
	Source   string `json:"source"`
}

// ASNProvider interface for Autonomous System queries.
type ASNProvider interface {
	Lookup(addr netip.Addr) (*ASNResult, error)
	LookupASN(asn int) (*ASNResult, error)
	Name() string
	Available() bool
	Close() error
}
