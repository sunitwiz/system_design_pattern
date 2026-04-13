# Cache (LRU / LFU) — Low-Level Design (Go)

## Problem Statement

Design an in-memory cache system supporting multiple eviction policies (LRU and LFU). The system should provide O(1) Get and Put operations, handle eviction transparently when at capacity, and manage multiple named caches through a central service. Thread safety is required for concurrent access.

---

## Core Entities and Relationships

```
CacheService (Singleton)
│
└── manages many ──► Cache (interface — Strategy Pattern)
                      │
                      ├── LRUCache implementation
                      │     ├── uses ──► doubly linked list (recency order)
                      │     └── uses ──► hash map (key → *node)
                      │
                      └── LFUCache implementation
                            ├── uses ──► hash map (key → *node)
                            ├── uses ──► frequency map (freq → doubly linked list)
                            └── tracks ──► minFrequency
```

---

## Class Diagram

### Cache (Interface + Factory)

```
<<interface>> Cache
├── Get(key)     → (any, bool)
├── Put(key, value)
├── Delete(key)  → bool
├── Size()       → int
├── Capacity()   → int
├── Clear()
└── String()     → string

Factory: NewCache(policy, capacity) → Cache
```

`EvictionPolicy` is an enum: `LRU | LFU`

The factory creates the correct implementation based on the eviction policy. Concrete types (`lruCache`, `lfuCache`) are unexported — callers interact only through the `Cache` interface.

---

### LRUCache (Least Recently Used)

```
lruCache
├── mu        sync.RWMutex
├── capacity  int
├── nodeMap   map[string]*lruNode    (O(1) lookup)
├── head      *lruNode               (sentinel — MRU end)
└── tail      *lruNode               (sentinel — LRU end)

lruNode
├── key    string
├── value  any
├── prev   *lruNode
└── next   *lruNode
```

---

### LFUCache (Least Frequently Used)

```
lfuCache
├── mu        sync.RWMutex
├── capacity  int
├── minFreq   int                         (O(1) eviction target)
├── nodeMap   map[string]*lfuNode         (O(1) lookup)
└── freqMap   map[int]*freqList           (frequency → doubly linked list)

lfuNode
├── key        string
├── value      any
├── frequency  int
├── prev       *lfuNode
└── next       *lfuNode

freqList
├── head  *lfuNode    (sentinel — MRU end within this frequency)
├── tail  *lfuNode    (sentinel — LRU end within this frequency)
└── size  int
```

---

### CacheService (Singleton + Orchestrator)

```
CacheService
├── mu      sync.Mutex
├── caches  map[string]Cache
│
├── CreateCache(name, policy, capacity)  → error
├── GetCache(name)                       → (Cache, error)
├── DeleteCache(name)                    → error
├── ListCaches()                         → []string
└── ViewStatus()
```

**Singleton** — `sync.Once` guarantees exactly one instance across all goroutines.

---

## LRU Algorithm

### How It Works

Evicts the **least recently used** item when the cache is full. Every access (Get or Put) moves the item to the front of the list, making the tail always the eviction candidate.

### Visual Diagram

```
  nodeMap (hash map)                 Doubly Linked List
  ┌───────────────┐
  │ "A" → node_A ─┼──┐          head ⇄ [C] ⇄ [B] ⇄ [A] ⇄ tail
  │ "B" → node_B ─┼──┤           ↑                      ↑
  │ "C" → node_C ─┼──┘          MRU                    LRU
  └───────────────┘            (front)              (evict here)

  Get("A"):
    head ⇄ [A] ⇄ [C] ⇄ [B] ⇄ tail     ← A moves to front

  Put("D") when full:
    Evict B (tail) → head ⇄ [D] ⇄ [A] ⇄ [C] ⇄ tail
```

### Operations

| Operation | Action |
|-----------|--------|
| **Get** | Move node to front of list, return value |
| **Put (exists)** | Update value, move to front |
| **Put (new, not full)** | Insert at front |
| **Put (new, full)** | Evict tail node, insert at front |
| **Delete** | Remove node from list and map |

---

## LFU Algorithm

### How It Works

Evicts the **least frequently used** item. Each access increments the item's frequency counter. When evicting, the item with the lowest frequency is removed. Ties are broken by LRU order within the same frequency bucket.

### Visual Diagram

```
  nodeMap                    freqMap (frequency → doubly linked list)
  ┌───────────────┐
  │ "A" → node_A ─┼──┐     freq=3: head ⇄ [A] ⇄ tail
  │ "B" → node_B ─┼──┤     freq=2: head ⇄ [C] ⇄ [B] ⇄ tail
  │ "C" → node_C ─┼──┤                    MRU         LRU
  │ "D" → node_D ─┼──┘     freq=1: head ⇄ [D] ⇄ tail  ← minFreq
  └───────────────┘

  Get("B"):
    Remove B from freq=2 list
    Increment B.frequency to 3
    Add B to front of freq=3 list
    freq=3: head ⇄ [B] ⇄ [A] ⇄ tail
    freq=2: head ⇄ [C] ⇄ tail

  Put("E") when full:
    Evict from minFreq=1 → remove D (tail of freq=1 list)
    Add E with freq=1, set minFreq=1
```

### Operations

| Operation | Action |
|-----------|--------|
| **Get** | Increment frequency, move to new freq list's front |
| **Put (exists)** | Update value, increment frequency, move to new freq list |
| **Put (new, not full)** | Insert with freq=1, set minFreq=1 |
| **Put (new, full)** | Evict tail of minFreq list, insert with freq=1 |
| **Delete** | Remove from freq list and node map |

---

## Complexity Analysis

| Operation | LRU | LFU | Notes |
|-----------|:---:|:---:|-------|
| **Get** | O(1) | O(1) | Hash map lookup + linked list pointer manipulation |
| **Put** | O(1) | O(1) | Hash map insert + linked list insert (+ eviction if full) |
| **Delete** | O(1) | O(1) | Hash map delete + linked list node removal |
| **Space** | O(n) | O(n) | n = capacity; LFU has additional freq map overhead |

Both implementations achieve O(1) for all operations by combining hash maps with doubly linked lists. The hash map provides O(1) key lookup; the doubly linked list provides O(1) insertion, deletion, and reordering.

---

## Flows

### Get Flow

```
Client calls Get(key)
  │
  ├── Lock mutex
  ├── Lookup key in nodeMap
  │     └── Not found? ──► return (nil, false)
  │
  ├── [LRU] Move node to front of list
  │   [LFU] Remove from old freq list, increment freq, add to new freq list
  │
  ├── Unlock mutex
  └── Return (value, true)
```

### Put Flow

```
Client calls Put(key, value)
  │
  ├── Lock mutex
  ├── Key exists?
  │     ├── Yes ──► Update value, treat as access (move/increment)
  │     │
  │     └── No ──► At capacity?
  │                  ├── Yes ──► Evict
  │                  │            ├── [LRU] Remove tail node (least recently used)
  │                  │            └── [LFU] Remove tail of minFreq list (least frequently used)
  │                  │
  │                  └── Insert new node
  │                       ├── [LRU] Add to front of list
  │                       └── [LFU] Add to freq=1 list, set minFreq=1
  │
  ├── Store in nodeMap
  ├── Unlock mutex
  └── Return
```

---

## Design Patterns Used

| Pattern | Where | Why |
|---------|-------|-----|
| **Singleton** | `CacheService` via `sync.Once` | One service instance system-wide; safe concurrent initialization |
| **Factory** | `NewCache(policy, capacity)` | Centralized creation; hides concrete types behind `Cache` interface |
| **Strategy** | `Cache` interface with LRU/LFU | Eviction behavior varies by policy; new policy = new implementation, no service changes |
| **Interface Segregation** | `Cache` interface | Callers depend only on Get/Put/Delete; unaware of eviction internals |

---

## Project Structure

```
Cache/
├── main.go                          # Demo driver
├── cache/cache.go                   # Cache interface, EvictionPolicy enum, factory
├── cache/lru.go                     # LRU implementation (doubly linked list + hash map)
├── cache/lfu.go                     # LFU implementation (frequency map + hash map)
└── cacheservice/cacheservice.go     # Singleton orchestrator, named cache management
```

---

## How to Run

```bash
go build ./...
go run main.go
```

## Sample Output

```
=== Cache System Demo ===

--- LRU Cache (capacity=3) ---

Putting 3 items: A=1, B=2, C=3
  State: LRU[3/3] MRU→LRU: [C:3 → B:2 → A:1]

Get B (moves B to front):
  Get(B) = 2, found=true
  State: LRU[3/3] MRU→LRU: [B:2 → C:3 → A:1]

Put D=4 (evicts least recently used):
  State: LRU[3/3] MRU→LRU: [D:4 → B:2 → C:3]

Get A (should miss — was evicted):
  Get(A) = <nil>, found=false

Update existing: Put B=20
  State: LRU[3/3] MRU→LRU: [B:20 → D:4 → C:3]

Delete C:
  State: LRU[2/3] MRU→LRU: [B:20 → D:4]

--- LFU Cache (capacity=3) ---

Putting 3 items: X=10, Y=20, Z=30
  State: LFU[3/3] freq=1{Z:30, Y:20, X:10}

Access X twice, Y once (X freq=3, Y freq=2, Z freq=1):
  State: LFU[3/3] freq=3{X:10} | freq=2{Y:20} | freq=1{Z:30}

Put W=40 (evicts Z — lowest frequency):
  State: LFU[3/3] freq=3{X:10} | freq=2{Y:20} | freq=1{W:40}

Get Z (should miss — was evicted):
  Get(Z) = <nil>, found=false

--- Eviction Comparison (LRU vs LFU) ---
Same operations on both: Put A,B,C → Get A, Get A, Get B, Get C → Put D

  Before eviction:
    LRU[3/3] MRU→LRU: [C:3 → B:2 → A:1]
    LFU[3/3] freq=3{A:1} | freq=2{C:3, B:2}

  After Put(D):
    LRU evicts least recently used  → LRU[3/3] MRU→LRU: [D:4 → C:3 → B:2]
    LFU evicts least frequently used → LFU[3/3] freq=3{A:1} | freq=2{C:3} | freq=1{D:4}

  LRU evicted A (least recently used, despite being most accessed)
  LFU evicted B (least frequently used, LRU tiebreaker among freq=2)

========================================
        CACHE SERVICE STATUS
========================================

  lfu-compare:       LFU[3/3] freq=3{A:1} | freq=2{C:3} | freq=1{D:4}
  lru-compare:       LRU[3/3] MRU→LRU: [D:4 → C:3 → B:2]
  product-cache:     LFU[3/3] freq=3{X:10} | freq=2{Y:20} | freq=1{W:40}
  session-cache:     LRU[2/3] MRU→LRU: [B:20 → D:4]

  Total Caches: 4
========================================
```

## Thread Safety

All `Cache` methods acquire `sync.RWMutex` before accessing shared state. `Get`, `Put`, `Delete`, and `Clear` use a write lock (`Lock`) since they modify internal data structures. `Size` uses a read lock (`RLock`) since it only reads the map length. The `CacheService` singleton is initialized via `sync.Once` and uses its own `sync.Mutex` for managing the named cache registry. This makes the system safe for concurrent goroutine access without external synchronization.
