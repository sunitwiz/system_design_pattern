package algorithm

import (
	"rate_limiter/config"
	"sync"
	"time"
)

type slidingWindow struct {
	mu          sync.Mutex
	maxRequests int
	windowSize  time.Duration
	requestLog  map[string][]time.Time
}

func newSlidingWindow(cfg config.RateLimiterConfig) *slidingWindow {
	return &slidingWindow{
		maxRequests: cfg.MaxRequests,
		windowSize:  cfg.WindowSize,
		requestLog:  make(map[string][]time.Time),
	}
}

func (sw *slidingWindow) Allow(key string) bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-sw.windowSize)

	var valid []time.Time
	for _, t := range sw.requestLog[key] {
		if t.After(windowStart) {
			valid = append(valid, t)
		}
	}
	sw.requestLog[key] = valid

	if len(valid) >= sw.maxRequests {
		return false
	}

	sw.requestLog[key] = append(sw.requestLog[key], now)
	return true
}

func (sw *slidingWindow) GetLimit() int                    { return sw.maxRequests }
func (sw *slidingWindow) GetWindowSize() time.Duration     { return sw.windowSize }
func (sw *slidingWindow) GetAlgorithmType() AlgorithmType  { return SlidingWindowType }
