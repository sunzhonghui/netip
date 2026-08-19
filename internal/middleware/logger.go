package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// InitLogger initializes standard slog structured logger.
func InitLogger(levelStr string) *slog.Logger {
	var level slog.Level
	switch levelStr {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// StructuredLogger returns a Gin middleware that logs HTTP requests in JSON format.
func StructuredLogger(logClientIP bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		reqID := c.GetString("RequestID")

		if raw != "" {
			path = path + "?" + raw
		}

		attrs := []slog.Attr{
			slog.String("request_id", reqID),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Int64("latency_ms", latency.Milliseconds()),
		}

		if logClientIP {
			clientIP := c.GetString("ClientIP")
			if clientIP == "" {
				clientIP = c.ClientIP()
			}
			attrs = append(attrs, slog.String("client_ip", clientIP))
		}

		// Don't log healthz spam at error level
		if status >= 500 {
			slog.LogAttrs(c.Request.Context(), slog.LevelError, "HTTP Request", attrs...)
		} else if status >= 400 {
			slog.LogAttrs(c.Request.Context(), slog.LevelWarn, "HTTP Request", attrs...)
		} else {
			slog.LogAttrs(c.Request.Context(), slog.LevelInfo, "HTTP Request", attrs...)
		}
	}
}
