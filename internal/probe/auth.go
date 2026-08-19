package probe

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	HeaderProbeID        = "X-Probe-ID"
	HeaderProbeTimestamp = "X-Probe-Timestamp"
	HeaderProbeSignature = "X-Probe-Signature"
	SignatureWindowSec   = 30
)

// GenerateSignature generates HMAC-SHA256 signature for probe RPC request.
// payload = timestamp + "\n" + method + "\n" + path + "\n" + sha256(body)
func GenerateSignature(secret string, timestamp int64, method, path string, body []byte) string {
	bodyHash := sha256.Sum256(body)
	bodyHashHex := hex.EncodeToString(bodyHash[:])

	message := fmt.Sprintf("%d\n%s\n%s\n%s", timestamp, strings.ToUpper(method), path, bodyHashHex)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifySignature validates that the request signature is authentic and within the time window.
func VerifySignature(secret string, probeID, timestampStr, signature, method, path string, body []byte) error {
	if secret == "" {
		return fmt.Errorf("probe secret not configured")
	}

	ts, err := strconv.ParseInt(strings.TrimSpace(timestampStr), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %w", err)
	}

	now := time.Now().Unix()
	if ts < now-SignatureWindowSec || ts > now+SignatureWindowSec {
		return fmt.Errorf("probe request timestamp expired (diff: %d s)", now-ts)
	}

	expectedSig := GenerateSignature(secret, ts, method, path, body)
	if !hmac.Equal([]byte(strings.ToLower(signature)), []byte(expectedSig)) {
		return fmt.Errorf("invalid probe signature")
	}

	return nil
}
