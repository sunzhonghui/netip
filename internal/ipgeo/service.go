package ipgeo

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/netip"
	"path/filepath"
	"sync"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"
)

type cacheEntry struct {
	result    *GeoResult
	expiresAt time.Time
}

// GeoService aggregates multiple geo providers with in-memory caching and online fallback.
type GeoService struct {
	providers  []GeoProvider
	mu         sync.RWMutex
	cache      map[string]cacheEntry
	httpClient *http.Client
}

// NewGeoService initializes GeoService searching for database files in dataDir.
func NewGeoService(dataDir string) *GeoService {
	svc := &GeoService{
		cache:      make(map[string]cacheEntry),
		httpClient: ssrf.NewSafeHTTPClient(config.DefaultHTTPTimeout),
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

// Reload reloads the geo database providers from dataDir dynamically.
func (s *GeoService) Reload(dataDir string) {
	ipdbDir := filepath.Join(dataDir, "ipdb")
	var newProviders []GeoProvider

	// 1. Try ip2region
	ip2regionPath := filepath.Join(ipdbDir, "ip2region.xdb")
	if p, err := NewIP2RegionProvider(ip2regionPath); err == nil {
		newProviders = append(newProviders, p)
		slog.Info("Reloaded ip2region database", "path", ip2regionPath)
	}

	// 2. Try MaxMind City
	maxmindPath := filepath.Join(ipdbDir, "GeoLite2-City.mmdb")
	if p, err := NewMaxMindCityProvider(maxmindPath); err == nil {
		newProviders = append(newProviders, p)
		slog.Info("Reloaded MaxMind City database", "path", maxmindPath)
	}

	s.mu.Lock()
	oldProviders := s.providers
	s.providers = newProviders
	s.cache = make(map[string]cacheEntry)
	s.mu.Unlock()

	for _, p := range oldProviders {
		_ = p.Close()
	}
}

// Close releases all underlying database resources.
func (s *GeoService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
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

	// 1. Try local databases (ip2region, MaxMind)
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

	// 2. If country is still empty (e.g. no local DB files mounted), fallback to online query
	if merged.Country == "" {
		if onlineRes := s.lookupOnline(addr); onlineRes != nil && onlineRes.Country != "" {
			merged.Country = onlineRes.Country
			merged.CountryCode = onlineRes.CountryCode
			merged.Province = onlineRes.Province
			merged.City = onlineRes.City
			merged.ISP = onlineRes.ISP
			merged.Source = onlineRes.Source
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

func (s *GeoService) lookupOnline(addr netip.Addr) *GeoResult {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultHTTPTimeout)
	defer cancel()

	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", addr.String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var data struct {
		Status      string `json:"status"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		RegionName  string `json:"regionName"`
		City        string `json:"city"`
		ISP         string `json:"isp"`
		Org         string `json:"org"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil || data.Status != "success" {
		return nil
	}

	isp := data.ISP
	if isp == "" {
		isp = data.Org
	}

	return &GeoResult{
		Country:     data.Country,
		CountryCode: data.CountryCode,
		Province:    data.RegionName,
		City:        data.City,
		ISP:         isp,
		Source:      "ip-api",
	}
}
