package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"netip/internal/config"
	"netip/internal/middleware"
	"netip/internal/server"
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "Show version and exit")
	flag.BoolVar(&showVersion, "version", false, "Show version and exit")
	flag.Parse()

	if showVersion {
		fmt.Printf("NetIP v%s (commit: %s, built: %s)\n", config.Version, config.Commit, config.BuildTime)
		os.Exit(0)
	}

	cfg := config.LoadConfig()

	// Check CLI subcommand: e.g. "netip server" or "netip probe"
	args := flag.Args()
	if len(args) > 0 {
		subcmd := args[0]
		if subcmd == "probe" {
			cfg.Mode = "probe"
		} else if subcmd == "server" {
			cfg.Mode = "server"
		}
	}

	// Initialize structured logging
	middleware.InitLogger(cfg.LogLevel)

	srv := server.NewServer(cfg)
	if err := srv.Run(); err != nil {
		slog.Error("NetIP server failed", "err", err.Error())
		os.Exit(1)
	}
}
