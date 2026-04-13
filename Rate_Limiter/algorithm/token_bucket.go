package algorithm

import (
	"rate_limiter/config"
	"sync"
	"time"
)

type tokenBucket struct {
	mu         sync.Mutex
	capacity   int
	refillRate float64
	tokens     map[string]float64
	lastRefill map[string]time.Time
}

func newTokenBucket(cfg config.RateLimiterConfig) *tokenBucket {
	return &tokenBucket{
		capacity:   cfg.BucketCapacity,
		refillRate: cfg.RefillRate,
		tokens:     make(map[string]float64),
		lastRefill: make(map[string]time.Time),
	}
}

func (tb *tokenBucket) Allow(key string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()

	if _, exists := tb.tokens[key]; !exists {
		tb.tokens[key] = float64(tb.capacity)
		tb.lastRefill[key] = now
	}

	elapsed := now.Sub(tb.lastRefill[key]).Seconds()
	tb.tokens[key] += elapsed * tb.refillRate
	if tb.tokens[key] > float64(tb.capacity) {
		tb.tokens[key] = float64(tb.capacity)
	}
	tb.lastRefill[key] = now

	if tb.tokens[key] >= 1.0 {
		tb.tokens[key] -= 1.0
		return true
	}

	return false
}

func (tb *tokenBucket) GetLimit() int              { return tb.capacity }
func (tb *tokenBucket) GetWindowSize() time.Duration { return 0 }
func (tb *tokenBucket) GetAlgorithmType() AlgorithmType { return TokenBucketType }
