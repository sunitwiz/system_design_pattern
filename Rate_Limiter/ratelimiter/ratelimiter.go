package ratelimiter

import (
	"fmt"
	"rate_limiter/algorithm"
	"sync"
)

type RateLimiterService struct {
	mu       sync.Mutex
	limiters map[string]algorithm.RateLimiter
}

var (
	instance *RateLimiterService
	once     sync.Once
)

func GetInstance() *RateLimiterService {
	once.Do(func() {
		instance = &RateLimiterService{
			limiters: make(map[string]algorithm.RateLimiter),
		}
	})
	return instance
}

func ResetInstance() {
	once = sync.Once{}
	instance = nil
}

func (s *RateLimiterService) RegisterRoute(route string, limiter algorithm.RateLimiter) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.limiters[route] = limiter
	fmt.Printf("  Registered route %-15s → %s (limit: %d)\n",
		route, limiter.GetAlgorithmType(), limiter.GetLimit())
}

func (s *RateLimiterService) AllowRequest(route string, clientID string) bool {
	s.mu.Lock()
	limiter, exists := s.limiters[route]
	s.mu.Unlock()

	if !exists {
		return true
	}

	return limiter.Allow(clientID)
}

func (s *RateLimiterService) GetStats(route string) map[string]interface{} {
	s.mu.Lock()
	limiter, exists := s.limiters[route]
	s.mu.Unlock()

	if !exists {
		return nil
	}

	return map[string]interface{}{
		"route":      route,
		"algorithm":  limiter.GetAlgorithmType().String(),
		"limit":      limiter.GetLimit(),
		"windowSize": limiter.GetWindowSize().String(),
	}
}

func (s *RateLimiterService) RemoveRoute(route string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.limiters[route]; !exists {
		return fmt.Errorf("route %s not found", route)
	}

	delete(s.limiters, route)
	fmt.Printf("  Removed route %s\n", route)
	return nil
}

func (s *RateLimiterService) ViewStatus() {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("========================================")
	fmt.Println("      RATE LIMITER STATUS")
	fmt.Println("========================================")

	if len(s.limiters) == 0 {
		fmt.Println("  No routes configured.")
		return
	}

	fmt.Printf("  %-18s %-22s %s\n", "Route", "Algorithm", "Limit")
	fmt.Printf("  %-18s %-22s %s\n", "------------------", "----------------------", "-----")

	for route, limiter := range s.limiters {
		windowInfo := ""
		if limiter.GetWindowSize() > 0 {
			windowInfo = fmt.Sprintf("%d / %s", limiter.GetLimit(), limiter.GetWindowSize())
		} else {
			windowInfo = fmt.Sprintf("%d capacity", limiter.GetLimit())
		}
		fmt.Printf("  %-18s %-22s %s\n", route, limiter.GetAlgorithmType(), windowInfo)
	}

	fmt.Println("========================================")
}
