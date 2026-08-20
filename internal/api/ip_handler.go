package api

import (
	"net/netip"
	"strings"

	"netip/internal/asn"
	"netip/internal/clientip"
	"netip/internal/config"
	"netip/internal/ipgeo"

	"github.com/gin-gonic/gin"
)

// IPHandler handles IP, Geo and ASN lookup requests.
type IPHandler struct {
	cfg        *config.AppConfig
	geoService *ipgeo.GeoService
	asnService *asn.ASNService
}

// NewIPHandler creates a new IPHandler.
func NewIPHandler(cfg *config.AppConfig, geoSvc *ipgeo.GeoService, asnSvc *asn.ASNService) *IPHandler {
	return &IPHandler{
		cfg:        cfg,
		geoService: geoSvc,
		asnService: asnSvc,
	}
}

// IPDetailsResponse represents the response structure for IP queries.
type IPDetailsResponse struct {
	IP          string            `json:"ip"`
	Version     int               `json:"version"`
	Country     string            `json:"country"`
	CountryCode string            `json:"country_code"`
	Province    string            `json:"province"`
	City        string            `json:"city"`
	ISP         string            `json:"isp"`
	ASN         int               `json:"asn"`
	ASName      string            `json:"as_name"`
	Network     string            `json:"network"`
	Sources     map[string]string `json:"sources"`
}

// Me handles GET /api/v1/me (current client IP + Geo + ASN).
func (h *IPHandler) Me(c *gin.Context) {
	ip := clientip.ClientIP(c.Request, h.cfg.TrustedProxyCIDRs)
	if !ip.IsValid() {
		BadRequest(c, "INVALID_CLIENT_IP", "无法获取客户端公网 IP 地址")
		return
	}

	details := h.resolveIP(ip)
	Success(c, details)
}

// LookupIP handles GET /api/v1/ip/:ip.
func (h *IPHandler) LookupIP(c *gin.Context) {
	ipStr := strings.TrimSpace(c.Param("ip"))
	if ipStr == "" {
		BadRequest(c, "EMPTY_IP", "IP 地址不能为空")
		return
	}

	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		BadRequest(c, "INVALID_IP", "无效的 IP 地址格式")
		return
	}

	if addr.Is4In6() {
		addr = addr.Unmap()
	}

	details := h.resolveIP(addr)
	Success(c, details)
}

// LookupASN handles GET /api/v1/asn/:query.
func (h *IPHandler) LookupASN(c *gin.Context) {
	query := strings.TrimSpace(c.Param("query"))
	if query == "" {
		BadRequest(c, "EMPTY_QUERY", "ASN 或 IP 查询参数不能为空")
		return
	}

	res, err := h.asnService.LookupQuery(query)
	if err != nil {
		BadRequest(c, "INVALID_ASN_QUERY", err.Error())
		return
	}

	Success(c, res)
}

func (h *IPHandler) resolveIP(addr netip.Addr) *IPDetailsResponse {
	version := 4
	if addr.Is6() {
		version = 6
	}

	geoRes := h.geoService.Lookup(addr)
	asnRes := h.asnService.LookupIP(addr)

	sources := make(map[string]string)
	if geoRes != nil && geoRes.Source != "" {
		sources["geo"] = geoRes.Source
	}
	if asnRes != nil && asnRes.Source != "" {
		sources["asn"] = asnRes.Source
	}

	resp := &IPDetailsResponse{
		IP:      addr.String(),
		Version: version,
		Sources: sources,
	}

	if geoRes != nil {
		resp.Country = geoRes.Country
		resp.CountryCode = geoRes.CountryCode
		resp.Province = geoRes.Province
		resp.City = geoRes.City
		resp.ISP = geoRes.ISP
	}

	if asnRes != nil {
		resp.ASN = asnRes.ASN
		resp.ASName = asnRes.ASName
		resp.Network = asnRes.Network
		if resp.Country == "" && asnRes.Country != "" {
			resp.Country = asnRes.Country
			resp.CountryCode = asnRes.Country
		}
		if resp.ISP == "" && asnRes.ASName != "" {
			resp.ISP = asnRes.ASName
		}
	}

	return resp
}
