package api

import (
	"context"
	"strings"
	"sync"

	"netip/internal/config"
	"netip/internal/pingx"
	"netip/internal/probe"
	"netip/internal/ratelimit"
	"netip/internal/tcping"

	"github.com/gin-gonic/gin"
)

// PingHandler handles ICMP and TCP ping operations.
type PingHandler struct {
	pingService   *pingx.PingService
	tcpingService *tcping.TCPingService
	probeManager  *probe.Manager
	semaphore     *ratelimit.Semaphore
}

// NewPingHandler creates a PingHandler.
func NewPingHandler(
	pingSvc *pingx.PingService,
	tcpingSvc *tcping.TCPingService,
	probeMgr *probe.Manager,
	sem *ratelimit.Semaphore,
) *PingHandler {
	return &PingHandler{
		pingService:   pingSvc,
		tcpingService: tcpingSvc,
		probeManager:  probeMgr,
		semaphore:     sem,
	}
}

// MultiNodePingResponse aggregates ping responses from multiple nodes.
type MultiNodePingResponse struct {
	Target      string             `json:"target"`
	ResolvedIP  string             `json:"resolved_ip"`
	Nodes       []pingx.PingResult `json:"nodes"`
}

// MultiNodeTCPingResponse aggregates tcping responses from multiple nodes.
type MultiNodeTCPingResponse struct {
	Target      string                 `json:"target"`
	Port        int                    `json:"port"`
	ResolvedIP  string                 `json:"resolved_ip"`
	Nodes       []tcping.TCPingResult `json:"nodes"`
}

// Ping handles POST /api/v1/ping.
func (h *PingHandler) Ping(c *gin.Context) {
	var req pingx.PingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "请输入有效的目标地址")
		return
	}

	req.Target = strings.TrimSpace(req.Target)
	if req.Target == "" {
		BadRequest(c, "EMPTY_TARGET", "目标地址不能为空")
		return
	}

	if !h.semaphore.TryAcquire() {
		RateLimitRejectionsTotal.WithLabelValues("concurrency_ping").Inc()
		ConcurrencyLimited(c, "当前 Ping 诊断并发已达上限，请稍后重试")
		return
	}
	defer h.semaphore.Release()

	ActiveNetworkTests.WithLabelValues("ping").Inc()
	defer ActiveNetworkTests.WithLabelValues("ping").Dec()
	PingRequestsTotal.WithLabelValues("icmp").Inc()

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultPingOverallTimeout)
	defer cancel()

	localRes, err := h.pingService.Ping(ctx, req)
	if err != nil {
		BadRequest(c, "PING_FAILED", err.Error())
		return
	}

	resp := &MultiNodePingResponse{
		Target:     req.Target,
		ResolvedIP: localRes.ResolvedIP,
		Nodes:      []pingx.PingResult{*localRes},
	}

	// Dispatch to Probes if configured
	probeNodes := h.probeManager.Nodes()
	if len(probeNodes) > 0 {
		var wg sync.WaitGroup
		var probeResults []pingx.PingResult
		var mu sync.Mutex

		for _, node := range probeNodes {
			wg.Add(1)
			go func(n probe.NodeConfig) {
				defer wg.Done()
				var pResp pingx.PingResult
				err := h.probeManager.Client().Call(ctx, n, "/probe/v1/ping", req, &pResp)

				mu.Lock()
				defer mu.Unlock()
				if err == nil {
					pResp.Node = n.Name + " " + n.ISP
					probeResults = append(probeResults, pResp)
				} else {
					probeResults = append(probeResults, pingx.PingResult{
						Target: req.Target,
						Node:   n.Name + " " + n.ISP,
						Error:  err.Error(),
					})
				}
			}(node)
		}

		wg.Wait()
		resp.Nodes = append(resp.Nodes, probeResults...)
	}

	Success(c, resp)
}

// TCPing handles POST /api/v1/tcping.
func (h *PingHandler) TCPing(c *gin.Context) {
	var req tcping.TCPingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "请输入有效的目标地址和端口")
		return
	}

	req.Target = strings.TrimSpace(req.Target)
	if req.Target == "" {
		BadRequest(c, "EMPTY_TARGET", "目标地址不能为空")
		return
	}

	if !h.semaphore.TryAcquire() {
		RateLimitRejectionsTotal.WithLabelValues("concurrency_tcping").Inc()
		ConcurrencyLimited(c, "当前 TCPing 诊断并发已达上限，请稍后重试")
		return
	}
	defer h.semaphore.Release()

	ActiveNetworkTests.WithLabelValues("tcping").Inc()
	defer ActiveNetworkTests.WithLabelValues("tcping").Dec()
	PingRequestsTotal.WithLabelValues("tcp").Inc()

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultTCPingOverallTimeout)
	defer cancel()

	localRes, err := h.tcpingService.Ping(ctx, req)
	if err != nil {
		BadRequest(c, "TCPING_FAILED", err.Error())
		return
	}

	resp := &MultiNodeTCPingResponse{
		Target:     req.Target,
		Port:       req.Port,
		ResolvedIP: localRes.ResolvedIP,
		Nodes:      []tcping.TCPingResult{*localRes},
	}

	// Dispatch to Probes if configured
	probeNodes := h.probeManager.Nodes()
	if len(probeNodes) > 0 {
		var wg sync.WaitGroup
		var probeResults []tcping.TCPingResult
		var mu sync.Mutex

		for _, node := range probeNodes {
			wg.Add(1)
			go func(n probe.NodeConfig) {
				defer wg.Done()
				var pResp tcping.TCPingResult
				err := h.probeManager.Client().Call(ctx, n, "/probe/v1/tcping", req, &pResp)

				mu.Lock()
				defer mu.Unlock()
				if err == nil {
					pResp.Node = n.Name + " " + n.ISP
					probeResults = append(probeResults, pResp)
				} else {
					probeResults = append(probeResults, tcping.TCPingResult{
						Target: req.Target,
						Port:   req.Port,
						Node:   n.Name + " " + n.ISP,
						Samples: []tcping.TCPingSample{
							{Success: false, Error: err.Error()},
						},
					})
				}
			}(node)
		}

		wg.Wait()
		resp.Nodes = append(resp.Nodes, probeResults...)
	}

	Success(c, resp)
}
