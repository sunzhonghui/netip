package ratelimit

import (
	"sync"
	"time"
)

// tokenBucket tracks tokens for a single client.
type tokenBucket struct {
	tokens     float64
	capacity   float64
	refillRate float64 // tokens per second
	lastUpdate time.Time
}

// IPRateLimiter is an in-memory token bucket rate limiter indexed by IP.
type IPRateLimiter struct {
	mu           sync.Mutex
	buckets      map[string]*tokenBucket
	ratePerMin   int
	capacity     float64
	refillRate   float64
	cleanupEvery time.Duration
	stopCh       chan struct{}
}

// NewIPRateLimiter creates an in-memory token bucket rate limiter.
func NewIPRateLimiter(ratePerMinute int) *IPRateLimiter {
	if ratePerMinute <= 0 {
		ratePerMinute = 60
	}
	capacity := float64(ratePerMinute)
	refillRate := float64(ratePerMinute) / 60.0

	limiter := &IPRateLimiter{
		buckets:      make(map[string]*tokenBucket),
		ratePerMin:   ratePerMinute,
		capacity:     capacity,
		refillRate:   refillRate,
		cleanupEvery: 5 * time.Minute,
		stopCh:       make(chan struct{}),
	}

	go limiter.cleanupLoop()
	return limiter
}

// Allow checks whether a request from the given IP is allowed.
func (l *IPRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	bucket, exists := l.buckets[ip]
	if !exists {
		l.buckets[ip] = &tokenBucket{
			tokens:     l.capacity - 1.0,
			capacity:   l.capacity,
			refillRate: l.refillRate,
			lastUpdate: now,
		}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(bucket.lastUpdate).Seconds()
	bucket.tokens += elapsed * bucket.refillRate
	if bucket.tokens > bucket.capacity {
		bucket.tokens = bucket.capacity
	}
	bucket.lastUpdate = now

	if bucket.tokens >= 1.0 {
		bucket.tokens -= 1.0
		return true
	}

	return false
}

// Stop terminates the background cleanup goroutine.
func (l *IPRateLimiter) Stop() {
	select {
	case <-l.stopCh:
	default:
		close(l.stopCh)
	}
}

func (l *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(l.cleanupEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.mu.Lock()
			now := time.Now()
			for ip, bucket := range l.buckets {
				// Remove if inactive for 10 minutes and full tokens
				if now.Sub(bucket.lastUpdate) > 10*time.Minute {
					delete(l.buckets, ip)
				}
			}
			l.mu.Unlock()
		case <-l.stopCh:
			return
		}
	}
}
