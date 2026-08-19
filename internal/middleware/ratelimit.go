package middleware

import (
	"netip/internal/api"
	"netip/internal/clientip"
	"netip/internal/config"
	"netip/internal/ratelimit"

	"github.com/gin-gonic/gin"
)

// RateLimit creates a Gin middleware that enforces an IP-based token bucket rate limit.
func RateLimit(cfg *config.AppConfig, limiter *ratelimit.IPRateLimiter, metricLabel string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter == nil {
			c.Next()
			return
		}

		ip := clientip.ClientIP(c.Request, cfg.TrustedProxyCIDRs)
		ipStr := ip.String()
		if !ip.IsValid() {
			ipStr = c.ClientIP()
		}

		c.Set("ClientIP", ipStr)

		if !limiter.Allow(ipStr) {
			api.RateLimitRejectionsTotal.WithLabelValues(metricLabel).Inc()
			c.Header("Retry-After", "60")
			api.RateLimited(c, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
