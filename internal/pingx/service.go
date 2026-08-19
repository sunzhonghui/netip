package pingx

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/netip"
	"os"
	"sync"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// PingRequest contains parameters for an ICMP ping.
type PingRequest struct {
	Target    string `json:"target"`
	Count     int    `json:"count"`
	IPVersion string `json:"ip_version"` // "auto", "4", "6"
}

// PingResult represents the aggregated results of a Ping.
type PingResult struct {
	Target      string    `json:"target"`
	ResolvedIP  string    `json:"resolved_ip"`
	Node        string    `json:"node"`
	Sent        int       `json:"sent"`
	Received    int       `json:"received"`
	LossPercent float64   `json:"loss_percent"`
	MinMs       float64   `json:"min_ms"`
	AvgMs       float64   `json:"avg_ms"`
	MaxMs       float64   `json:"max_ms"`
	Samples     []float64 `json:"samples"`
	Error       string    `json:"error,omitempty"`
}

// PingService manages ICMP ping operations.
type PingService struct {
	mu sync.Mutex
}

// NewPingService creates a new PingService.
func NewPingService() *PingService {
	return &PingService{}
}

// Ping executes ICMP pings against target.
func (s *PingService) Ping(ctx context.Context, req PingRequest) (*PingResult, error) {
	if req.Count <= 0 {
		req.Count = 4
	}
	if req.Count > config.DefaultPingMaxCount {
		req.Count = config.DefaultPingMaxCount
	}

	safeIPs, err := ssrf.ResolveSafeIPs(ctx, req.Target)
	if err != nil {
		return nil, err
	}

	// Select IP according to version preference
	var targetIP netip.Addr
	for _, ip := range safeIPs {
		if req.IPVersion == "4" && ip.Is4() {
			targetIP = ip
			break
		} else if req.IPVersion == "6" && ip.Is6() {
			targetIP = ip
			break
		}
	}

	if !targetIP.IsValid() {
		// Default to first IP (prefer IPv6 if dual-stack auto)
		for _, ip := range safeIPs {
			if ip.Is6() {
				targetIP = ip
				break
			}
		}
		if !targetIP.IsValid() {
			targetIP = safeIPs[0]
		}
	}

	res := &PingResult{
		Target:      req.Target,
		ResolvedIP:  targetIP.String(),
		Node:        "local",
		Sent:        req.Count,
		Received:    0,
		Samples:     make([]float64, 0, req.Count),
		LossPercent: 100.0,
	}

	isIPv6 := targetIP.Is6()
	network := "ip4:icmp"
	listenAddr := "0.0.0.0"
	var icmpType icmp.Type = ipv4.ICMPTypeEcho
	if isIPv6 {
		network = "ip6:ipv6-icmp"
		listenAddr = "::"
		icmpType = ipv6.ICMPTypeEchoRequest
	}

	// Listen for ICMP packets. First try raw socket, fallback to "udp" ping if unprivileged
	conn, err := icmp.ListenPacket(network, listenAddr)
	if err != nil {
		// Try unprivileged UDP ICMP
		udpNet := "udp4"
		if isIPv6 {
			udpNet = "udp6"
		}
		conn, err = icmp.ListenPacket(udpNet, listenAddr)
		if err != nil {
			return s.fallbackPing(ctx, req, targetIP)
		}
	}
	defer conn.Close()

	id := os.Getpid() & 0xffff
	var latencies []float64

	for seq := 1; seq <= req.Count; seq++ {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}

		payload := make([]byte, 32)
		rand.Read(payload)

		wm := icmp.Message{
			Type: icmpType,
			Code: 0,
			Body: &icmp.Echo{
				ID:   id,
				Seq:  seq,
				Data: payload,
			},
		}

		wb, err := wm.Marshal(nil)
		if err != nil {
			continue
		}

		var dstAddr net.Addr
		if conn.LocalAddr().Network() == "udp" || conn.LocalAddr().Network() == "udp4" || conn.LocalAddr().Network() == "udp6" {
			dstAddr = &net.UDPAddr{IP: net.ParseIP(targetIP.String()), Port: 0}
		} else {
			dstAddr = &net.IPAddr{IP: net.ParseIP(targetIP.String())}
		}

		start := time.Now()
		_ = conn.SetDeadline(time.Now().Add(config.DefaultPingPacketTimeout))

		if _, err := conn.WriteTo(wb, dstAddr); err != nil {
			continue
		}

		rb := make([]byte, 1500)
		for {
			n, _, err := conn.ReadFrom(rb)
			if err != nil {
				// Timeout or read error
				break
			}

			rtt := time.Since(start)
			protoNum := 1 // ICMPv4
			if isIPv6 {
				protoNum = 58 // ICMPv6
			}

			rm, err := icmp.ParseMessage(protoNum, rb[:n])
			if err != nil {
				continue
			}

			if echo, ok := rm.Body.(*icmp.Echo); ok {
				// Check echo sequence
				if echo.Seq == seq {
					rttMs := math.Round(float64(rtt.Microseconds())/100.0) / 10.0 // 1 decimal place ms
					latencies = append(latencies, rttMs)
					res.Received++
					break
				}
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	res.Samples = latencies
	if res.Sent > 0 {
		res.LossPercent = math.Round((float64(res.Sent-res.Received)/float64(res.Sent)*100.0)*10.0) / 10.0
	}

	if len(latencies) > 0 {
		min := latencies[0]
		max := latencies[0]
		sum := 0.0
		for _, l := range latencies {
			if l < min {
				min = l
			}
			if l > max {
				max = l
			}
			sum += l
		}
		res.MinMs = min
		res.MaxMs = max
		res.AvgMs = math.Round((sum/float64(len(latencies)))*10.0) / 10.0
	}

	return res, nil
}

// fallbackPing measures TCP connect latency if ICMP raw socket permissions are restricted.
func (s *PingService) fallbackPing(ctx context.Context, req PingRequest, targetIP netip.Addr) (*PingResult, error) {
	res := &PingResult{
		Target:      req.Target,
		ResolvedIP:  targetIP.String(),
		Node:        "local (tcp-fallback)",
		Sent:        req.Count,
		Received:    0,
		Samples:     make([]float64, 0, req.Count),
		LossPercent: 100.0,
	}

	port := 80
	if targetIP.Is6() {
		port = 443
	}

	var latencies []float64
	for seq := 1; seq <= req.Count; seq++ {
		start := time.Now()
		targetAddr := net.JoinHostPort(targetIP.String(), fmt.Sprintf("%d", port))

		var dialer net.Dialer
		dialer.Timeout = config.DefaultPingPacketTimeout

		conn, err := dialer.DialContext(ctx, "tcp", targetAddr)
		rtt := time.Since(start)
		if err == nil {
			conn.Close()
			rttMs := math.Round(float64(rtt.Microseconds())/100.0) / 10.0
			latencies = append(latencies, rttMs)
			res.Received++
		}
		time.Sleep(100 * time.Millisecond)
	}

	res.Samples = latencies
	if res.Sent > 0 {
		res.LossPercent = math.Round((float64(res.Sent-res.Received)/float64(res.Sent)*100.0)*10.0) / 10.0
	}

	if len(latencies) > 0 {
		min := latencies[0]
		max := latencies[0]
		sum := 0.0
		for _, l := range latencies {
			if l < min {
				min = l
			}
			if l > max {
				max = l
			}
			sum += l
		}
		res.MinMs = min
		res.MaxMs = max
		res.AvgMs = math.Round((sum/float64(len(latencies)))*10.0) / 10.0
	}

	return res, nil
}
