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

func (f *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	return true, time.Second
}
