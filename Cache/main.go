package main

import (
	"cache/cache"
	"cache/cacheservice"
	"fmt"
)

func main() {
	fmt.Println("=== Cache System Demo ===\n")

	svc := cacheservice.GetInstance()

	fmt.Println("--- LRU Cache (capacity=3) ---")
	svc.CreateCache("session-cache", cache.LRU, 3)
	lru, _ := svc.GetCache("session-cache")

	fmt.Println("\nPutting 3 items: A=1, B=2, C=3")
	lru.Put("A", 1)
	lru.Put("B", 2)
	lru.Put("C", 3)
	fmt.Printf("  State: %s\n", lru)

	fmt.Println("\nGet B (moves B to front):")
	val, ok := lru.Get("B")
	fmt.Printf("  Get(B) = %v, found=%v\n", val, ok)
	fmt.Printf("  State: %s\n", lru)

	fmt.Println("\nPut D=4 (evicts least recently used):")
	lru.Put("D", 4)
	fmt.Printf("  State: %s\n", lru)

	fmt.Println("\nGet A (should miss — was evicted):")
	val, ok = lru.Get("A")
	fmt.Printf("  Get(A) = %v, found=%v\n", val, ok)

	fmt.Println("\nUpdate existing: Put B=20")
	lru.Put("B", 20)
	fmt.Printf("  State: %s\n", lru)

	fmt.Println("\nDelete C:")
	lru.Delete("C")
	fmt.Printf("  State: %s\n", lru)

	fmt.Println("\n--- LFU Cache (capacity=3) ---")
	svc.CreateCache("product-cache", cache.LFU, 3)
	lfu, _ := svc.GetCache("product-cache")

	fmt.Println("\nPutting 3 items: X=10, Y=20, Z=30")
	lfu.Put("X", 10)
	lfu.Put("Y", 20)
	lfu.Put("Z", 30)
	fmt.Printf("  State: %s\n", lfu)

	fmt.Println("\nAccess X twice, Y once (X freq=3, Y freq=2, Z freq=1):")
	lfu.Get("X")
	lfu.Get("X")
	lfu.Get("Y")
	fmt.Printf("  State: %s\n", lfu)

	fmt.Println("\nPut W=40 (evicts Z — lowest frequency):")
	lfu.Put("W", 40)
	fmt.Printf("  State: %s\n", lfu)

	fmt.Println("\nGet Z (should miss — was evicted):")
	val, ok = lfu.Get("Z")
	fmt.Printf("  Get(Z) = %v, found=%v\n", val, ok)

	fmt.Println("\n--- Eviction Comparison (LRU vs LFU) ---")
	fmt.Println("Same operations on both: Put A,B,C → Get A, Get A, Get B, Get C → Put D\n")

	svc.CreateCache("lru-compare", cache.LRU, 3)
	svc.CreateCache("lfu-compare", cache.LFU, 3)
	lruCmp, _ := svc.GetCache("lru-compare")
	lfuCmp, _ := svc.GetCache("lfu-compare")

	for _, c := range []cache.Cache{lruCmp, lfuCmp} {
		c.Put("A", 1)
		c.Put("B", 2)
		c.Put("C", 3)
		c.Get("A")
		c.Get("A")
		c.Get("B")
		c.Get("C")
	}

	fmt.Println("  Before eviction:")
	fmt.Printf("    %s\n", lruCmp)
	fmt.Printf("    %s\n", lfuCmp)

	lruCmp.Put("D", 4)
	lfuCmp.Put("D", 4)

	fmt.Println("\n  After Put(D):")
	fmt.Printf("    LRU evicts least recently used  → %s\n", lruCmp)
	fmt.Printf("    LFU evicts least frequently used → %s\n", lfuCmp)
	fmt.Println("\n  LRU evicted A (least recently used, despite being most accessed)")
	fmt.Println("  LFU evicted B (least frequently used, LRU tiebreaker among freq=2)")

	fmt.Println()
	svc.ViewStatus()
}
