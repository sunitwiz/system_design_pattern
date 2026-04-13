# Rate Limiter — Low-Level Design (Go)

## Problem Statement

Design a rate limiter that controls the number of requests a client can make to API endpoints within a given time window. The system should support multiple rate-limiting algorithms (Token Bucket, Sliding Window Log, Fixed Window Counter), allow per-route configuration, and be thread-safe for concurrent access.

---

## Core Entities and Relationships

```
RateLimiterService (Singleton)
│
└── manages many ──► route → RateLimiter (interface)
                               │
                               ├── TokenBucket
                               ├── SlidingWindowLog
                               └── FixedWindowCounter
                               
                     Each configured via ──► RateLimiterConfig
```

---

## Class Diagram

### RateLimiterConfig

```
RateLimiterConfig
├── MaxRequests     int            (for fixed/sliding window)
├── WindowSize      time.Duration  (for fixed/sliding window)
├── BucketCapacity  int            (for token bucket)
└── RefillRate      float64        (for token bucket)
```

A single config struct covers all algorithms. Each algorithm reads only the fields it needs.

---

### RateLimiter (Interface + Factory)

```
<<interface>> RateLimiter
├── Allow(key string)     → bool
├── GetLimit()            → int
├── GetWindowSize()       → time.Duration
└── GetAlgorithmType()    → AlgorithmType

Concrete types (unexported — created only via factory):
├── tokenBucket
├── slidingWindow
└── fixedWindow

Factory: NewRateLimiter(algorithmType, config) → RateLimiter
```

`AlgorithmType` is an enum: `TokenBucketType | SlidingWindowType | FixedWindowType`

**Why a factory?** The concrete structs are unexported. The only way to create a rate limiter is through `NewRateLimiter()`, which validates the type and returns the correct implementation.

---

### Token Bucket

```
tokenBucket
├── capacity    int
├── refillRate  float64        (tokens per second)
├── tokens      map[string]float64     (per-client token count)
├── lastRefill  map[string]time.Time   (per-client last refill timestamp)
└── mu          sync.Mutex

Allow(key):
  1. Calculate elapsed time since last refill
  2. Add elapsed × refillRate tokens (capped at capacity)
  3. If tokens ≥ 1 → consume one token, return true
  4. Otherwise → return false
```

**How it works:**

```
Capacity = 5, RefillRate = 1 token/sec

Time 0s:  [●●●●●]  5 tokens → Request ALLOWED → [●●●●○]
Time 0s:  [●●●●○]  4 tokens → Request ALLOWED → [●●●○○]
Time 0s:  [●●●○○]  3 tokens → Request ALLOWED → [●●○○○]
Time 0s:  [●●○○○]  2 tokens → Request ALLOWED → [●○○○○]
Time 0s:  [●○○○○]  1 token  → Request ALLOWED → [○○○○○]
Time 0s:  [○○○○○]  0 tokens → Request BLOCKED
Time 3s:  [●●●○○]  3 tokens refilled → Request ALLOWED → [●●○○○]
```

**Pros:** Smooth burst handling, memory efficient
**Cons:** Does not enforce a strict window-based limit

---

### Fixed Window Counter

```
fixedWindow
├── maxRequests  int
├── windowSize   time.Duration
├── counters     map[string]*windowEntry   (per-client)
└── mu           sync.Mutex

windowEntry { windowStart time.Time, count int }

Allow(key):
  1. If no entry or current window expired → start new window, count = 1, return true
  2. If count < maxRequests → increment, return true
  3. Otherwise → return false
```

**How it works:**

```
MaxRequests = 3, Window = 60s

Window [0s ─────────── 60s]
  Req @5s   → count=1 ✓ ALLOWED
  Req @20s  → count=2 ✓ ALLOWED
  Req @45s  → count=3 ✓ ALLOWED
  Req @50s  → count=3 ✗ BLOCKED

Window [60s ─────────── 120s]    ← counter resets
  Req @65s  → count=1 ✓ ALLOWED
```

**Pros:** Simple, low memory
**Cons:** Boundary burst problem — 2x burst possible at window edges

---

### Sliding Window Log

```
slidingWindow
├── maxRequests  int
├── windowSize   time.Duration
├── requestLog   map[string][]time.Time   (per-client timestamps)
└── mu           sync.Mutex

Allow(key):
  1. Remove all timestamps older than (now - windowSize)
  2. If remaining count ≥ maxRequests → return false
  3. Otherwise → append current time, return true
```

**How it works:**

```
MaxRequests = 3, Window = 60s

Timeline:  ──────|─────────────── now ──────|
                 now-60s                    now

  @T=10s  Log: [10]           count=1  ✓ ALLOWED
  @T=30s  Log: [10, 30]      count=2  ✓ ALLOWED
  @T=50s  Log: [10, 30, 50]  count=3  ✓ ALLOWED
  @T=55s  Log: [10, 30, 50]  count=3  ✗ BLOCKED
  @T=75s  Log: [30, 50]      count=2  ✓ ALLOWED  (T=10 expired)
```

**Pros:** Most accurate, no boundary burst problem
**Cons:** Higher memory usage (stores all timestamps)

---

### RateLimiterService (Singleton + Orchestrator)

```
RateLimiterService
├── mu        sync.Mutex
├── limiters  map[string]RateLimiter   (route → limiter)

Methods:
├── RegisterRoute(route, limiter)
├── AllowRequest(route, clientID)     → bool
├── GetStats(route)                   → map[string]interface{}
├── RemoveRoute(route)                → error
└── ViewStatus()
```

**Singleton** — `sync.Once` guarantees exactly one instance across all goroutines.

**Mutex** — every public method locks `mu` before reading/writing the `limiters` map. This prevents race conditions like concurrent route registration or stats access.

---

## Flows

### Request Flow

```
Client calls AllowRequest(route, clientID)
  │
  ├── Lock mutex
  ├── Lookup limiter for route
  │     └── Not found? → return true (no rate limit)
  ├── Unlock mutex
  │
  ├── Call limiter.Allow(clientID)
  │     │
  │     ├── Token Bucket:
  │     │     ├── Refill tokens based on elapsed time
  │     │     ├── tokens ≥ 1? → consume, return true
  │     │     └── Otherwise → return false
  │     │
  │     ├── Fixed Window:
  │     │     ├── Window expired? → reset counter
  │     │     ├── count < max? → increment, return true
  │     │     └── Otherwise → return false
  │     │
  │     └── Sliding Window:
  │           ├── Evict expired timestamps
  │           ├── count < max? → add timestamp, return true
  │           └── Otherwise → return false
  │
  └── Return allowed (true/false)
```

### Route Registration Flow

```
Admin calls RegisterRoute(route, limiter)
  │
  ├── Lock mutex
  ├── Store limiter in map[route]
  ├── Unlock mutex
  └── Print confirmation
```

---

## Algorithm Comparison

| Feature | Token Bucket | Fixed Window | Sliding Window Log |
|---------|:---:|:---:|:---:|
| Accuracy | Medium | Low | High |
| Memory usage | Low | Low | High |
| Boundary burst | No | Yes | No |
| Allows bursts | Yes (up to capacity) | No | No |
| Implementation complexity | Medium | Simple | Medium |
| Best for | APIs needing burst tolerance | Simple rate caps | Strict per-window limits |

---

## Design Patterns Used

| Pattern | Where | Why |
|---------|-------|-----|
| **Singleton** | `RateLimiterService` via `sync.Once` | One service instance system-wide; safe concurrent initialization |
| **Factory** | `NewRateLimiter()` | Centralized object creation; hides concrete types behind the `RateLimiter` interface |
| **Strategy** | `RateLimiter` interface | Rate-limiting algorithm varies per route; new algorithm = new struct, no service changes |
| **Interface Segregation** | `RateLimiter` interface | Each algorithm implements the same contract; service doesn't know internals |

---

## Project Structure

```
Rate_Limiter/
├── main.go                          # Demo driver
├── config/config.go                 # RateLimiterConfig struct
├── algorithm/algorithm.go           # RateLimiter interface, AlgorithmType enum, factory
├── algorithm/token_bucket.go        # Token bucket implementation
├── algorithm/sliding_window.go      # Sliding window log implementation
├── algorithm/fixed_window.go        # Fixed window counter implementation
└── ratelimiter/ratelimiter.go       # Singleton orchestrator, thread-safe operations
```

---

## How to Run

```bash
go build ./...
go run main.go
```

## Sample Output

```
=== Rate Limiter System Demo ===

--- Registering Routes ---
  Registered route /api/search      → Token Bucket (limit: 10)
  Registered route /api/login       → Fixed Window Counter (limit: 5)
  Registered route /api/data        → Sliding Window Log (limit: 10)

========================================
      RATE LIMITER STATUS
========================================
  Route              Algorithm              Limit
  ------------------ ---------------------- -----
  /api/search        Token Bucket           10 capacity
  /api/login         Fixed Window Counter   5 / 1m0s
  /api/data          Sliding Window Log     10 / 1m0s
========================================

--- Token Bucket: /api/search (capacity=10, refill=2/sec) ---
  client-A → 12 requests: 10 allowed, 2 blocked

--- Fixed Window: /api/login (5 requests per 60s) ---
  client-B → 8 requests: 5 allowed, 3 blocked

--- Sliding Window: /api/data (10 requests per 60s) ---
  client-C → 12 requests: 10 allowed, 2 blocked

--- Multiple Clients on Same Route ---
  Client-X on /api/login:
  client-X → 6 requests: 5 allowed, 1 blocked
  Client-Y on /api/login:
  client-Y → 6 requests: 5 allowed, 1 blocked

--- Rate Limit Recovery (Token Bucket) ---
  Sending 10 requests to exhaust bucket...
  Request 11 (exhausted): allowed=false
  Waiting 3 seconds for token refill (rate=2/sec → ~6 new tokens)...
  After waiting:
    Request 1: allowed=true
    Request 2: allowed=true
    Request 3: allowed=true
    Request 4: allowed=true
    Request 5: allowed=true
    Request 6: allowed=true
    Request 7: allowed=false

--- Route Stats ---
  /api/search     algorithm=Token Bucket           limit=10  window=0s
  /api/login      algorithm=Fixed Window Counter    limit=5   window=1m0s
  /api/data       algorithm=Sliding Window Log      limit=10  window=1m0s

--- Edge Case: Unregistered Route ---
  /api/unknown → allowed=true (no limiter = always allowed)

--- Removing a Route ---
  Removed route /api/data
  /api/data after removal → allowed=true

--- Final Status ---
========================================
      RATE LIMITER STATUS
========================================
  Route              Algorithm              Limit
  ------------------ ---------------------- -----
  /api/search        Token Bucket           10 capacity
  /api/login         Fixed Window Counter   5 / 1m0s
========================================
```

## Thread Safety

All `RateLimiterService` methods acquire `sync.Mutex` before accessing the shared `limiters` map. Each algorithm implementation also has its own `sync.Mutex` protecting per-client state (tokens, counters, request logs). The singleton itself is initialized via `sync.Once`. This makes the system safe for concurrent goroutine access without external synchronization.
