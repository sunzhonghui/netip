package ipdb

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"netip/internal/asn"
	"netip/internal/config"
	"netip/internal/ipgeo"
	"netip/internal/security/ssrf"
)

// Updater manages scheduled automatic downloads and hot-reloads of IP databases.
type Updater struct {
	cfg        *config.AppConfig
	geoSvc     *ipgeo.GeoService
	asnSvc     *asn.ASNService
	httpClient *http.Client
}

// NewUpdater creates a new IP database updater.
func NewUpdater(cfg *config.AppConfig, geoSvc *ipgeo.GeoService, asnSvc *asn.ASNService) *Updater {
	return &Updater{
		cfg:        cfg,
		geoSvc:     geoSvc,
		asnSvc:     asnSvc,
		httpClient: ssrf.NewSafeHTTPClient(3 * time.Minute),
	}
}

// Start launches the background updater goroutine if enabled.
func (u *Updater) Start(ctx context.Context) {
	if !u.cfg.IPDBAutoUpdate {
		slog.Info("IP database auto-update is disabled (IPDB_AUTO_UPDATE=false)")
		return
	}

	ipdbDir := filepath.Join(u.cfg.DataDir, "ipdb")
	if err := os.MkdirAll(ipdbDir, 0755); err != nil {
		slog.Warn("Failed to create ipdb directory", "dir", ipdbDir, "err", err.Error())
	}

	go func() {
		slog.Info("IP database auto-updater started", "interval", u.cfg.IPDBUpdateInterval)

		// Initial check: if files don't exist, download immediately
		u.checkAndDownloadAll(ctx)

		ticker := time.NewTicker(u.cfg.IPDBUpdateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				slog.Info("IP database auto-updater stopped")
				return
			case <-ticker.C:
				slog.Info("Running scheduled IP database update...")
				u.checkAndDownloadAll(ctx)
			}
		}
	}()
}

func (u *Updater) checkAndDownloadAll(ctx context.Context) {
	ipdbDir := filepath.Join(u.cfg.DataDir, "ipdb")
	updatedAny := false

	// 1. ip2region.xdb
	ip2regionPath := filepath.Join(ipdbDir, "ip2region.xdb")
	if u.shouldDownload(ip2regionPath) {
		urls := []string{
			"https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb",
			"https://fastly.jsdelivr.net/gh/lionsoul2014/ip2region@master/data/ip2region.xdb",
			"https://raw.gitmirror.com/lionsoul2014/ip2region/master/data/ip2region.xdb",
		}
		if u.cfg.IP2RegionURL != "" {
			urls = append([]string{u.cfg.IP2RegionURL}, urls...)
		}

		if err := u.downloadWithMirrors(ctx, urls, ip2regionPath, 1024*1024); err == nil {
			slog.Info("Successfully updated ip2region.xdb", "path", ip2regionPath)
			updatedAny = true
		} else {
			slog.Warn("Failed to download ip2region.xdb", "err", err.Error())
		}
	}

	// 2. GeoLite2-ASN.mmdb
	asnPath := filepath.Join(ipdbDir, "GeoLite2-ASN.mmdb")
	if u.shouldDownload(asnPath) {
		asnURLs := []string{
			"https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-ASN.mmdb",
			"https://git.io/GeoLite2-ASN.mmdb",
		}
		if err := u.downloadWithMirrors(ctx, asnURLs, asnPath, 2*1024*1024); err == nil {
			slog.Info("Successfully updated GeoLite2-ASN.mmdb", "path", asnPath)
			updatedAny = true
		} else {
			slog.Info("GeoLite2-ASN.mmdb optional download skipped/failed", "err", err.Error())
		}
	}

	// 3. GeoLite2-City.mmdb
	cityPath := filepath.Join(ipdbDir, "GeoLite2-City.mmdb")
	if u.shouldDownload(cityPath) {
		cityURLs := []string{
			"https://raw.githubusercontent.com/P3TERX/GeoLite.mmdb/download/GeoLite2-City.mmdb",
			"https://git.io/GeoLite2-City.mmdb",
		}
		if err := u.downloadWithMirrors(ctx, cityURLs, cityPath, 10*1024*1024); err == nil {
			slog.Info("Successfully updated GeoLite2-City.mmdb", "path", cityPath)
			updatedAny = true
		} else {
			slog.Info("GeoLite2-City.mmdb optional download skipped/failed", "err", err.Error())
		}
	}

	// Hot reload services if any database was downloaded/updated
	if updatedAny {
		slog.Info("Hot-reloading GeoService and ASNService with updated databases...")
		u.geoSvc.Reload(u.cfg.DataDir)
		u.asnSvc.Reload(u.cfg.DataDir)
	}
}

func (u *Updater) shouldDownload(targetPath string) bool {
	info, err := os.Stat(targetPath)
	if err != nil {
		// File does not exist
		return true
	}

	// If file is older than update interval or smaller than 100KB (broken)
	if info.Size() < 100*1024 {
		return true
	}

	if time.Since(info.ModTime()) > u.cfg.IPDBUpdateInterval {
		return true
	}

	return false
}

func (u *Updater) downloadWithMirrors(ctx context.Context, urls []string, destPath string, minBytes int64) error {
	var lastErr error

	for _, downloadURL := range urls {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := u.downloadFile(ctx, downloadURL, destPath, minBytes)
		if err == nil {
			return nil
		}
		lastErr = err
		slog.Debug("Mirror download failed, trying next mirror", "url", downloadURL, "err", err.Error())
	}

	return fmt.Errorf("all mirrors failed: %w", lastErr)
}

func (u *Updater) downloadFile(ctx context.Context, url, destPath string, minBytes int64) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "NetIP-Updater/1.0 (https://github.com/sunzhonghui/netip)")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	tmpFile := destPath + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
		_ = os.Remove(tmpFile) // clean up if still exists
	}()

	written, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	if written < minBytes {
		return fmt.Errorf("file too small: got %d bytes, expected >= %d", written, minBytes)
	}

	_ = out.Close()

	// Atomic replace
	return os.Rename(tmpFile, destPath)
}
