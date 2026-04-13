package cache

import (
	"fmt"
	"strings"
	"sync"
)

type lruNode struct {
	key   string
	value any
	prev  *lruNode
	next  *lruNode
}

type lruCache struct {
	mu       sync.RWMutex
	capacity int
	nodeMap  map[string]*lruNode
	head     *lruNode
	tail     *lruNode
}

var _ Cache = (*lruCache)(nil)

func newLRUCache(capacity int) *lruCache {
	head := &lruNode{}
	tail := &lruNode{}
	head.next = tail
	tail.prev = head
	return &lruCache{
		capacity: capacity,
		nodeMap:  make(map[string]*lruNode),
		head:     head,
		tail:     tail,
	}
}

func (c *lruCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.nodeMap[key]
	if !exists {
		return nil, false
	}
	c.moveToFront(node)
	return node.value, true
}

func (c *lruCache) Put(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.nodeMap[key]; exists {
		node.value = value
		c.moveToFront(node)
		return
	}

	if len(c.nodeMap) >= c.capacity {
		c.evict()
	}

	node := &lruNode{key: key, value: value}
	c.addToFront(node)
	c.nodeMap[key] = node
}

func (c *lruCache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.nodeMap[key]
	if !exists {
		return false
	}
	c.removeNode(node)
	delete(c.nodeMap, key)
	return true
}

func (c *lruCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.nodeMap)
}

func (c *lruCache) Capacity() int {
	return c.capacity
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.head.next = c.tail
	c.tail.prev = c.head
	c.nodeMap = make(map[string]*lruNode)
}

func (c *lruCache) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var items []string
	node := c.head.next
	for node != c.tail {
		items = append(items, fmt.Sprintf("%s:%v", node.key, node.value))
		node = node.next
	}
	return fmt.Sprintf("LRU[%d/%d] MRU→LRU: [%s]", len(c.nodeMap), c.capacity, strings.Join(items, " → "))
}

func (c *lruCache) addToFront(node *lruNode) {
	node.prev = c.head
	node.next = c.head.next
	c.head.next.prev = node
	c.head.next = node
}

func (c *lruCache) removeNode(node *lruNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *lruCache) moveToFront(node *lruNode) {
	c.removeNode(node)
	c.addToFront(node)
}

func (c *lruCache) evict() {
	victim := c.tail.prev
	if victim == c.head {
		return
	}
	c.removeNode(victim)
	delete(c.nodeMap, victim.key)
}
