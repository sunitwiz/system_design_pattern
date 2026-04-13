package cache

import "fmt"

type EvictionPolicy int


const (
	LRU EvictionPolicy = iota
	LFU
)

func NewCache(policy EvictionPolicy, capacity int) (Cache, error) {
	if capacity <= 0 {
		return nil, fmt.Errorf("capacity must be positive, got %d", capacity)
	}
	switch policy {
	case LRU:
		return newLRUCache(capacity), nil
	case LFU:
		return newLFUCache(capacity), nil
	default:
		return nil, fmt.Errorf("unknown eviction policy: %d", policy)
	}
}
