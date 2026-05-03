package ratelimiter

import "time"

type Limiter interface {
	Allow(ip string) (bool, int, time.Duration)
}

type Config struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}
