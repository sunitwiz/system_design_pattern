package algorithm

import (
	"rate_limiter/config"
	"sync"
	"time"
)

type windowEntry struct {
	windowStart time.Time
	count       int
}

type fixedWindow struct {
	mu          sync.Mutex
	maxRequests int
	windowSize  time.Duration
	counters    map[string]*windowEntry
}

func newFixedWindow(cfg config.RateLimiterConfig) *fixedWindow {
	return &fixedWindow{
		maxRequests: cfg.MaxRequests,
		windowSize:  cfg.WindowSize,
		counters:    make(map[string]*windowEntry),
	}
}

func (fw *fixedWindow) Allow(key string) bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	now := time.Now()

	entry, exists := fw.counters[key]
	if !exists || now.Sub(entry.windowStart) >= fw.windowSize {
		fw.counters[key] = &windowEntry{windowStart: now, count: 1}
		return true
	}

	if entry.count >= fw.maxRequests {
		return false
	}

	entry.count++
	return true
}

func (fw *fixedWindow) GetLimit() int                    { return fw.maxRequests }
func (fw *fixedWindow) GetWindowSize() time.Duration     { return fw.windowSize }
func (fw *fixedWindow) GetAlgorithmType() AlgorithmType  { return FixedWindowType }
