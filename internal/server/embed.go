package server

import (
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"netip/web"

	"github.com/gin-gonic/gin"
)

// ServeEmbeddedSPA serves static assets from STATIC_DIR on disk if available, or falls back to embedded web.DistFS.
func ServeEmbeddedSPA(r *gin.Engine) {
	var targetFS fs.FS

	// 1. Check STATIC_DIR environment variable or mapped /web/dist volume directory
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		if stat, err := os.Stat("/web/dist"); err == nil && stat.IsDir() {
			staticDir = "/web/dist"
		} else if stat, err := os.Stat("./web/dist"); err == nil && stat.IsDir() {
			staticDir = "./web/dist"
		}
	}

	if staticDir != "" {
		if stat, err := os.Stat(staticDir); err == nil && stat.IsDir() {
			slog.Info("Serving frontend static assets from mapped host directory", "path", staticDir)
			targetFS = os.DirFS(staticDir)
		}
	}

	// 2. Fallback to embedded dist filesystem inside the Go binary
	if targetFS == nil {
		sub, err := fs.Sub(web.DistFS, "dist")
		if err == nil {
			slog.Info("Serving frontend static assets from embedded binary FS")
			targetFS = sub
		}
	}

	if targetFS == nil {
		slog.Error("No valid static asset filesystem found")
		return
	}

	fileServer := http.FileServer(http.FS(targetFS))

	r.NoRoute(func(c *gin.Context) {
		path := strings.TrimPrefix(c.Request.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Don't intercept API or probe routes
		if strings.HasPrefix(c.Request.URL.Path, "/api/") ||
			strings.HasPrefix(c.Request.URL.Path, "/probe/") ||
			c.Request.URL.Path == "/metrics" ||
			c.Request.URL.Path == "/healthz" ||
			c.Request.URL.Path == "/readyz" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "Endpoint not found",
				},
			})
			return
		}

		// Try opening requested file in target FS
		if f, err := targetFS.Open(path); err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// Fallback to index.html for SPA client-side routes
		indexFile, err := targetFS.Open("index.html")
		if err != nil {
			c.String(http.StatusNotFound, "Frontend static build not found")
			return
		}
		defer indexFile.Close()

		if rs, ok := indexFile.(io.ReadSeeker); ok {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Header("Cache-Control", "no-cache")
			http.ServeContent(c.Writer, c.Request, "index.html", time.Time{}, rs)
			return
		}

		content, err := io.ReadAll(indexFile)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read index.html")
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	})
}
