package api

import (
	"context"
	"strings"

	"netip/internal/config"
	"netip/internal/sslcheck"

	"github.com/gin-gonic/gin"
)

// SSLHandler handles SSL/TLS certificate diagnostic endpoints.
type SSLHandler struct {
	sslService *sslcheck.SSLService
}

// NewSSLHandler creates a new SSLHandler.
func NewSSLHandler(sslSvc *sslcheck.SSLService) *SSLHandler {
	return &SSLHandler{
		sslService: sslSvc,
	}
}

// Check handles POST /api/v1/ssl.
func (h *SSLHandler) Check(c *gin.Context) {
	var req sslcheck.SSLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "请输入有效的域名或主机名")
		return
	}

	req.Hostname = strings.TrimSpace(req.Hostname)
	if req.Hostname == "" {
		BadRequest(c, "EMPTY_HOSTNAME", "主机名不能为空")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultTLSDialTimeout)
	defer cancel()

	result, err := h.sslService.Inspect(ctx, req)
	if err != nil {
		BadRequest(c, "SSL_CHECK_FAILED", err.Error())
		return
	}

	Success(c, result)
}
