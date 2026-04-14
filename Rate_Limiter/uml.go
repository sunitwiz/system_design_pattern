package main

import "fmt"

func printClassDiagram() {
	fmt.Println(`classDiagram

    class AlgorithmType {
        <<enumeration>>
        TokenBucketType
        SlidingWindowType
        FixedWindowType
        func (a AlgorithmType) String() string
    }

    class RateLimiterConfig {
        MaxRequests    int
        WindowSize     time.Duration
        BucketCapacity int
        RefillRate     float64
    }

    class RateLimiter {
        <<interface>>
        Allow(key string) bool
        GetLimit() int
        GetWindowSize() time.Duration
        GetAlgorithmType() AlgorithmType
    }

    class windowEntry {
        windowStart time.Time
        count       int
    }

    class FixedWindow {
        mu          sync.Mutex
        maxRequests int
        windowSize  time.Duration
        counters    map[string]*windowEntry
        func newFixedWindow(cfg config.RateLimiterConfig) *FixedWindow
        func (fw *FixedWindow) Allow(key string) bool
        func (fw *FixedWindow) GetLimit() int
        func (fw *FixedWindow) GetWindowSize() time.Duration
        func (fw *FixedWindow) GetAlgorithmType() AlgorithmType
    }

    class SlidingWindow {
        mu          sync.Mutex
        maxRequests int
        windowSize  time.Duration
        requestLog  map[string][]time.Time
        func newSlidingWindow(cfg config.RateLimiterConfig) *SlidingWindow
        func (sw *SlidingWindow) Allow(key string) bool
        func (sw *SlidingWindow) GetLimit() int
        func (sw *SlidingWindow) GetWindowSize() time.Duration
        func (sw *SlidingWindow) GetAlgorithmType() AlgorithmType
    }

    class TokenBucket {
        mu         sync.Mutex
        capacity   int
        refillRate float64
        tokens     map[string]float64
        lastRefill map[string]time.Time
        func newTokenBucket(cfg config.RateLimiterConfig) *TokenBucket
        func (tb *TokenBucket) Allow(key string) bool
        func (tb *TokenBucket) GetLimit() int
        func (tb *TokenBucket) GetWindowSize() time.Duration
        func (tb *TokenBucket) GetAlgorithmType() AlgorithmType
    }

    class RateLimiterService {
        mu       sync.Mutex
        limiters map[string]algorithm.RateLimiter
        func GetInstance() *RateLimiterService
        func (s *RateLimiterService) RegisterRoute(route string, limiter algorithm.RateLimiter)
        func (s *RateLimiterService) AllowRequest(route string, clientID string) bool
        func (s *RateLimiterService) GetStats(route string) map[string]interface{}
        func (s *RateLimiterService) RemoveRoute(route string) error
        func (s *RateLimiterService) ViewStatus()
    }

    FixedWindow ..|> RateLimiter : implements
    SlidingWindow ..|> RateLimiter : implements
    TokenBucket ..|> RateLimiter : implements
    FixedWindow *-- windowEntry : tracks per-key windows
    FixedWindow --> AlgorithmType
    SlidingWindow --> AlgorithmType
    TokenBucket --> AlgorithmType
    RateLimiterService o-- RateLimiter : manages per route`)
}
