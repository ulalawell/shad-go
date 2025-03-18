package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	maxCount   int
	interval   time.Duration
	requests   chan request
	stop       chan struct{}
	timestamps []time.Time
	pending    []request
}

type request struct {
	ctx  context.Context
	done chan error
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	l := &Limiter{
		maxCount:   maxCount,
		interval:   interval,
		requests:   make(chan request),
		stop:       make(chan struct{}),
		timestamps: make([]time.Time, 0, maxCount),
		pending:    make([]request, 0),
	}

	go l.run()
	return l
}

func (l *Limiter) Acquire(ctx context.Context) error {
	req := request{
		ctx:  ctx,
		done: make(chan error, 1),
	}

	select {
	case l.requests <- req:
	case <-l.stop:
		return ErrStopped
	case <-ctx.Done():
		return ctx.Err()
	}

	select {
	case err := <-req.done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *Limiter) Stop() {
	close(l.stop)
}

func (l *Limiter) run() {
	for {
		select {
		case req := <-l.requests:
			now := time.Now()

			for len(l.timestamps) > 0 && now.Sub(l.timestamps[0]) > l.interval {
				l.timestamps = l.timestamps[1:]
			}

			if len(l.timestamps) < l.maxCount {

				l.timestamps = append(l.timestamps, now)
				req.done <- nil
			} else {
				l.pending = append(l.pending, req)
			}

		case <-time.After(l.interval):
			now := time.Now()
			for len(l.timestamps) > 0 && now.Sub(l.timestamps[0]) > l.interval {
				l.timestamps = l.timestamps[1:]
			}

			for len(l.pending) > 0 && len(l.timestamps) < l.maxCount {
				req := l.pending[0]
				l.pending = l.pending[1:]
				l.timestamps = append(l.timestamps, now)
				req.done <- nil
			}

		case <-l.stop:
			return
		}
	}
}
