package main

import (
	"fmt"
	"rate_limiter/algorithm"
	"rate_limiter/config"
	"rate_limiter/ratelimiter"
	"time"
)

func main() {
	fmt.Println("=== Rate Limiter System Demo ===")
	fmt.Println()

	service := ratelimiter.GetInstance()

	fmt.Println("--- Registering Routes ---")

	searchLimiter, _ := algorithm.NewRateLimiter(algorithm.TokenBucketType, config.RateLimiterConfig{
		BucketCapacity: 10,
		RefillRate:     2.0,
	})
	service.RegisterRoute("/api/search", searchLimiter)

	loginLimiter, _ := algorithm.NewRateLimiter(algorithm.FixedWindowType, config.RateLimiterConfig{
		MaxRequests: 5,
		WindowSize:  60 * time.Second,
	})
	service.RegisterRoute("/api/login", loginLimiter)

	dataLimiter, _ := algorithm.NewRateLimiter(algorithm.SlidingWindowType, config.RateLimiterConfig{
		MaxRequests: 10,
		WindowSize:  60 * time.Second,
	})
	service.RegisterRoute("/api/data", dataLimiter)

	fmt.Println()
	service.ViewStatus()

	fmt.Println("\n--- Token Bucket: /api/search (capacity=10, refill=2/sec) ---")
	simulateRequests(service, "/api/search", "client-A", 12)

	fmt.Println("\n--- Fixed Window: /api/login (5 requests per 60s) ---")
	simulateRequests(service, "/api/login", "client-B", 8)

	fmt.Println("\n--- Sliding Window: /api/data (10 requests per 60s) ---")
	simulateRequests(service, "/api/data", "client-C", 12)

	fmt.Println("\n--- Multiple Clients on Same Route ---")
	fmt.Println("  Client-X on /api/login:")
	simulateRequests(service, "/api/login", "client-X", 6)
	fmt.Println("  Client-Y on /api/login:")
	simulateRequests(service, "/api/login", "client-Y", 6)

	fmt.Println("\n--- Rate Limit Recovery (Token Bucket) ---")
	fmt.Println("  Sending 10 requests to exhaust bucket...")
	for i := 0; i < 10; i++ {
		service.AllowRequest("/api/search", "client-recovery")
	}
	fmt.Printf("  Request 11 (exhausted): allowed=%v\n", service.AllowRequest("/api/search", "client-recovery"))

	fmt.Println("  Waiting 3 seconds for token refill (rate=2/sec → ~6 new tokens)...")
	time.Sleep(3 * time.Second)

	fmt.Println("  After waiting:")
	for i := 1; i <= 7; i++ {
		allowed := service.AllowRequest("/api/search", "client-recovery")
		fmt.Printf("    Request %d: allowed=%v\n", i, allowed)
	}

	fmt.Println("\n--- Route Stats ---")
	for _, route := range []string{"/api/search", "/api/login", "/api/data"} {
		stats := service.GetStats(route)
		fmt.Printf("  %-15s algorithm=%-22s limit=%d  window=%s\n",
			stats["route"], stats["algorithm"], stats["limit"], stats["windowSize"])
	}

	fmt.Println("\n--- Edge Case: Unregistered Route ---")
	allowed := service.AllowRequest("/api/unknown", "client-Z")
	fmt.Printf("  /api/unknown → allowed=%v (no limiter = always allowed)\n", allowed)

	fmt.Println("\n--- Removing a Route ---")
	if err := service.RemoveRoute("/api/data"); err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
	fmt.Printf("  /api/data after removal → allowed=%v\n",
		service.AllowRequest("/api/data", "client-C"))

	fmt.Println("\n--- Final Status ---")
	service.ViewStatus()
}

func simulateRequests(service *ratelimiter.RateLimiterService, route string, clientID string, count int) {
	allowed, blocked := 0, 0
	for i := 0; i < count; i++ {
		if service.AllowRequest(route, clientID) {
			allowed++
		} else {
			blocked++
		}
	}
	fmt.Printf("  %s → %d requests: %d allowed, %d blocked\n", clientID, count, allowed, blocked)
}
