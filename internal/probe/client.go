package probe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"netip/internal/config"
	"netip/internal/security/ssrf"
)

// NodeConfig represents a single remote probe node's credentials.
type NodeConfig struct {
	ID     string `yaml:"id" json:"id"`
	Name   string `yaml:"name" json:"name"`
	ISP    string `yaml:"isp" json:"isp"`
	URL    string `yaml:"url" json:"url"`
	Secret string `yaml:"secret" json:"-"`
}

// ProbeClient handles signed HTTP requests to a remote probe node.
type ProbeClient struct {
	httpClient *http.Client
}

// NewProbeClient creates a ProbeClient.
func NewProbeClient() *ProbeClient {
	return &ProbeClient{
		httpClient: ssrf.NewSafeHTTPClient(config.DefaultProbeRequestTimeout),
	}
}

// Call sends a signed request to a probe endpoint and unmarshals the response.
func (c *ProbeClient) Call(ctx context.Context, node NodeConfig, path string, reqBody interface{}, respData interface{}) error {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	targetURL := strings.TrimRight(node.URL, "/") + path
	req, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}

	ts := time.Now().Unix()
	sig := GenerateSignature(node.Secret, ts, "POST", path, bodyBytes)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(HeaderProbeID, node.ID)
	req.Header.Set(HeaderProbeTimestamp, strconv.FormatInt(ts, 10))
	req.Header.Set(HeaderProbeSignature, sig)
	req.Header.Set("User-Agent", ssrf.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("probe request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("probe node returned status %d: %s", resp.StatusCode, string(body))
	}

	var envelope struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
		Error   *struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("failed to decode probe response: %w", err)
	}

	if !envelope.Success {
		if envelope.Error != nil {
			return fmt.Errorf("probe error [%s]: %s", envelope.Error.Code, envelope.Error.Message)
		}
		return fmt.Errorf("probe returned unsuccessful response")
	}

	if err := json.Unmarshal(envelope.Data, respData); err != nil {
		return fmt.Errorf("failed to unmarshal probe data: %w", err)
	}

	return nil
}
