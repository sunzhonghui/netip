package api

import (
	"context"
	"strings"

	"netip/internal/config"
	"netip/internal/ratelimit"
	"netip/internal/speedtest"

	"github.com/gin-gonic/gin"
)

// SpeedHandler handles HTTP speed tests.
type SpeedHandler struct {
	speedService *speedtest.SpeedService
	semaphore    *ratelimit.Semaphore
}

// NewSpeedHandler creates a SpeedHandler.
func NewSpeedHandler(speedSvc *speedtest.SpeedService, sem *ratelimit.Semaphore) *SpeedHandler {
	return &SpeedHandler{
		speedService: speedSvc,
		semaphore:    sem,
	}
}

// SpeedRequest body.
type SpeedRequest struct {
	Target string `json:"target" binding:"required"`
}

// Test handles POST /api/v1/speed.
func (h *SpeedHandler) Test(c *gin.Context) {
	var req SpeedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "请输入有效的测试目标网址")
		return
	}

	req.Target = strings.TrimSpace(req.Target)
	if req.Target == "" {
		BadRequest(c, "EMPTY_TARGET", "测试目标不能为空")
		return
	}

	if !h.semaphore.TryAcquire() {
		RateLimitRejectionsTotal.WithLabelValues("concurrency_speed").Inc()
		ConcurrencyLimited(c, "当前测速并发已达上限，请稍后重试")
		return
	}
	defer h.semaphore.Release()

	ActiveNetworkTests.WithLabelValues("speed").Inc()
	defer ActiveNetworkTests.WithLabelValues("speed").Dec()

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultSpeedTestOverallTimeout)
	defer cancel()

	result, err := h.speedService.Test(ctx, req.Target)
	if err != nil {
		BadRequest(c, "SPEED_TEST_FAILED", err.Error())
		return
	}

	if result.DownloadBytes > 0 {
		SpeedTestBytesTotal.Add(float64(result.DownloadBytes))
	}

	Success(c, result)
}
