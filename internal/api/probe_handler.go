package api

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"netip/internal/config"
	"netip/internal/dnsx"
	"netip/internal/pingx"
	"netip/internal/probe"
	"netip/internal/tcping"
	"netip/internal/webcheck"

	"github.com/gin-gonic/gin"
)

// ProbeHandler handles RPC diagnostic requests on remote probe nodes.
type ProbeHandler struct {
	cfg           *config.AppConfig
	dnsService    *dnsx.DNSService
	pingService   *pingx.PingService
	tcpingService *tcping.TCPingService
	httpChecker   *webcheck.HTTPChecker
}

// NewProbeHandler creates a ProbeHandler.
func NewProbeHandler(
	cfg *config.AppConfig,
	dnsSvc *dnsx.DNSService,
	pingSvc *pingx.PingService,
	tcpingSvc *tcping.TCPingService,
	httpChecker *webcheck.HTTPChecker,
) *ProbeHandler {
	return &ProbeHandler{
		cfg:           cfg,
		dnsService:    dnsSvc,
		pingService:   pingSvc,
		tcpingService: tcpingSvc,
		httpChecker:   httpChecker,
	}
}

// AuthMiddleware verifies HMAC-SHA256 signature on probe requests.
func (h *ProbeHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.cfg.ProbeSecret == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "PROBE_UNCONFIGURED",
					"message": "Probe secret not configured on this node",
				},
			})
			return
		}

		probeID := c.GetHeader(probe.HeaderProbeID)
		timestamp := c.GetHeader(probe.HeaderProbeTimestamp)
		signature := c.GetHeader(probe.HeaderProbeSignature)

		if probeID == "" || timestamp == "" || signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "MISSING_AUTH_HEADERS",
					"message": "Missing probe authentication headers",
				},
			})
			return
		}

		// Read and preserve request body for verification and subsequent handler
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_BODY",
					"message": "Failed to read request body",
				},
			})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		path := c.Request.URL.Path
		method := c.Request.Method

		if err := probe.VerifySignature(h.cfg.ProbeSecret, probeID, timestamp, signature, method, path, bodyBytes); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "INVALID_PROBE_SIGNATURE",
					"message": err.Error(),
				},
			})
			return
		}

		c.Next()
	}
}

// ProbeDNS handles POST /probe/v1/dns.
func (h *ProbeHandler) ProbeDNS(c *gin.Context) {
	var req DNSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "Invalid DNS request")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultDNSOverallTimeout)
	defer cancel()

	res, err := h.dnsService.Query(ctx, req.Name, req.Type)
	if err != nil {
		BadRequest(c, "DNS_FAILED", err.Error())
		return
	}

	Success(c, res)
}

// ProbePing handles POST /probe/v1/ping.
func (h *ProbeHandler) ProbePing(c *gin.Context) {
	var req pingx.PingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "Invalid Ping request")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultPingOverallTimeout)
	defer cancel()

	res, err := h.pingService.Ping(ctx, req)
	if err != nil {
		BadRequest(c, "PING_FAILED", err.Error())
		return
	}

	Success(c, res)
}

// ProbeTCPing handles POST /probe/v1/tcping.
func (h *ProbeHandler) ProbeTCPing(c *gin.Context) {
	var req tcping.TCPingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "Invalid TCPing request")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultTCPingOverallTimeout)
	defer cancel()

	res, err := h.tcpingService.Ping(ctx, req)
	if err != nil {
		BadRequest(c, "TCPING_FAILED", err.Error())
		return
	}

	Success(c, res)
}

// ProbeHTTP handles POST /probe/v1/http.
func (h *ProbeHandler) ProbeHTTP(c *gin.Context) {
	var req WebTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "INVALID_INPUT", "Invalid HTTP request")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), config.DefaultHTTPTimeout)
	defer cancel()

	res, err := h.httpChecker.Check(ctx, req.Target)
	if err != nil {
		BadRequest(c, "HTTP_FAILED", err.Error())
		return
	}

	Success(c, res)
}
