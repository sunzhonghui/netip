package api

import (
	"net/http"

	"netip/internal/config"

	"github.com/gin-gonic/gin"
)

// SystemHandler handles health checks, readiness, and version info.
type SystemHandler struct {
	cfg            *config.AppConfig
	checkReadiness func() []string
}

// NewSystemHandler creates a new SystemHandler.
func NewSystemHandler(cfg *config.AppConfig, readinessChecker func() []string) *SystemHandler {
	return &SystemHandler{
		cfg:            cfg,
		checkReadiness: readinessChecker,
	}
}

// Healthz returns basic liveness status.
func (h *SystemHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Readyz returns readiness status and any non-fatal warnings (e.g. missing Geo DB).
func (h *SystemHandler) Readyz(c *gin.Context) {
	var warnings []string
	if h.checkReadiness != nil {
		warnings = h.checkReadiness()
	}

	resp := gin.H{
		"status": "ok",
	}
	if len(warnings) > 0 {
		resp["warnings"] = warnings
	}

	c.JSON(http.StatusOK, resp)
}

// Version returns build version metadata.
func (h *SystemHandler) Version(c *gin.Context) {
	Success(c, gin.H{
		"version":    h.cfg.Version,
		"commit":     h.cfg.Commit,
		"build_time": h.cfg.BuildTime,
		"mode":       h.cfg.Mode,
	})
}
