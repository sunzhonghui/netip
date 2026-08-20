package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"netip/internal/api"
	"netip/internal/asn"
	"netip/internal/config"
	"netip/internal/dnsx"
	"netip/internal/ipdb"
	"netip/internal/ipgeo"
	"netip/internal/middleware"
	"netip/internal/pingx"
	"netip/internal/probe"
	"netip/internal/ratelimit"
	"netip/internal/sslcheck"
	"netip/internal/speedtest"
	"netip/internal/tcping"
	"netip/internal/webcheck"
	"netip/internal/whois"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server encapsulates the NetIP HTTP server.
type Server struct {
	cfg         *config.AppConfig
	engine      *gin.Engine
	httpServer  *http.Server
	geoSvc      *ipgeo.GeoService
	asnSvc      *asn.ASNService
	ipdbUpdater *ipdb.Updater
}

// NewServer initializes NetIP application services and routing.
func NewServer(cfg *config.AppConfig) *Server {
	if cfg.LogLevel != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	// Initialize core services
	geoSvc := ipgeo.NewGeoService(cfg.DataDir)
	asnSvc := asn.NewASNService(cfg.DataDir)
	dnsSvc := dnsx.NewDNSService(cfg.DNSResolvers)
	pingSvc := pingx.NewPingService()
	tcpingSvc := tcping.NewTCPingService(cfg.AllowedTCPPorts)
	sslSvc := sslcheck.NewSSLService(cfg.AllowedTCPPorts)
	whoisSvc := whois.NewWHOISService()
	ipv6Checker := webcheck.NewIPv6Checker()
	httpChecker := webcheck.NewHTTPChecker()
	speedSvc := speedtest.NewSpeedService()
	probeMgr := probe.NewManager(cfg.ProbeConfigFile)

	// Rate limiters
	generalLimiter := ratelimit.NewIPRateLimiter(cfg.RateLimitGeneral)
	ipLimiter := ratelimit.NewIPRateLimiter(cfg.RateLimitIP)
	dnsLimiter := ratelimit.NewIPRateLimiter(cfg.RateLimitDNS)
	pingLimiter := ratelimit.NewIPRateLimiter(cfg.RateLimitPing)
	webLimiter := ratelimit.NewIPRateLimiter(cfg.RateLimitWeb)
	specialIPLimiter := ratelimit.NewIPRateLimiter(cfg.RateLimitSpecialIP)

	// Concurrency semaphores
	pingSem := ratelimit.NewSemaphore(cfg.ConcurrencyPing)
	httpSem := ratelimit.NewSemaphore(cfg.ConcurrencyHTTP)
	speedSem := ratelimit.NewSemaphore(cfg.ConcurrencySpeed)
	whoisSem := ratelimit.NewSemaphore(cfg.ConcurrencyWHOIS)

	// Initialize Handlers
	systemHandler := api.NewSystemHandler(cfg, func() []string {
		var warnings []string
		if !geoSvc.HasProviders() {
			warnings = append(warnings, "geo database unavailable")
		}
		if !asnSvc.HasProviders() {
			warnings = append(warnings, "asn database unavailable")
		}
		return warnings
	})
	ipHandler := api.NewIPHandler(cfg, geoSvc, asnSvc)
	dnsHandler := api.NewDNSHandler(dnsSvc, probeMgr)
	pingHandler := api.NewPingHandler(pingSvc, tcpingSvc, probeMgr, pingSem)
	sslHandler := api.NewSSLHandler(sslSvc)
	whoisHandler := api.NewWHOISHandler(whoisSvc, whoisSem)
	webHandler := api.NewWebHandler(ipv6Checker, httpChecker, httpSem)
	speedHandler := api.NewSpeedHandler(speedSvc, speedSem)
	probeHandler := api.NewProbeHandler(cfg, dnsSvc, pingSvc, tcpingSvc, httpChecker)

	// Base middlewares
	engine.Use(middleware.RequestID())
	engine.Use(middleware.StructuredLogger(cfg.LogClientIP))
	engine.Use(middleware.SecurityHeaders())

	if cfg.Mode == "probe" {
		slog.Info("Starting in PROBE worker mode", "probe_id", cfg.ProbeID, "probe_name", cfg.ProbeName)

		engine.GET("/healthz", systemHandler.Healthz)

		probeGroup := engine.Group("/probe/v1")
		probeGroup.Use(probeHandler.AuthMiddleware())
		{
			probeGroup.POST("/dns", probeHandler.ProbeDNS)
			probeGroup.POST("/ping", probeHandler.ProbePing)
			probeGroup.POST("/tcping", probeHandler.ProbeTCPing)
			probeGroup.POST("/http", probeHandler.ProbeHTTP)
		}
	} else {
		// Server Mode
		slog.Info("Starting in SERVER master mode")

		// Special plain IP router on subdomains 4., 6., test.
		engine.Use(middleware.HostRouter(cfg, specialIPLimiter))
		engine.Use(middleware.CORS(cfg.PublicAppOrigin))

		// Observability
		engine.GET("/healthz", systemHandler.Healthz)
		engine.GET("/readyz", systemHandler.Readyz)
		if cfg.MetricsEnabled {
			engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
		}

		// API v1
		apiV1 := engine.Group("/api/v1")
		apiV1.Use(middleware.RateLimit(cfg, generalLimiter, "general"))
		{
			apiV1.GET("/version", systemHandler.Version)
			apiV1.GET("/me", middleware.RateLimit(cfg, ipLimiter, "ip"), ipHandler.Me)
			apiV1.GET("/ip/:ip", middleware.RateLimit(cfg, ipLimiter, "ip"), ipHandler.LookupIP)
			apiV1.GET("/asn/:query", middleware.RateLimit(cfg, ipLimiter, "ip"), ipHandler.LookupASN)

			apiV1.POST("/dns", middleware.RateLimit(cfg, dnsLimiter, "dns"), dnsHandler.Query)
			apiV1.POST("/ping", middleware.RateLimit(cfg, pingLimiter, "ping"), pingHandler.Ping)
			apiV1.POST("/tcping", middleware.RateLimit(cfg, pingLimiter, "tcping"), pingHandler.TCPing)
			apiV1.POST("/ipv6-check", middleware.RateLimit(cfg, webLimiter, "ipv6_check"), webHandler.IPv6Check)
			apiV1.POST("/http", middleware.RateLimit(cfg, webLimiter, "http"), webHandler.HTTPCheck)
			apiV1.POST("/ssl", middleware.RateLimit(cfg, webLimiter, "ssl"), sslHandler.Check)
			apiV1.POST("/whois", middleware.RateLimit(cfg, webLimiter, "whois"), whoisHandler.Query)
			apiV1.POST("/speed", middleware.RateLimit(cfg, webLimiter, "speed"), speedHandler.Test)
		}

		// Embedded Vue 3 Single Page Application
		ServeEmbeddedSPA(engine)
	}

	httpSrv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: engine,
	}

	ipdbUpdater := ipdb.NewUpdater(cfg, geoSvc, asnSvc)

	return &Server{
		cfg:         cfg,
		engine:      engine,
		httpServer:  httpSrv,
		geoSvc:      geoSvc,
		asnSvc:      asnSvc,
		ipdbUpdater: ipdbUpdater,
	}
}

// Run starts the HTTP server and handles graceful shutdown on SIGINT/SIGTERM.
func (s *Server) Run() error {
	slog.Info(fmt.Sprintf("NetIP server listening on %s", s.cfg.HTTPAddr),
		"mode", s.cfg.Mode,
		"public_domain", s.cfg.PublicDomain,
		"version", s.cfg.Version,
	)

	// Start IP database auto-updater in background
	updaterCtx, cancelUpdater := context.WithCancel(context.Background())
	defer cancelUpdater()
	s.ipdbUpdater.Start(updaterCtx)

	errCh := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-quit:
		slog.Info("Received shutdown signal", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.DefaultGracefulShutdownTimeout)
	defer cancel()

	slog.Info("Shutting down server gracefully...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("Server forced shutdown error", "err", err.Error())
	}

	s.geoSvc.Close()
	s.asnSvc.Close()
	slog.Info("NetIP server shutdown complete")
	return nil
}
