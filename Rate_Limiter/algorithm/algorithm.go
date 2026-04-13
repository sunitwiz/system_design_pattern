package algorithm

import (
	"fmt"
	"rate_limiter/config"
	"time"
)

type AlgorithmType int

const (
	TokenBucketType AlgorithmType = iota
	SlidingWindowType
	FixedWindowType
)

func (a AlgorithmType) String() string {
	switch a {
	case TokenBucketType:
		return "Token Bucket"
	case SlidingWindowType:
		return "Sliding Window Log"
	case FixedWindowType:
		return "Fixed Window Counter"
	default:
		return "Unknown"
	}
}

type RateLimiter interface {
	Allow(key string) bool
	GetLimit() int
	GetWindowSize() time.Duration
	GetAlgorithmType() AlgorithmType
}

func NewRateLimiter(algorithmType AlgorithmType, cfg config.RateLimiterConfig) (RateLimiter, error) {
	switch algorithmType {
	case TokenBucketType:
		return newTokenBucket(cfg), nil
	case SlidingWindowType:
		return newSlidingWindow(cfg), nil
	case FixedWindowType:
		return newFixedWindow(cfg), nil
	default:
		return nil, fmt.Errorf("unknown algorithm type: %d", algorithmType)
	}
}
