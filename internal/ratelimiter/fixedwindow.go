package ratelimiter

import (
	"fmt"
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

func (f *FixedWindowRateLimiter) Allow(ip string) bool {

	f.mu.Lock()
	defer f.mu.Unlock()
	now := time.Now()
	val, exists := f.counters[ip]
	// fmt.Println(exists)
	fmt.Println(len(f.counters))
	fmt.Println(" wb")
	if !exists || now.After(val.expiresAt) {
		// fmt.Println(ip, "   ip")
		f.counters[ip] = &CurrentWindow{
			count:     1,
			expiresAt: now.Add(f.windowSize),
		}
		return true

	}
	if val.count < f.limit {
		val.count++
		return true
	}

	return false

}
