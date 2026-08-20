package config

import (
	"net/netip"
	"os"
	"strconv"
	"strings"
	"time"
)

// AppConfig contains all runtime configurations for NetIP.
type AppConfig struct {
	// Mode: "server" or "probe"
	Mode string

	// Server settings
	HTTPAddr        string
	PublicDomain    string
	PublicAppOrigin string
	IPv4Domain      string
	IPv6Domain      string
	TestDomain      string

	// Trusted proxies for client IP extraction
	TrustedProxyCIDRs []netip.Prefix

	// Storage & IP databases
	DataDir string

	// Diagnostic Tool Settings
	DNSResolvers    []string
	AllowedTCPPorts map[int]bool

	// Observability
	LogLevel       string
	LogClientIP    bool
	MetricsEnabled bool

	// Probe Master Settings
	ProbeConfigFile string

	// Probe Node Settings
	ProbeID     string
	ProbeName   string
	ProbeISP    string
	ProbeSecret string

	// Rate Limits (requests per minute per IP)
	RateLimitGeneral   int
	RateLimitIP        int
	RateLimitDNS       int
	RateLimitPing      int
	RateLimitWeb       int
	RateLimitSpecialIP int

	// Global Concurrency Semaphores
	ConcurrencyPing  int
	ConcurrencyHTTP  int
	ConcurrencySpeed int
	ConcurrencyWHOIS int

	// IP Database Auto-Update
	IPDBAutoUpdate     bool
	IPDBUpdateInterval time.Duration
	MaxMindLicenseKey  string
	IP2RegionURL       string

	// Build info
	Version   string
	Commit    string
	BuildTime string
}

// Global build variables injected by -ldflags
var (
	Version   = "0.1.0"
	Commit    = "dev"
	BuildTime = "unknown"
)

// LoadConfig loads configuration from environment variables with sensible defaults.
func LoadConfig() *AppConfig {
	cfg := &AppConfig{
		Mode:            getEnv("MODE", "server"),
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		PublicDomain:    getEnv("PUBLIC_DOMAIN", "ip.ipw.3x.cx"),
		PublicAppOrigin: getEnv("PUBLIC_APP_ORIGIN", "https://ip.ipw.3x.cx"),
		IPv4Domain:      getEnv("IPV4_DOMAIN", "4.ipw.3x.cx"),
		IPv6Domain:      getEnv("IPV6_DOMAIN", "6.ipw.3x.cx"),
		TestDomain:      getEnv("TEST_DOMAIN", "test.ipw.3x.cx"),
		DataDir:         getEnv("DATA_DIR", "./data"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		LogClientIP:     getEnvBool("LOG_CLIENT_IP", true),
		MetricsEnabled:  getEnvBool("METRICS_ENABLED", true),
		ProbeConfigFile: getEnv("PROBE_CONFIG", "./config/probes.yaml"),

		// Probe Node credentials
		ProbeID:     getEnv("PROBE_ID", ""),
		ProbeName:   getEnv("PROBE_NAME", ""),
		ProbeISP:    getEnv("PROBE_ISP", ""),
		ProbeSecret: getEnv("PROBE_SECRET", ""),

		// Rate limits
		RateLimitGeneral:   getEnvInt("RATE_LIMIT_GENERAL", 60),
		RateLimitIP:        getEnvInt("RATE_LIMIT_IP", 30),
		RateLimitDNS:       getEnvInt("RATE_LIMIT_DNS", 20),
		RateLimitPing:      getEnvInt("RATE_LIMIT_PING", 10),
		RateLimitWeb:       getEnvInt("RATE_LIMIT_WEB", 10),
		RateLimitSpecialIP: getEnvInt("RATE_LIMIT_SPECIAL_IP", 120),

		// Concurrency limits
		ConcurrencyPing:  getEnvInt("CONCURRENCY_PING", 50),
		ConcurrencyHTTP:  getEnvInt("CONCURRENCY_HTTP", 30),
		ConcurrencySpeed: getEnvInt("CONCURRENCY_SPEED", 10),
		ConcurrencyWHOIS: getEnvInt("CONCURRENCY_WHOIS", 10),

		// IP Database Auto-Update
		IPDBAutoUpdate:     getEnvBool("IPDB_AUTO_UPDATE", true),
		IPDBUpdateInterval: getEnvDuration("IPDB_UPDATE_INTERVAL", 7*24*time.Hour),
		MaxMindLicenseKey:  getEnv("MAXMIND_LICENSE_KEY", ""),
		IP2RegionURL:       getEnv("IP2REGION_URL", ""),

		Version:   Version,
		Commit:    Commit,
		BuildTime: BuildTime,
	}

	// Parse Trusted Proxy CIDRs
	trustedStr := getEnv("TRUSTED_PROXY_CIDRS", "")
	cfg.TrustedProxyCIDRs = parseCIDRList(trustedStr)

	// Parse DNS Resolvers
	dnsStr := getEnv("DNS_RESOLVERS", "system,1.1.1.1,8.8.8.8,223.5.5.5,119.29.29.29")
	cfg.DNSResolvers = parseStringList(dnsStr)

	// Parse Allowed TCP Ports
	portsStr := getEnv("ALLOWED_TCP_PORTS", "22,53,80,443,465,587,993,995,3389")
	cfg.AllowedTCPPorts = parsePortMap(portsStr)

	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
			return n
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		val = strings.ToLower(strings.TrimSpace(val))
		return val == "true" || val == "1" || val == "yes"
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(strings.TrimSpace(val)); err == nil {
			return d
		}
	}
	return defaultVal
}

func parseCIDRList(raw string) []netip.Prefix {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var prefixes []netip.Prefix
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		prefix, err := netip.ParsePrefix(p)
		if err == nil {
			prefixes = append(prefixes, prefix)
		} else {
			// Try single IP as /32 or /128
			addr, err := netip.ParseAddr(p)
			if err == nil {
				bits := 32
				if addr.Is6() {
					bits = 128
				}
				prefixes = append(prefixes, netip.PrefixFrom(addr, bits))
			}
		}
	}
	return prefixes
}

func parseStringList(raw string) []string {
	var list []string
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			list = append(list, p)
		}
	}
	return list
}

func parsePortMap(raw string) map[int]bool {
	ports := make(map[int]bool)
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if port, err := strconv.Atoi(p); err == nil && port > 0 && port <= 65535 {
			ports[port] = true
		}
	}
	return ports
}
