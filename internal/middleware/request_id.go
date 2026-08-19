package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderXRequestID = "X-Request-ID"

// RequestID middleware ensures every request has a unique Request ID.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set("RequestID", reqID)
		c.Header(HeaderXRequestID, reqID)
		c.Next()
	}
}
