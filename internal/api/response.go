package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError details a structured error response.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// APIResponse represents the standard JSON envelope for NetIP APIs.
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
	RequestID string      `json:"request_id"`
}

// GetRequestID retrieves the X-Request-ID from gin context.
func GetRequestID(c *gin.Context) string {
	if reqID, ok := c.Get("RequestID"); ok {
		if s, ok := reqID.(string); ok && s != "" {
			return s
		}
	}
	return c.GetHeader("X-Request-ID")
}

// Success returns a 200 OK JSON response with the standard envelope.
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      data,
		RequestID: GetRequestID(c),
	})
}

// Error returns an error response with appropriate HTTP status code.
func Error(c *gin.Context, statusCode int, code string, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
		RequestID: GetRequestID(c),
	})
}

// Common error helpers
func BadRequest(c *gin.Context, code, message string) {
	Error(c, http.StatusBadRequest, code, message)
}

func Forbidden(c *gin.Context, code, message string) {
	Error(c, http.StatusForbidden, code, message)
}

func NotFound(c *gin.Context, code, message string) {
	Error(c, http.StatusNotFound, code, message)
}

func RateLimited(c *gin.Context, message string) {
	Error(c, http.StatusTooManyRequests, "RATE_LIMITED", message)
}

func ConcurrencyLimited(c *gin.Context, message string) {
	Error(c, http.StatusServiceUnavailable, "CONCURRENCY_LIMIT", message)
}

func ServerError(c *gin.Context, code, message string) {
	Error(c, http.StatusInternalServerError, code, message)
}

func UpstreamError(c *gin.Context, code, message string) {
	Error(c, http.StatusBadGateway, code, message)
}

func TimeoutError(c *gin.Context, code, message string) {
	Error(c, http.StatusRequestTimeout, code, message)
}
