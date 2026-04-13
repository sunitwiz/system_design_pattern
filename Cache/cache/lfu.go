package cache

import (
	"fmt"
	"strings"
	"sync"
)

type lfuNode struct {
	key       string
	value     any
	frequency int
	prev      *lfuNode
	next      *lfuNode
}

type freqList struct {
	head *lfuNode
	tail *lfuNode
	size int
}

func newFreqList() *freqList {
	head := &lfuNode{}
	tail := &lfuNode{}
	head.next = tail
	tail.prev = head
	return &freqList{head: head, tail: tail}
}

func (fl *freqList) addToFront(node *lfuNode) {
	node.prev = fl.head
	node.next = fl.head.next
	fl.head.next.prev = node
	fl.head.next = node
	fl.size++
}

func (fl *freqList) removeNode(node *lfuNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
	fl.size--
}

func (fl *freqList) removeTail() *lfuNode {
	if fl.size == 0 {
		return nil
	}
	victim := fl.tail.prev
	fl.removeNode(victim)
	return victim
}

func (fl *freqList) isEmpty() bool {
	return fl.size == 0
}

type lfuCache struct {
	mu       sync.RWMutex
	capacity int
	minFreq  int
	nodeMap  map[string]*lfuNode
	freqMap  map[int]*freqList
}

var _ Cache = (*lfuCache)(nil)

func newLFUCache(capacity int) *lfuCache {
	return &lfuCache{
		capacity: capacity,
		nodeMap:  make(map[string]*lfuNode),
		freqMap:  make(map[int]*freqList),
	}
}

func (c *lfuCache) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.nodeMap[key]
	if !exists {
		return nil, false
	}
	c.incrementFrequency(node)
	return node.value, true
}

func (c *lfuCache) Put(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.nodeMap[key]; exists {
		node.value = value
		c.incrementFrequency(node)
		return
	}

	if len(c.nodeMap) >= c.capacity {
		c.evict()
	}

	node := &lfuNode{key: key, value: value, frequency: 1}
	c.nodeMap[key] = node

	if c.freqMap[1] == nil {
		c.freqMap[1] = newFreqList()
	}
	c.freqMap[1].addToFront(node)
	c.minFreq = 1
}

func (c *lfuCache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	node, exists := c.nodeMap[key]
	if !exists {
		return false
	}

	fl := c.freqMap[node.frequency]
	fl.removeNode(node)
	if fl.isEmpty() {
		delete(c.freqMap, node.frequency)
	}
	delete(c.nodeMap, key)
	return true
}

func (c *lfuCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.nodeMap)
}

func (c *lfuCache) Capacity() int {
	return c.capacity
}

func (c *lfuCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.nodeMap = make(map[string]*lfuNode)
	c.freqMap = make(map[int]*freqList)
	c.minFreq = 0
}

func (c *lfuCache) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	maxFreq := 0
	for freq := range c.freqMap {
		if freq > maxFreq {
			maxFreq = freq
		}
	}

	var parts []string
	for freq := maxFreq; freq >= 1; freq-- {
		fl, exists := c.freqMap[freq]
		if !exists || fl.isEmpty() {
			continue
		}
		var items []string
		node := fl.head.next
		for node != fl.tail {
			items = append(items, fmt.Sprintf("%s:%v", node.key, node.value))
			node = node.next
		}
		parts = append(parts, fmt.Sprintf("freq=%d{%s}", freq, strings.Join(items, ", ")))
	}

	return fmt.Sprintf("LFU[%d/%d] %s", len(c.nodeMap), c.capacity, strings.Join(parts, " | "))
}

func (c *lfuCache) incrementFrequency(node *lfuNode) {
	oldFreq := node.frequency
	fl := c.freqMap[oldFreq]
	fl.removeNode(node)

	if fl.isEmpty() {
		delete(c.freqMap, oldFreq)
		if c.minFreq == oldFreq {
			c.minFreq++
		}
	}

	node.frequency++

	if c.freqMap[node.frequency] == nil {
		c.freqMap[node.frequency] = newFreqList()
	}
	c.freqMap[node.frequency].addToFront(node)
}

func (c *lfuCache) evict() {
	fl := c.freqMap[c.minFreq]
	if fl == nil {
		return
	}
	victim := fl.removeTail()
	if victim == nil {
		return
	}
	if fl.isEmpty() {
		delete(c.freqMap, c.minFreq)
	}
	delete(c.nodeMap, victim.key)
}
