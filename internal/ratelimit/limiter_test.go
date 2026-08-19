package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestIPRateLimiter_Allow(t *testing.T) {
	limiter := NewIPRateLimiter(2) // 2 requests per minute = 2 capacity
	defer limiter.Stop()

	ip := "192.0.2.1"

	// First two should pass
	if !limiter.Allow(ip) {
		t.Errorf("expected first request to be allowed")
	}
	if !limiter.Allow(ip) {
		t.Errorf("expected second request to be allowed")
	}

	// Third should be rejected
	if limiter.Allow(ip) {
		t.Errorf("expected third request to be rejected due to rate limit")
	}

	// Different IP should still be allowed
	if !limiter.Allow("192.0.2.2") {
		t.Errorf("expected different IP to be allowed")
	}
}

func TestSemaphore_Concurrency(t *testing.T) {
	sem := NewSemaphore(2)

	if !sem.TryAcquire() {
		t.Fatalf("failed to acquire first slot")
	}
	if !sem.TryAcquire() {
		t.Fatalf("failed to acquire second slot")
	}

	// Third attempt must fail
	if sem.TryAcquire() {
		t.Fatalf("expected third acquire to fail")
	}

	// Release one
	sem.Release()

	// Now try acquire should succeed
	if !sem.TryAcquire() {
		t.Fatalf("expected acquire after release to succeed")
	}

	sem.Release()
	sem.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	if err := sem.Acquire(ctx); err != nil {
		t.Fatalf("acquire failed: %v", err)
	}
	sem.Release()
}
