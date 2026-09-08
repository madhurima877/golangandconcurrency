package main

import (
	"log"
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	window      time.Duration
	windowStart time.Time
	limit       int
	counter     int
	mu          sync.Mutex
}

func NewFixedWindowLimiter(window time.Duration, limit int) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		window:      window,
		limit:       limit,
		counter:     0,
		windowStart: time.Now(),
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
	limiter := NewFixedWindowLimiter(10*time.Second, 5)
	for i := 1; i <= 100; i++ {
		if limiter.Allow() {
			log.Println("Allowed")
		} else {
			log.Println("Rejected")
		}
		time.Sleep(time.Second)
	}

}
