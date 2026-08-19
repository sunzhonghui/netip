package server

import (
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"netip/web"

	"github.com/gin-gonic/gin"
)

// ServeEmbeddedSPA serves embedded static assets from web.DistFS with SPA fallback to index.html.
func ServeEmbeddedSPA(r *gin.Engine) {
	distSubFS, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		return
	}

	fileServer := http.FileServer(http.FS(distSubFS))

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

		// Try opening the requested file in embedded FS
		if f, err := distSubFS.Open(path); err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// Fallback to index.html for Vue SPA client-side routes
		indexFile, err := distSubFS.Open("index.html")
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
