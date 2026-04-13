package config

import "time"

type RateLimiterConfig struct {
	MaxRequests    int
	WindowSize     time.Duration
	BucketCapacity int
	RefillRate     float64
}
