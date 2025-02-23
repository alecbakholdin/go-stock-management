package ratelimiter_test

import (
	"context"
	"stock-management/internal/task/ratelimiter"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOnePerSecond(t *testing.T) {
	t.Parallel()
	start := time.Now()
	gap := time.Millisecond * 50
	r := ratelimiter.New(gap)

	for i := range 5 {
		r.Acquire(t.Context())
		assert.Greater(t, time.Since(start), gap*time.Duration(i))
		assert.Less(t, time.Since(start), gap*time.Duration(i+1))
	}
}

func TestMultipleCallers(t *testing.T) {
	t.Parallel()
	start := time.Now()
	gap := time.Millisecond * 50
	r := ratelimiter.New(gap)
	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r.Acquire(t.Context())
		}()
	}
	wg.Wait()
	assert.Greater(t, time.Since(start), gap*4)
	assert.Less(t, time.Since(start), gap*5)
}

func TestCancelledContextBeforeAcquire(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	start := time.Now()
	gap := time.Second
	r := ratelimiter.New(gap)
	cancel()
	r.Acquire(ctx)
	r.Acquire(ctx)
	assert.Less(t, time.Since(start), gap)
}

func TestCancelledContextDuringAcquire(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(t.Context())
	start := time.Now()
	gap := time.Second
	r := ratelimiter.New(gap)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		r.Acquire(ctx)
		wg.Done()
		r.Acquire(ctx)
	}()
	wg.Wait()
	cancel()
	assert.Less(t, time.Since(start), gap)
}
