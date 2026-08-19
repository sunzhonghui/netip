package probe

import (
	"strconv"
	"testing"
	"time"
)

func TestProbeSignatureComprehensive(t *testing.T) {
	secret := "my-secret-key-xyz"
	probeID := "cn-beijing-01"
	method := "POST"
	path := "/probe/v1/dns"
	body := []byte(`{"name":"example.com","type":"A"}`)

	now := time.Now().Unix()
	sig := GenerateSignature(secret, now, method, path, body)

	// 1. Valid
	err := VerifySignature(secret, probeID, strconv.FormatInt(now, 10), sig, method, path, body)
	if err != nil {
		t.Fatalf("expected valid signature to pass, got: %v", err)
	}

	// 2. Wrong secret
	err = VerifySignature("wrong-secret", probeID, strconv.FormatInt(now, 10), sig, method, path, body)
	if err == nil {
		t.Errorf("expected wrong secret to fail")
	}

	// 3. Expired timestamp (40s ago)
	expired := now - 40
	expiredSig := GenerateSignature(secret, expired, method, path, body)
	err = VerifySignature(secret, probeID, strconv.FormatInt(expired, 10), expiredSig, method, path, body)
	if err == nil {
		t.Errorf("expected expired timestamp to fail")
	}

	// 4. Future timestamp (40s in future)
	future := now + 40
	futureSig := GenerateSignature(secret, future, method, path, body)
	err = VerifySignature(secret, probeID, strconv.FormatInt(future, 10), futureSig, method, path, body)
	if err == nil {
		t.Errorf("expected future timestamp to fail")
	}

	// 5. Tampered body
	tamperedBody := []byte(`{"name":"example.com","type":"AAAA"}`)
	err = VerifySignature(secret, probeID, strconv.FormatInt(now, 10), sig, method, path, tamperedBody)
	if err == nil {
		t.Errorf("expected tampered body to fail")
	}
}
