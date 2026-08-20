package asn

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/netip"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"
)

type cacheEntry struct {
	result    *ASNResult
	expiresAt time.Time
}

// ASNService manages ASN lookups via local mmdb and online fallback.
type ASNService struct {
	providers  []ASNProvider
	mu         sync.RWMutex
	cache      map[string]cacheEntry
	httpClient *http.Client
}

// Well-known popular ASNs for instant local resolution
var wellKnownASNs = map[int]ASNResult{
	4134:   {ASN: 4134, ASName: "CHINANET-BACKBONE (China Telecom)", Country: "CN", Registry: "APNIC"},
	4809:   {ASN: 4809, ASName: "China Telecom CN2", Country: "CN", Registry: "APNIC"},
	4837:   {ASN: 4837, ASName: "CHINA169-BACKBONE (China Unicom)", Country: "CN", Registry: "APNIC"},
	9929:   {ASN: 9929, ASName: "China Unicom Industrial Backbone", Country: "CN", Registry: "APNIC"},
	9808:   {ASN: 9808, ASName: "CMNET-GD (China Mobile)", Country: "CN", Registry: "APNIC"},
	56040:  {ASN: 56040, ASName: "China Mobile International", Country: "HK", Registry: "APNIC"},
	4538:   {ASN: 4538, ASName: "CERNET (China Education and Research Network)", Country: "CN", Registry: "APNIC"},
	7497:   {ASN: 7497, ASName: "CSTNET (China Science and Technology Network)", Country: "CN", Registry: "APNIC"},
	24424:  {ASN: 24424, ASName: "China Broadnet (CBN)", Country: "CN", Registry: "APNIC"},
	13335:  {ASN: 13335, ASName: "CLOUDFLARENET", Country: "US", Registry: "ARIN"},
	15169:  {ASN: 15169, ASName: "GOOGLE", Country: "US", Registry: "ARIN"},
	16509:  {ASN: 16509, ASName: "AMAZON-02", Country: "US", Registry: "ARIN"},
	8075:   {ASN: 8075, ASName: "MICROSOFT-CORP", Country: "US", Registry: "ARIN"},
	32934:  {ASN: 32934, ASName: "FACEBOOK", Country: "US", Registry: "ARIN"},
	20940:  {ASN: 20940, ASName: "AKAMAI-ASN1", Country: "NL", Registry: "RIPE"},
	54113:  {ASN: 54113, ASName: "FASTLY", Country: "US", Registry: "ARIN"},
	714:    {ASN: 714, ASName: "APPLE-ENGINEERING", Country: "US", Registry: "ARIN"},
	37963:  {ASN: 37963, ASName: "ALIBABA-CN-NET", Country: "CN", Registry: "APNIC"},
	45102:  {ASN: 45102, ASName: "ALIBABA-US-NET", Country: "US", Registry: "ARIN"},
	132203: {ASN: 132203, ASName: "TENCENT-NET-AP", Country: "CN", Registry: "APNIC"},
}

// NewASNService initializes ASN service.
func NewASNService(dataDir string) *ASNService {
	svc := &ASNService{
		cache:      make(map[string]cacheEntry),
		httpClient: ssrf.NewSafeHTTPClient(config.DefaultHTTPTimeout),
	}

	ipdbDir := filepath.Join(dataDir, "ipdb")
	maxmindASNPath := filepath.Join(ipdbDir, "GeoLite2-ASN.mmdb")
	if p, err := NewMaxMindASNProvider(maxmindASNPath); err == nil {
		svc.providers = append(svc.providers, p)
		slog.Info("Loaded MaxMind ASN database", "path", maxmindASNPath)
	} else {
		slog.Info("MaxMind ASN database not loaded (optional)", "path", maxmindASNPath, "err", err.Error())
	}

	return svc
}

// Reload reloads ASN providers from dataDir dynamically.
func (s *ASNService) Reload(dataDir string) {
	ipdbDir := filepath.Join(dataDir, "ipdb")
	maxmindASNPath := filepath.Join(ipdbDir, "GeoLite2-ASN.mmdb")
	var newProviders []ASNProvider

	if p, err := NewMaxMindASNProvider(maxmindASNPath); err == nil {
		newProviders = append(newProviders, p)
		slog.Info("Reloaded MaxMind ASN database", "path", maxmindASNPath)
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

// Close releases resources.
func (s *ASNService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range s.providers {
		_ = p.Close()
	}
}

// HasProviders checks if local ASN database is available.
func (s *ASNService) HasProviders() bool {
	return len(s.providers) > 0
}

// LookupIP queries ASN for an IP address.
func (s *ASNService) LookupIP(addr netip.Addr) *ASNResult {
	if !addr.IsValid() {
		return &ASNResult{Source: "none"}
	}

	key := "ip:" + addr.String()
	if cached := s.getFromCache(key); cached != nil {
		return cached
	}

	// Try local providers
	for _, p := range s.providers {
		if !p.Available() {
			continue
		}
		if res, err := p.Lookup(addr); err == nil && res != nil {
			s.saveToCache(key, res)
			return res
		}
	}

	// Fallback: try online lookup
	res := s.lookupIPOnline(addr)
	s.saveToCache(key, res)
	return res
}

// LookupQuery handles both ASN strings (AS4134, 4134) and IP addresses.
func (s *ASNService) LookupQuery(query string) (*ASNResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("empty query")
	}

	// Case 1: IP address
	if addr, err := netip.ParseAddr(query); err == nil {
		return s.LookupIP(addr), nil
	}

	// Case 2: ASN (AS4134 or 4134)
	cleanASN := strings.ToUpper(query)
	cleanASN = strings.TrimPrefix(cleanASN, "AS")
	asnNum, err := strconv.Atoi(cleanASN)
	if err != nil || asnNum <= 0 {
		return nil, fmt.Errorf("invalid ASN format: %s", query)
	}

	key := fmt.Sprintf("asn:%d", asnNum)
	if cached := s.getFromCache(key); cached != nil {
		return cached, nil
	}

	// Check well-known ASNs
	if wk, ok := wellKnownASNs[asnNum]; ok {
		res := wk
		res.Source = "builtin"
		s.saveToCache(key, &res)
		return &res, nil
	}

	// Fallback to RDAP online query
	res := s.lookupASNOnline(asnNum)
	s.saveToCache(key, res)
	return res, nil
}

func (s *ASNService) lookupIPOnline(addr netip.Addr) *ASNResult {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultHTTPTimeout)
	defer cancel()

	// 1. Try ip-api for instant ASN & ASName resolution
	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", addr.String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err == nil {
		resp, err := s.httpClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var data struct {
				Status      string `json:"status"`
				CountryCode string `json:"countryCode"`
				AS          string `json:"as"`
				ISP         string `json:"isp"`
				Org         string `json:"org"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&data); err == nil && data.Status == "success" {
				asnNum := 0
				asName := data.AS
				if strings.HasPrefix(strings.ToUpper(data.AS), "AS") {
					parts := strings.SplitN(data.AS, " ", 2)
					cleanASN := strings.TrimPrefix(strings.ToUpper(parts[0]), "AS")
					if n, err := strconv.Atoi(cleanASN); err == nil {
						asnNum = n
					}
					if len(parts) > 1 {
						asName = parts[1]
					}
				}
				if asName == "" {
					asName = data.ISP
				}
				return &ASNResult{
					ASN:     asnNum,
					ASName:  asName,
					Country: data.CountryCode,
					Source:  "ip-api",
				}
			}
		}
	}

	// 2. Fallback to RDAP IP endpoint
	urlRDAP := fmt.Sprintf("https://rdap.arin.net/registry/ip/%s", addr.String())
	reqRDAP, err := http.NewRequestWithContext(ctx, "GET", urlRDAP, nil)
	if err != nil {
		return &ASNResult{Source: "none"}
	}
	reqRDAP.Header.Set("Accept", "application/rdap+json, application/json")

	respRDAP, err := s.httpClient.Do(reqRDAP)
	if err != nil || respRDAP.StatusCode != http.StatusOK {
		return &ASNResult{Source: "none"}
	}
	defer respRDAP.Body.Close()

	var rdapData struct {
		Name      string `json:"name"`
		Handle    string `json:"handle"`
		Country   string `json:"country"`
		StartAddr string `json:"startAddress"`
		EndAddr   string `json:"endAddress"`
	}

	if err := json.NewDecoder(respRDAP.Body).Decode(&rdapData); err == nil {
		return &ASNResult{
			ASName:  rdapData.Name,
			Country: rdapData.Country,
			Network: fmt.Sprintf("%s - %s", rdapData.StartAddr, rdapData.EndAddr),
			Source:  "rdap",
		}
	}

	return &ASNResult{Source: "none"}
}

func (s *ASNService) lookupASNOnline(asnNum int) *ASNResult {
	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultHTTPTimeout)
	defer cancel()

	url := fmt.Sprintf("https://rdap.arin.net/registry/autnum/%d", asnNum)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return &ASNResult{ASN: asnNum, Source: "none"}
	}
	req.Header.Set("Accept", "application/rdap+json, application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return &ASNResult{
			ASN:    asnNum,
			ASName: fmt.Sprintf("AS%d", asnNum),
			Source: "none",
		}
	}
	defer resp.Body.Close()

	var data struct {
		Name    string `json:"name"`
		Handle  string `json:"handle"`
		Country string `json:"country"`
		Port43  string `json:"port43"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
		return &ASNResult{
			ASN:     asnNum,
			ASName:  data.Name,
			Country: data.Country,
			Source:  "rdap",
		}
	}

	return &ASNResult{
		ASN:    asnNum,
		ASName: fmt.Sprintf("AS%d", asnNum),
		Source: "none",
	}
}

func (s *ASNService) getFromCache(key string) *ASNResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if entry, ok := s.cache[key]; ok && time.Now().Before(entry.expiresAt) {
		return entry.result
	}
	return nil
}

func (s *ASNService) saveToCache(key string, res *ASNResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.cache) > 20000 {
		s.cache = make(map[string]cacheEntry)
	}
	s.cache[key] = cacheEntry{
		result:    res,
		expiresAt: time.Now().Add(1 * time.Hour),
	}
}
