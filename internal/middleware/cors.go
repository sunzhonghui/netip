package middleware

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// isOriginAllowed checks if an origin is permitted by allowedOrigin.
// It allows exact matches, wildcard "*", or origins that share the same base domain.
func isOriginAllowed(origin, allowedOrigin string) bool {
	if origin == "" || allowedOrigin == "" {
		return false
	}
	if allowedOrigin == "*" || origin == allowedOrigin {
		return true
	}

	u1, err1 := url.Parse(origin)
	u2, err2 := url.Parse(allowedOrigin)
	if err1 != nil || err2 != nil {
		return false
	}

	h1 := strings.ToLower(u1.Hostname())
	h2 := strings.ToLower(u2.Hostname())
	if h1 == h2 {
		return true
	}

	// Extract base domain from allowedOrigin, stripping common prefix like "ip."
	baseDomain := h2
	if strings.HasPrefix(baseDomain, "ip.") {
		baseDomain = strings.TrimPrefix(baseDomain, "ip.")
	}

	// Allow the base domain itself or any subdomain of it
	return h1 == baseDomain || strings.HasSuffix(h1, "."+baseDomain)
}

// CORS handles Cross-Origin Resource Sharing.
func CORS(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Header("Vary", "Origin")

		if origin != "" && (allowedOrigin == "" || isOriginAllowed(origin, allowedOrigin)) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-Request-ID, X-Probe-ID, X-Probe-Timestamp, X-Probe-Signature, Origin, Accept")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
