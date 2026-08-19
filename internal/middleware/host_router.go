package middleware

import (
	"net"
	"net/http"
	"strings"

	"netip/internal/clientip"
	"netip/internal/config"
	"netip/internal/ratelimit"

	"github.com/gin-gonic/gin"
)

// HostRouter handles special plain text & JSON IP API requests for subdomains:
// 4.ipw.3x.cx, 6.ipw.3x.cx, test.ipw.3x.cx.
func HostRouter(cfg *config.AppConfig, limiter *ratelimit.IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		host = strings.ToLower(strings.TrimSpace(host))

		isSpecialHost := false
		if (cfg.IPv4Domain != "" && host == strings.ToLower(cfg.IPv4Domain)) ||
			(cfg.IPv6Domain != "" && host == strings.ToLower(cfg.IPv6Domain)) ||
			(cfg.TestDomain != "" && host == strings.ToLower(cfg.TestDomain)) {
			isSpecialHost = true
		}

		path := c.Request.URL.Path

		if isSpecialHost || path == "/ip" || path == "/json" {
			// Extract client IP
			clientIP := clientip.ClientIP(c.Request, cfg.TrustedProxyCIDRs)
			if !clientIP.IsValid() {
				c.String(http.StatusBadRequest, "Invalid or unresolvable client IP\n")
				c.Abort()
				return
			}

			// Store client IP in gin context for logger
			c.Set("ClientIP", clientIP.String())

			// Rate limit check
			if limiter != nil && !limiter.Allow(clientIP.String()) {
				c.Header("Retry-After", "60")
				c.String(http.StatusTooManyRequests, "Rate limit exceeded\n")
				c.Abort()
				return
			}

			// Set CORS for frontend detection
			origin := c.GetHeader("Origin")
			if cfg.PublicAppOrigin != "" {
				if origin == cfg.PublicAppOrigin || cfg.PublicAppOrigin == "*" {
					c.Header("Access-Control-Allow-Origin", origin)
				} else {
					c.Header("Access-Control-Allow-Origin", cfg.PublicAppOrigin)
				}
			} else {
				c.Header("Access-Control-Allow-Origin", "*")
			}
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, HEAD")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Origin, Accept")

			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}

			// Path matching
			if path == "/json" {
				version := 4
				if clientIP.Is6() {
					version = 6
				}
				c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
				c.JSON(http.StatusOK, gin.H{
					"ip":      clientIP.String(),
					"version": version,
				})
				c.Abort()
				return
			}

			if path == "/" || path == "/ip" {
				c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
				c.Header("Content-Type", "text/plain; charset=utf-8")
				c.String(http.StatusOK, clientIP.String()+"\n")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
