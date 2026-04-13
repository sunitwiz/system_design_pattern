package cacheservice

import (
	"cache/cache"
	"fmt"
	"sort"
	"sync"
)

type CacheService struct {
	mu     sync.Mutex
	caches map[string]cache.Cache
}

var (
	instance *CacheService
	once     sync.Once
)

func GetInstance() *CacheService {
	once.Do(func() {
		instance = &CacheService{
			caches: make(map[string]cache.Cache),
		}
	})
	return instance
}

func ResetInstance() {
	once = sync.Once{}
	instance = nil
}

func (cs *CacheService) CreateCache(name string, policy cache.EvictionPolicy, capacity int) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, exists := cs.caches[name]; exists {
		return fmt.Errorf("cache %q already exists", name)
	}

	c, err := cache.NewCache(policy, capacity)
	if err != nil {
		return err
	}
	cs.caches[name] = c
	return nil
}

func (cs *CacheService) GetCache(name string) (cache.Cache, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	c, exists := cs.caches[name]
	if !exists {
		return nil, fmt.Errorf("cache %q not found", name)
	}
	return c, nil
}

func (cs *CacheService) DeleteCache(name string) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, exists := cs.caches[name]; !exists {
		return fmt.Errorf("cache %q not found", name)
	}
	delete(cs.caches, name)
	return nil
}

func (cs *CacheService) ListCaches() []string {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	names := make([]string, 0, len(cs.caches))
	for name := range cs.caches {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (cs *CacheService) ViewStatus() {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	fmt.Println("========================================")
	fmt.Println("        CACHE SERVICE STATUS")
	fmt.Println("========================================")

	if len(cs.caches) == 0 {
		fmt.Println("  No caches configured.")
		return
	}

	names := make([]string, 0, len(cs.caches))
	for name := range cs.caches {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		c := cs.caches[name]
		fmt.Printf("\n  %-18s %s\n", name+":", c)
	}

	fmt.Printf("\n  Total Caches: %d\n", len(cs.caches))
	fmt.Println("========================================")
}
