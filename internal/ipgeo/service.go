package ipgeo

import (
	"log/slog"
	"net/netip"
	"path/filepath"
	"sync"
	"time"
)

type cacheEntry struct {
	result    *GeoResult
	expiresAt time.Time
}

// GeoService aggregates multiple geo providers with in-memory caching.
type GeoService struct {
	providers []GeoProvider
	mu        sync.RWMutex
	cache     map[string]cacheEntry
}

// NewGeoService initializes GeoService searching for database files in dataDir.
func NewGeoService(dataDir string) *GeoService {
	svc := &GeoService{
		cache: make(map[string]cacheEntry),
	}

	ipdbDir := filepath.Join(dataDir, "ipdb")

	// 1. Try ip2region
	ip2regionPath := filepath.Join(ipdbDir, "ip2region.xdb")
	if p, err := NewIP2RegionProvider(ip2regionPath); err == nil {
		svc.providers = append(svc.providers, p)
		slog.Info("Loaded ip2region database", "path", ip2regionPath)
	} else {
		slog.Info("ip2region database not loaded (optional)", "path", ip2regionPath, "err", err.Error())
	}

	// 2. Try MaxMind City
	maxmindPath := filepath.Join(ipdbDir, "GeoLite2-City.mmdb")
	if p, err := NewMaxMindCityProvider(maxmindPath); err == nil {
		svc.providers = append(svc.providers, p)
		slog.Info("Loaded MaxMind City database", "path", maxmindPath)
	} else {
		slog.Info("MaxMind City database not loaded (optional)", "path", maxmindPath, "err", err.Error())
	}

	return svc
}

// Close releases all underlying database resources.
func (s *GeoService) Close() {
	for _, p := range s.providers {
		_ = p.Close()
	}
}

// HasProviders returns whether at least one geo database is loaded.
func (s *GeoService) HasProviders() bool {
	return len(s.providers) > 0
}

// Lookup queries all available providers and combines results.
func (s *GeoService) Lookup(addr netip.Addr) *GeoResult {
	if !addr.IsValid() {
		return &GeoResult{Source: "none"}
	}

	ipStr := addr.String()

	// Check cache
	s.mu.RLock()
	if entry, ok := s.cache[ipStr]; ok && time.Now().Before(entry.expiresAt) {
		s.mu.RUnlock()
		return entry.result
	}
	s.mu.RUnlock()

	merged := &GeoResult{
		Source: "none",
	}

	for _, p := range s.providers {
		if !p.Available() {
			continue
		}

		res, err := p.Lookup(addr)
		if err != nil || res == nil {
			continue
		}

		// Prefer values from earlier providers (e.g. ip2region for CN ISP)
		if merged.Country == "" && res.Country != "" {
			merged.Country = res.Country
		}
		if merged.CountryCode == "" && res.CountryCode != "" {
			merged.CountryCode = res.CountryCode
		}
		if merged.Province == "" && res.Province != "" {
			merged.Province = res.Province
		}
		if merged.City == "" && res.City != "" {
			merged.City = res.City
		}
		if merged.ISP == "" && res.ISP != "" {
			merged.ISP = res.ISP
		}
		if merged.Source == "none" {
			merged.Source = res.Source
		} else if merged.Source != res.Source {
			merged.Source = merged.Source + "+" + res.Source
		}
	}

	// Cache result for 1 hour
	s.mu.Lock()
	if len(s.cache) > 20000 {
		// Quick prune on high memory
		s.cache = make(map[string]cacheEntry)
	}
	s.cache[ipStr] = cacheEntry{
		result:    merged,
		expiresAt: time.Now().Add(1 * time.Hour),
	}
	s.mu.Unlock()

	return merged
}
