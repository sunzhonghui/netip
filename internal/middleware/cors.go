package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS handles Cross-Origin Resource Sharing.
func CORS(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if allowedOrigin != "" {
			if origin == allowedOrigin || allowedOrigin == "*" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				c.Header("Access-Control-Allow-Origin", allowedOrigin)
			}
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
