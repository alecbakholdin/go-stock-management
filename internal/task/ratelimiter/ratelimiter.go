package ratelimiter

import (
	"context"
	"time"
)

type RateLimiter struct {
	reqGap   time.Duration
	rateChan chan int
}

func New(reqGap time.Duration) *RateLimiter {
	r := &RateLimiter{
		reqGap:   reqGap,
		rateChan: make(chan int),
	}
	go func() {
		r.rateChan <- 1
	}()
	return r
}

// blocking request called on this rate limiter
func (r *RateLimiter) Acquire(ctx context.Context) {
	select {
	case <-r.rateChan:
		go func() {
			time.Sleep(r.reqGap)
			r.rateChan <- 1
		}()
		return
	case <-ctx.Done():
		return
	}
}
