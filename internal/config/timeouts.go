package config

import "time"

// Centralized timeout constants across NetIP services.
// Handlers and services must not hardcode timeout durations.
const (
	// DNS Query Timeouts
	DefaultDNSSingleResolverTimeout = 3 * time.Second
	DefaultDNSOverallTimeout        = 5 * time.Second

	// Ping Timeouts
	DefaultPingPacketTimeout = 2 * time.Second
	DefaultPingOverallTimeout = 10 * time.Second
	DefaultPingMaxCount      = 5

	// TCPing Timeouts
	DefaultTCPingConnectTimeout = 2 * time.Second
	DefaultTCPingOverallTimeout = 8 * time.Second
	DefaultTCPingMaxCount       = 5

	// TLS / SSL Timeouts
	DefaultTLSDialTimeout = 8 * time.Second

	// HTTP & IPv6 Check Timeouts
	DefaultHTTPTimeout       = 10 * time.Second
	DefaultIPv6CheckTimeout  = 10 * time.Second
	DefaultMaxHTTPRedirects  = 3

	// Speed Test Timeouts and Limits
	DefaultSpeedTestOverallTimeout = 12 * time.Second
	DefaultSpeedTestMaxBytes       = 5 * 1024 * 1024 // 5 MB

	// WHOIS / RDAP Timeouts
	DefaultWHOISTimeout = 10 * time.Second

	// Remote Probe Timeouts
	DefaultProbeRequestTimeout = 12 * time.Second

	// Server Graceful Shutdown Timeout
	DefaultGracefulShutdownTimeout = 10 * time.Second

	// Client IP detection frontend timeout
	DefaultFrontendDetectionTimeout = 4 * time.Second
)
