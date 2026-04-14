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

type FixedWindow struct {
	mu          sync.Mutex
	maxRequests int
	windowSize  time.Duration
	counters    map[string]*windowEntry
}

func newFixedWindow(cfg config.RateLimiterConfig) *FixedWindow {
	return &FixedWindow{
		maxRequests: cfg.MaxRequests,
		windowSize:  cfg.WindowSize,
		counters:    make(map[string]*windowEntry),
	}
}

func (fw *FixedWindow) Allow(key string) bool {
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

func (fw *FixedWindow) GetLimit() int                    { return fw.maxRequests }
func (fw *FixedWindow) GetWindowSize() time.Duration     { return fw.windowSize }
func (fw *FixedWindow) GetAlgorithmType() AlgorithmType  { return FixedWindowType }
