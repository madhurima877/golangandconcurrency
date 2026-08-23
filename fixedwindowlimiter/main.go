package main

import (
	"log"
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	limit       int
	counter     int
	window      time.Duration
	windowStart time.Time
	mu          sync.Mutex
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:       limit,
		window:      window,
		windowStart: time.Now(),
		counter:     0,
	}

}
func (f *FixedWindowLimiter) Allow() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if time.Since(f.windowStart) >= f.window {
		f.windowStart = time.Now()
		f.counter = 0
	}
	if f.counter >= f.limit {
		return false
	}
	f.counter++
	return true
}
func main() {
	limiter := NewFixedWindowLimiter(5, 10*time.Second)
	for i := 1; i <= 100; i++ {
		if limiter.Allow() {
			log.Println("Allowed")
		} else {
			log.Println("Rejected")
		}
		time.Sleep(time.Second)
	}
}
