package ipgeo

import (
	"fmt"
	"net/netip"
	"os"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// IP2RegionProvider searches IPv4 locations via ip2region xdb.
type IP2RegionProvider struct {
	searcher *xdb.Searcher
	filePath string
}

// NewIP2RegionProvider initializes an ip2region searcher if the file exists.
func NewIP2RegionProvider(dbPath string) (*IP2RegionProvider, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("ip2region database file not found at %s", dbPath)
	}

	ver, err := xdb.VersionFromName("IPv4")
	if err != nil {
		return nil, fmt.Errorf("failed to get xdb ipv4 version: %w", err)
	}

	cBuff, err := xdb.LoadContentFromFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load ip2region buffer: %w", err)
	}

	searcher, err := xdb.NewWithBuffer(ver, cBuff)
	if err != nil {
		return nil, fmt.Errorf("failed to create ip2region searcher: %w", err)
	}

	return &IP2RegionProvider{
		searcher: searcher,
		filePath: dbPath,
	}, nil
}

func (p *IP2RegionProvider) Name() string {
	return "ip2region"
}

func (p *IP2RegionProvider) Available() bool {
	return p != nil && p.searcher != nil
}

func (p *IP2RegionProvider) Close() error {
	if p.searcher != nil {
		p.searcher.Close()
	}
	return nil
}

// Lookup queries the ip2region xdb database.
func (p *IP2RegionProvider) Lookup(addr netip.Addr) (*GeoResult, error) {
	if !p.Available() {
		return nil, fmt.Errorf("ip2region not available")
	}

	if addr.Is4In6() {
		addr = addr.Unmap()
	}

	// ip2region xdb only supports IPv4
	if !addr.Is4() {
		return nil, fmt.Errorf("ip2region only supports IPv4")
	}

	region, err := p.searcher.Search(addr.String())
	if err != nil {
		return nil, err
	}

	// ip2region format: Country|Region|Province|City|ISP
	// Example: 中国|0|广东省|广州市|电信
	parts := strings.Split(region, "|")
	res := &GeoResult{
		Source: "ip2region",
	}

	if len(parts) > 0 && parts[0] != "0" {
		res.Country = parts[0]
		if parts[0] == "中国" {
			res.CountryCode = "CN"
		}
	}
	if len(parts) > 2 && parts[2] != "0" {
		res.Province = parts[2]
	}
	if len(parts) > 3 && parts[3] != "0" {
		res.City = parts[3]
	}
	if len(parts) > 4 && parts[4] != "0" {
		res.ISP = parts[4]
	}

	return res, nil
}
