package ratelimit

import (
	"context"
	"errors"
	"sync"
)

var ErrConcurrencyLimitReached = errors.New("concurrency limit reached")

// Semaphore implements a weighted counting semaphore with try-acquire capabilities.
type Semaphore struct {
	mu      sync.Mutex
	current int
	max     int
	ch      chan struct{}
}

// NewSemaphore creates a semaphore with a maximum concurrency limit.
func NewSemaphore(max int) *Semaphore {
	if max <= 0 {
		max = 10
	}
	return &Semaphore{
		max: max,
		ch:  make(chan struct{}, max),
	}
}

// TryAcquire attempts to acquire a slot immediately without blocking.
// Returns true if acquired, false otherwise.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		s.mu.Lock()
		s.current++
		s.mu.Unlock()
		return true
	default:
		return false
	}
}

// Acquire blocks until a slot is available or context is cancelled.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.ch <- struct{}{}:
		s.mu.Lock()
		s.current++
		s.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release releases an acquired slot.
func (s *Semaphore) Release() {
	select {
	case <-s.ch:
		s.mu.Lock()
		if s.current > 0 {
			s.current--
		}
		s.mu.Unlock()
	default:
	}
}

// Current returns the number of currently active slots.
func (s *Semaphore) Current() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}
