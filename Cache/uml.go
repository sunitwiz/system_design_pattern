package main

import "fmt"

func printClassDiagram() {
	fmt.Println(`classDiagram

    class EvictionPolicy {
        <<enumeration>>
        LRU
        LFU
    }

    class Cache {
        <<interface>>
        Get(key string) (any, bool)
        Put(key string, value any)
        Delete(key string) bool
        Size() int
        Capacity() int
        Clear()
        String() string
    }

    class lruNode {
        key   string
        value any
        prev  *lruNode
        next  *lruNode
    }

    class lruCache {
        mu       sync.RWMutex
        capacity int
        nodeMap  map[string]*lruNode
        head     *lruNode
        tail     *lruNode
        func newLRUCache(capacity int) *lruCache
        func (c *lruCache) Get(key string) (any, bool)
        func (c *lruCache) Put(key string, value any)
        func (c *lruCache) Delete(key string) bool
        func (c *lruCache) Size() int
        func (c *lruCache) Capacity() int
        func (c *lruCache) Clear()
        func (c *lruCache) String() string
        func (c *lruCache) addToFront(node *lruNode)
        func (c *lruCache) removeNode(node *lruNode)
        func (c *lruCache) moveToFront(node *lruNode)
        func (c *lruCache) evict()
    }

    class lfuNode {
        key       string
        value     any
        frequency int
        prev      *lfuNode
        next      *lfuNode
    }

    class freqList {
        head *lfuNode
        tail *lfuNode
        size int
        func newFreqList() *freqList
        func (fl *freqList) addToFront(node *lfuNode)
        func (fl *freqList) removeNode(node *lfuNode)
        func (fl *freqList) removeTail() *lfuNode
        func (fl *freqList) isEmpty() bool
    }

    class lfuCache {
        mu       sync.RWMutex
        capacity int
        minFreq  int
        nodeMap  map[string]*lfuNode
        freqMap  map[int]*freqList
        func newLFUCache(capacity int) *lfuCache
        func (c *lfuCache) Get(key string) (any, bool)
        func (c *lfuCache) Put(key string, value any)
        func (c *lfuCache) Delete(key string) bool
        func (c *lfuCache) Size() int
        func (c *lfuCache) Capacity() int
        func (c *lfuCache) Clear()
        func (c *lfuCache) String() string
        func (c *lfuCache) incrementFrequency(node *lfuNode)
        func (c *lfuCache) evict()
    }

    class CacheService {
        mu     sync.Mutex
        caches map[string]cache.Cache
        func GetInstance() *CacheService
        func (cs *CacheService) CreateCache(name string, policy cache.EvictionPolicy, capacity int) error
        func (cs *CacheService) GetCache(name string) (cache.Cache, error)
        func (cs *CacheService) DeleteCache(name string) error
        func (cs *CacheService) ListCaches() []string
        func (cs *CacheService) ViewStatus()
    }

    lruCache ..|> Cache : implements
    lfuCache ..|> Cache : implements
    lruCache *-- lruNode : doubly-linked list
    lfuCache *-- lfuNode : uses nodes
    lfuCache *-- freqList : frequency buckets
    freqList *-- lfuNode : contains
    CacheService o-- Cache : manages`)
}
