package api

import (
	"context"
	"strings"

	"netip/internal/config"
	"netip/internal/ratelimit"
	"netip/internal/whois"

	"github.com/gin-gonic/gin"
)

// WHOISHandler handles WHOIS and RDAP queries.
type WHOISHandler struct {
	whoisService *whois.WHOISService
	semaphore    *ratelimit.Semaphore
}

// NewWHOISHandler creates a new WHOISHandler.
func NewWHOISHandler(whoisSvc *whois.WHOISService, sem *ratelimit.Semaphore) *WHOISHandler {
	return &WHOISHandler{
		whoisService: whoisSvc,
		semaphore:    sem,
	}
}

// WHOISRequest body.
type WHOISRequest struct {
	Target string `json:"target" binding:"required"`
}

// Query handles POST /api/v1/whois.
func (h *WHOISHandler) Query(c *gin.Context) {
	var req WHOISRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "请输入有效的域名或 IP 地址")
		return
	}

	req.Target = strings.TrimSpace(req.Target)
	if req.Target == "" {
		BadRequest(c, "EMPTY_TARGET", "目标地址不能为空")
		return
	}

	if !h.semaphore.TryAcquire() {
		RateLimitRejectionsTotal.WithLabelValues("concurrency_whois").Inc()
		ConcurrencyLimited(c, "当前 WHOIS 查询并发已达上限，请稍后重试")
		return
	}
	defer h.semaphore.Release()

	ActiveNetworkTests.WithLabelValues("whois").Inc()
	defer ActiveNetworkTests.WithLabelValues("whois").Dec()

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultWHOISTimeout)
	defer cancel()

	result, err := h.whoisService.Query(ctx, req.Target)
	if err != nil {
		BadRequest(c, "WHOIS_FAILED", err.Error())
		return
	}

	Success(c, result)
}
