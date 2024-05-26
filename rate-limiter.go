package main

import (
	"fmt"
	"sync"
	"time"
)

// SlidingWindowRateLimiter defines the structure of the rate limiter
type SlidingWindowRateLimiter struct {
	mu          sync.Mutex
	requests    []int64
	windowSize  int64 // in nanoseconds
	maxRequests int
}

// NewSlidingWindowRateLimiter initializes a new sliding window rate limiter
func NewSlidingWindowRateLimiter(windowSize time.Duration, maxRequests int) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		windowSize:  int64(windowSize),
		maxRequests: maxRequests,
	}
}

// AllowRequest checks if a new request can be allowed
func (rl *SlidingWindowRateLimiter) AllowRequest() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now().UnixNano()
	windowStart := now - rl.windowSize

	// Remove old requests that are outside the current window
	rl.removeOldRequests(windowStart)

	if len(rl.requests) < rl.maxRequests {
		rl.requests = append(rl.requests, now)
		return true
	}
	return false
}

// removeOldRequests removes requests that are outside the current sliding window
func (rl *SlidingWindowRateLimiter) removeOldRequests(windowStart int64) {
	newRequests := rl.requests[:0]
	for _, timestamp := range rl.requests {
		if timestamp >= windowStart {
			newRequests = append(newRequests, timestamp)
		}
	}
	rl.requests = newRequests
}

func main() {
	windowSize := time.Second * 10 // 10-second window
	maxRequests := 5               // maximum of 5 requests per window
	limiter := NewSlidingWindowRateLimiter(windowSize, maxRequests)

	for i := 0; i < 10; i++ {
		if limiter.AllowRequest() {
			fmt.Printf("Request %d allowed\n", i+1)
		} else {
			fmt.Printf("Request %d denied\n", i+1)
		}
		time.Sleep(time.Second)
	}
}
