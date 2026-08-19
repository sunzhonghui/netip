package api

import (
	"context"
	"strings"
	"sync"
	"time"

	"netip/internal/config"
	"netip/internal/dnsx"
	"netip/internal/probe"

	"github.com/gin-gonic/gin"
)

// DNSHandler handles DNS diagnostic requests.
type DNSHandler struct {
	dnsService   *dnsx.DNSService
	probeManager *probe.Manager
}

// NewDNSHandler creates a DNSHandler.
func NewDNSHandler(dnsSvc *dnsx.DNSService, probeMgr *probe.Manager) *DNSHandler {
	return &DNSHandler{
		dnsService:   dnsSvc,
		probeManager: probeMgr,
	}
}

// DNSRequest represents the JSON body for /api/v1/dns.
type DNSRequest struct {
	Name string `json:"name" binding:"required"`
	Type string `json:"type"`
}

// Query handles POST /api/v1/dns.
func (h *DNSHandler) Query(c *gin.Context) {
	var req DNSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "请输入有效的域名和查询类型")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Type == "" {
		req.Type = "A"
	}
	req.Type = strings.ToUpper(strings.TrimSpace(req.Type))

	DNSQueriesTotal.WithLabelValues(req.Type, "all").Inc()

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultDNSOverallTimeout)
	defer cancel()

	// 1. Local resolution
	localRes, err := h.dnsService.Query(ctx, req.Name, req.Type)
	if err != nil {
		BadRequest(c, "DNS_QUERY_FAILED", err.Error())
		return
	}

	// 2. Dispatch to Probes in parallel if available
	probeNodes := h.probeManager.Nodes()
	if len(probeNodes) > 0 {
		var wg sync.WaitGroup
		var probeResults []dnsx.ResolverResult
		var mu sync.Mutex

		for _, node := range probeNodes {
			wg.Add(1)
			go func(n probe.NodeConfig) {
				defer wg.Done()
				start := time.Now()
				var pResp struct {
					Results []dnsx.ResolverResult `json:"results"`
				}
				err := h.probeManager.Client().Call(ctx, n, "/probe/v1/dns", req, &pResp)
				latency := time.Since(start).Milliseconds()

				mu.Lock()
				defer mu.Unlock()
				if err == nil && len(pResp.Results) > 0 {
					for _, r := range pResp.Results {
						r.NodeId = n.ID
						r.NodeName = n.Name
						r.ISP = n.ISP
						probeResults = append(probeResults, r)
					}
				} else {
					errStr := "probe timeout or unreachable"
					if err != nil {
						errStr = err.Error()
					}
					probeResults = append(probeResults, dnsx.ResolverResult{
						NodeId:    n.ID,
						NodeName:  n.Name,
						ISP:       n.ISP,
						Resolver:  n.Name,
						LatencyMs: latency,
						Error:     errStr,
					})
				}
			}(node)
		}

		wg.Wait()
		localRes.Results = append(localRes.Results, probeResults...)
	}

	Success(c, localRes)
}
