package api

import (
	"context"
	"strings"

	"netip/internal/config"
	"netip/internal/ratelimit"
	"netip/internal/webcheck"

	"github.com/gin-gonic/gin"
)

// WebHandler handles website IPv6 and HTTP diagnostics.
type WebHandler struct {
	ipv6Checker *webcheck.IPv6Checker
	httpChecker *webcheck.HTTPChecker
	semaphore   *ratelimit.Semaphore
}

// NewWebHandler creates a WebHandler.
func NewWebHandler(
	ipv6Checker *webcheck.IPv6Checker,
	httpChecker *webcheck.HTTPChecker,
	sem *ratelimit.Semaphore,
) *WebHandler {
	return &WebHandler{
		ipv6Checker: ipv6Checker,
		httpChecker: httpChecker,
		semaphore:   sem,
	}
}

// WebTargetRequest body.
type WebTargetRequest struct {
	Target string `json:"target" binding:"required"`
}

// IPv6Check handles POST /api/v1/ipv6-check.
func (h *WebHandler) IPv6Check(c *gin.Context) {
	var req WebTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "请输入有效的目标网站域名")
		return
	}

	req.Target = strings.TrimSpace(req.Target)
	if req.Target == "" {
		BadRequest(c, "EMPTY_TARGET", "目标网站不能为空")
		return
	}

	if !h.semaphore.TryAcquire() {
		RateLimitRejectionsTotal.WithLabelValues("concurrency_http").Inc()
		ConcurrencyLimited(c, "当前检测并发已达上限，请稍后重试")
		return
	}
	defer h.semaphore.Release()

	ActiveNetworkTests.WithLabelValues("ipv6_check").Inc()
	defer ActiveNetworkTests.WithLabelValues("ipv6_check").Dec()

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultIPv6CheckTimeout)
	defer cancel()

	result, err := h.ipv6Checker.Check(ctx, req.Target)
	if err != nil {
		BadRequest(c, "IPV6_CHECK_FAILED", err.Error())
		return
	}

	Success(c, result)
}

// HTTPCheck handles POST /api/v1/http.
func (h *WebHandler) HTTPCheck(c *gin.Context) {
	var req WebTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "请输入有效的目标网址")
		return
	}

	req.Target = strings.TrimSpace(req.Target)
	if req.Target == "" {
		BadRequest(c, "EMPTY_TARGET", "目标网址不能为空")
		return
	}

	if !h.semaphore.TryAcquire() {
		RateLimitRejectionsTotal.WithLabelValues("concurrency_http").Inc()
		ConcurrencyLimited(c, "当前检测并发已达上限，请稍后重试")
		return
	}
	defer h.semaphore.Release()

	ActiveNetworkTests.WithLabelValues("http_check").Inc()
	defer ActiveNetworkTests.WithLabelValues("http_check").Dec()

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultHTTPTimeout)
	defer cancel()

	result, err := h.httpChecker.Check(ctx, req.Target)
	if err != nil {
		BadRequest(c, "HTTP_CHECK_FAILED", err.Error())
		return
	}

	Success(c, result)
}
