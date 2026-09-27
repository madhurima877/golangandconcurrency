package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	limiters map[string]*FixedWindow
	mu       sync.Mutex
	limit    int
	window   time.Duration
}
type FixedWindow struct {
	window      time.Duration
	limit       int
	counter     int
	windowStart time.Time
	mu          sync.Mutex
}

func NewFixedWindow(limit int, window time.Duration) *FixedWindow {
	return &FixedWindow{limit: limit, window: window, windowStart: time.Now(), counter: 0}

}

func (f *FixedWindow) Allow() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if time.Since(f.windowStart) >= f.window {
		f.windowStart = time.Now()
		f.counter = 0
	}
	if f.limit <= f.counter {
		return false
	}
	f.counter++
	return true
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*FixedWindow),
		limit:    limit,
		window:   window,
	}
}

func (r *RateLimiter) Allow(userID string) bool {
	r.mu.Lock()

	limiter, exists := r.limiters[userID]
	if !exists {
		limiter = NewFixedWindow(r.limit, r.window)
		r.limiters[userID] = limiter
	}
	r.mu.Unlock()
	return limiter.Allow()
}
func main() {
	limiter := NewRateLimiter(5, time.Second)
	http.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")

		if !limiter.Allow(userID) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("request allowed"))
	})

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
