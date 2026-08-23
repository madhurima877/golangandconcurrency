package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	tokens     int
	capacity   int
	refillrate int
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(capacity, refillrate int) *TokenBucket {
	return &TokenBucket{
		tokens:     capacity,
		capacity:   capacity,
		refillrate: refillrate,
		lastRefill: time.Now(),
	}
}
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	seconds := time.Since(tb.lastRefill).Seconds()
	tokensToAdd := int(seconds * float64(tb.refillrate))
	if tokensToAdd > 0 {
		tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)
		tb.lastRefill = time.Now()
	}
	if tb.tokens > 0 {
		tb.tokens -= 1
		return true
	}
	return false
}

func main() {
	limiter := NewTokenBucket(5, 1)
	for i := 1; i <= 20; i++ {
		if limiter.Allow() {
			fmt.Println("Request", i, "allowed")
		} else {
			fmt.Println("Request", i, "rejected")
		}
	}
	time.Sleep(3 * time.Second)

	fmt.Println("After waiting...")

	for i := 1; i <= 5; i++ {
		if limiter.Allow() {
			fmt.Println("Request", i, "allowed")
		} else {
			fmt.Println("Request", i, "rejected")
		}
	}

}
