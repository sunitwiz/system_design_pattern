package algorithm

import "time"

type RateLimiter interface {
	Allow(key string) bool
	GetLimit() int
	GetWindowSize() time.Duration
	GetAlgorithmType() AlgorithmType
}
