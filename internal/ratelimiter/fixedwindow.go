package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	mu         sync.Mutex
	counters   map[string]*CurrentWindow
	limit      int
	windowSize time.Duration
}
type CurrentWindow struct {
	count     int
	expiresAt time.Time
}

func NewFixedWindowRateLimiter(limit int, windowSize time.Duration) *FixedWindowRateLimiter {

	return &FixedWindowRateLimiter{
		counters:   make(map[string]*CurrentWindow),
		limit:      limit,
		windowSize: windowSize,
	}

}

func (f *FixedWindowRateLimiter) Allow(ip string) (bool, int, time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	val, exists := f.counters[ip]

	if !exists || now.After(val.expiresAt) {
		expiresAt := now.Add(f.windowSize)
		f.counters[ip] = &CurrentWindow{
			count:     1,
			expiresAt: expiresAt,
		}
		return true, f.limit - 1, time.Until(expiresAt)
	}

	if val.count < f.limit {
		val.count++
		return true, f.limit - val.count, time.Until(val.expiresAt)
	}

	return false, 0, time.Until(val.expiresAt)
}
