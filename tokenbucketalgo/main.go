package main

import (
	"fmt"
	"sync"
	"time"
)

type TokenBucket struct {
	fillrate   int
	capacity   int
	lastrefill time.Time
	token      int
	mu         sync.Mutex
}

func NewTokenBucket(fillrate int, capacity int) *TokenBucket {
	return &TokenBucket{
		fillrate:   fillrate,
		capacity:   capacity,
		token:      capacity,
		lastrefill: time.Now(),
	}

}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	secondsElapseds := time.Since(tb.lastrefill).Seconds()
	tokensToAdd := tb.fillrate * int(secondsElapseds)
	if tokensToAdd > 0 {
		tb.token = min(tb.capacity, tb.token+tokensToAdd)
		tb.lastrefill = time.Now()
	}
	if tb.token > 0 {
		tb.token -= 1
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
