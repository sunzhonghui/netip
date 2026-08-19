package api

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests processed by NetIP.",
		},
		[]string{"method", "path", "status"},
	)

	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of HTTP request durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	ActiveNetworkTests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_network_tests",
			Help: "Current number of concurrent active network diagnostics.",
		},
		[]string{"tool"},
	)

	DNSQueriesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dns_queries_total",
			Help: "Total number of DNS queries performed.",
		},
		[]string{"type", "resolver"},
	)

	PingRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ping_requests_total",
			Help: "Total number of Ping requests.",
		},
		[]string{"protocol"},
	)

	SpeedTestBytesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "speedtest_bytes_total",
			Help: "Total number of bytes transferred during speed tests.",
		},
	)

	RateLimitRejectionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_rejections_total",
			Help: "Total number of requests rejected by rate limiting or concurrency limits.",
		},
		[]string{"type"},
	)
)
