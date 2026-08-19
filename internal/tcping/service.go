package tcping

import (
	"context"
	"fmt"
	"math"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"
)

// TCPingRequest parameters.
type TCPingRequest struct {
	Target string `json:"target"`
	Port   int    `json:"port"`
	Count  int    `json:"count"`
}

// TCPingSample represents an individual TCP connection attempt.
type TCPingSample struct {
	Success   bool    `json:"success"`
	LatencyMs float64 `json:"latency_ms,omitempty"`
	Error     string  `json:"error,omitempty"`
}

// TCPingResult represents the overall TCPing diagnostic result.
type TCPingResult struct {
	Target     string         `json:"target"`
	Port       int            `json:"port"`
	ResolvedIP string         `json:"resolved_ip"`
	Node       string         `json:"node"`
	Success    int            `json:"success"`
	Failed     int            `json:"failed"`
	AvgMs      float64        `json:"avg_ms"`
	MinMs      float64        `json:"min_ms"`
	MaxMs      float64        `json:"max_ms"`
	Samples    []TCPingSample `json:"samples"`
}

// TCPingService handles TCPing probes.
type TCPingService struct {
	allowedPorts map[int]bool
}

// NewTCPingService creates a TCPingService.
func NewTCPingService(allowedPorts map[int]bool) *TCPingService {
	return &TCPingService{
		allowedPorts: allowedPorts,
	}
}

// Ping executes TCP connect latency measurements.
func (s *TCPingService) Ping(ctx context.Context, req TCPingRequest) (*TCPingResult, error) {
	req.Target = strings.TrimSpace(req.Target)
	if req.Target == "" {
		return nil, fmt.Errorf("target cannot be empty")
	}

	if req.Port <= 0 {
		req.Port = 80
	}

	if err := ssrf.ValidatePort(req.Port, s.allowedPorts); err != nil {
		return nil, err
	}

	if req.Count <= 0 {
		req.Count = 4
	}
	if req.Count > config.DefaultTCPingMaxCount {
		req.Count = config.DefaultTCPingMaxCount
	}

	safeIPs, err := ssrf.ResolveSafeIPs(ctx, req.Target)
	if err != nil {
		return nil, err
	}

	targetIP := safeIPs[0]
	targetAddr := net.JoinHostPort(targetIP.String(), strconv.Itoa(req.Port))

	res := &TCPingResult{
		Target:     req.Target,
		Port:       req.Port,
		ResolvedIP: targetIP.String(),
		Node:       "local",
		Samples:    make([]TCPingSample, 0, req.Count),
	}

	var totalMs float64
	minMs := math.MaxFloat64
	maxMs := 0.0

	for i := 0; i < req.Count; i++ {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}

		sample := s.probeOnce(ctx, targetAddr, targetIP)
		res.Samples = append(res.Samples, sample)

		if sample.Success {
			res.Success++
			totalMs += sample.LatencyMs
			if sample.LatencyMs < minMs {
				minMs = sample.LatencyMs
			}
			if sample.LatencyMs > maxMs {
				maxMs = sample.LatencyMs
			}
		} else {
			res.Failed++
		}

		time.Sleep(100 * time.Millisecond)
	}

	if res.Success > 0 {
		res.AvgMs = math.Round((totalMs/float64(res.Success))*10.0) / 10.0
		res.MinMs = minMs
		res.MaxMs = maxMs
	}

	return res, nil
}

func (s *TCPingService) probeOnce(ctx context.Context, targetAddr string, _ netip.Addr) TCPingSample {
	start := time.Now()

	var dialer net.Dialer
	dialer.Timeout = config.DefaultTCPingConnectTimeout

	conn, err := dialer.DialContext(ctx, "tcp", targetAddr)
	latency := time.Since(start)

	if err != nil {
		return TCPingSample{
			Success: false,
			Error:   err.Error(),
		}
	}
	defer conn.Close()

	latencyMs := math.Round(float64(latency.Microseconds())/100.0) / 10.0
	return TCPingSample{
		Success:   true,
		LatencyMs: latencyMs,
	}
}
