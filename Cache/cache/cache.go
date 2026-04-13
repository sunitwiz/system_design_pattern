package cache

import "fmt"

type EvictionPolicy int

const (
	LRU EvictionPolicy = iota
	LFU
)

func (e EvictionPolicy) String() string {
	switch e {
	case LRU:
		return "LRU"
	case LFU:
		return "LFU"
	default:
		return "Unknown"
	}
}

type Cache interface {
	Get(key string) (any, bool)
	Put(key string, value any)
	Delete(key string) bool
	Size() int
	Capacity() int
	Clear()
	String() string
}

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
