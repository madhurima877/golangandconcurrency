package main

import (
	"sync"
	"time"
)

type RateLimiter struct {
	capacity   int
	tokens     float64
	refillRate float64
	mu         sync.Mutex
	lastRefill time.Time
}

func NewRateLimiter(capacity int, refillRate float64) *RateLimiter {
	return &RateLimiter{
		capacity:   capacity,
		refillRate: refillRate,
		tokens:     0,
		lastRefill: time.Now(),
	}
}
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	seconds := time.Since(rl.lastRefill).Seconds()
	newtokens := seconds * rl.refillRate
	if newtokens > 0 {
		rl.lastRefill = time.Now()
		rl.tokens = min(float64(rl.capacity), rl.tokens+newtokens)
	}
	if rl.tokens > 0 {
		rl.tokens -= 1
		return true
	}
	return false
}
