# Elevator System — Low-Level Design (Go)

## Problem Statement

Design an elevator system that manages multiple elevators in a building. The system should handle floor requests from passengers, assign the optimal elevator using pluggable scheduling strategies, simulate elevator movement step-by-step, and allow admins to add/remove elevators and switch scheduling algorithms at runtime.

---

## Core Entities and Relationships

```
ElevatorController (Singleton)
│
├── has many ──► Elevator
│                 ├── CurrentFloor, Direction, Status
│                 └── Requests (sorted set of floors to visit)
│
├── uses one ──► ElevatorScheduler (Strategy interface)
│                 ├── NearestElevatorScheduler
│                 └── RoundRobinScheduler
│
└── receives ──► Request
                  ├── SourceFloor
                  ├── DestinationFloor
                  └── Direction (derived)
```

---

## Class Diagram

### Request

```
Request
├── SourceFloor       int
├── DestinationFloor  int
├── Direction         Direction (Up | Down | Idle)
├── Type              RequestType (External | Internal)
└── String()          → string

Factory: NewRequest(source, dest, type) → Request
Direction is derived automatically from source vs destination.
```

`Direction` is an enum: `Up | Down | Idle`
`RequestType` is an enum: `External | Internal`

---

### Elevator

```
Elevator
├── ID            int
├── CurrentFloor  int
├── Direction     Direction
├── Status        Status (Idle | Moving | Maintenance)
├── Requests      []int (sorted pending floors)
├── AddRequest(floor)      — adds floor, deduplicates, re-sorts
├── MoveOneStep()          — moves one floor toward next request
├── GetDirection()         → Direction
├── IsIdle()               → bool
└── String()               → string
```

`Status` is an enum: `StatusIdle | StatusMoving | StatusMaintenance`

**Movement logic:** Each call to `MoveOneStep()` increments or decrements `CurrentFloor` by 1 toward the next pending request. When the elevator arrives at a requested floor, that floor is removed from the queue. When no requests remain, the elevator goes idle.

**Request sorting:** When direction is Up, requests are sorted ascending (serve lowest first). When Down, sorted descending (serve highest first). This implements the SCAN/elevator algorithm.

---

### ElevatorScheduler (Strategy Pattern)

```
<<interface>> ElevatorScheduler
├── AssignElevator(elevators, request) → *Elevator
└── String()                           → string

Concrete strategies:
├── NearestElevatorScheduler — picks closest idle or same-direction elevator
└── RoundRobinScheduler      — cycles through elevators in rotation

Factory: NewScheduler(schedulerType) → ElevatorScheduler
```

**Nearest strategy:**
1. First pass: find idle elevators or elevators moving in the same direction toward the source floor — pick the closest
2. Fallback: if no ideal match, pick the closest elevator regardless of direction
3. Skip elevators in maintenance

**Round Robin strategy:**
- Maintains `lastIndex`, assigns to the next non-maintenance elevator in rotation
- Simple, fair distribution — does not consider proximity or direction

**Why a strategy?** Scheduling logic is decoupled from the controller. Adding a new algorithm (e.g., zone-based, load-balanced) means adding a new strategy — no changes to `ElevatorController`.

---

### ElevatorController (Singleton + Orchestrator)

```
ElevatorController
├── mu          sync.Mutex
├── Elevators   []*Elevator
├── Scheduler   ElevatorScheduler

Interfaces it implements:
├── ElevatorOperations
│   ├── RequestElevator(sourceFloor, destFloor) → *Elevator, error
│   ├── StepAll()                               — moves all elevators one step
│   └── ViewStatus()                            — prints system status
│
└── AdminOperations
    ├── AddElevator(id)
    ├── RemoveElevator(id)          → error
    └── SetScheduler(scheduler)
```

**Singleton** — `sync.Once` guarantees exactly one instance across all goroutines.

**Mutex** — every public method locks `mu` before reading/writing shared state (`Elevators`, `Scheduler`). This prevents race conditions like double-assigning an elevator or concurrent modification of the elevator list.

**Interface Segregation** — elevator operations and admin operations are defined as separate interfaces. A passenger-facing system only needs `ElevatorOperations`; a building management dashboard only needs `AdminOperations`.

---

## Flows

### Request Elevator

```
Client calls RequestElevator(sourceFloor, destFloor)
  │
  ├── Lock mutex
  ├── Create Request (direction derived from source → dest)
  ├── Scheduler.AssignElevator(elevators, request)
  │     ├── Nearest: find closest idle/same-direction elevator
  │     └── RoundRobin: pick next in rotation
  │
  ├── Assigned? ──No──► return error("no suitable elevator")
  │      │
  │     Yes
  │      │
  ├── elevator.AddRequest(sourceFloor)  → pick up passenger
  ├── elevator.AddRequest(destFloor)    → deliver passenger
  ├── Unlock mutex
  └── Return assigned elevator
```

### Step Simulation

```
Client calls StepAll()
  │
  ├── Lock mutex
  ├── For each elevator:
  │     └── elevator.MoveOneStep()
  │           ├── Move one floor in current direction
  │           ├── Arrived at requested floor? → remove from queue
  │           ├── Queue empty? → set Idle
  │           └── Otherwise → update direction toward next request
  ├── Unlock mutex
  └── Return
```

---

## Scheduling Strategy Explanations

| Strategy | How It Works | Best For |
|----------|-------------|----------|
| **Nearest** | Picks the elevator closest to the source floor that is either idle or moving in the same direction toward the source. Falls back to any closest elevator if no ideal match. | Minimizing wait time, real-world elevator systems |
| **Round Robin** | Cycles through elevators in order, skipping those in maintenance. Does not consider distance or direction. | Fair load distribution, simple systems |

---

## Design Patterns Used

| Pattern | Where | Why |
|---------|-------|-----|
| **Singleton** | `ElevatorController` via `sync.Once` | One controller instance system-wide; safe concurrent initialization |
| **Factory** | `NewScheduler()`, `NewRequest()`, `NewElevator()` | Centralized object creation; hides construction details |
| **Strategy** | `ElevatorScheduler` interface | Scheduling algorithm varies at runtime; new strategies = no controller changes |
| **Interface Segregation** | `ElevatorOperations` / `AdminOperations` | Clients depend only on the methods they need |

---

## Project Structure

```
Elevator_System/
├── main.go                        # Demo driver
├── request/request.go             # Request, Direction enum, RequestType enum
├── elevator/elevator.go           # Elevator struct, movement logic, Status enum
├── scheduler/
│   ├── scheduler.go               # ElevatorScheduler interface, SchedulerType, factory
│   ├── nearest.go                 # NearestElevatorScheduler (proximity-based)
│   └── round_robin.go             # RoundRobinScheduler (rotation-based)
└── controller/controller.go       # Singleton orchestrator, thread-safe operations
```

---

## How to Run

```bash
cd Elevator_System
go build ./...
go run main.go
```

## Sample Output

```
=== Elevator System Demo ===

--- Admin: Adding Elevators ---
  Added Elevator 1
  Added Elevator 2
  Added Elevator 3

--- Setting Nearest Scheduler ---
  Scheduler set to: NearestElevatorScheduler

--- Requesting Elevators (Nearest Scheduler) ---
  Person at floor 3 → 7: assigned to Elevator 1
  Person at floor 1 → 5: assigned to Elevator 2
  Person at floor 8 → 2: assigned to Elevator 3

--- Simulation: Stepping Elevators ---
  [Step 1]
  Elevator 1: Floor 2, Moving Up, Pending: [3 7]
  Elevator 2: Floor 1, Moving Up, Pending: [5]
  Elevator 3: Floor 1, Moving Up, Pending: [2 8]

--- Switching to Round Robin Scheduler ---
  Scheduler set to: RoundRobinScheduler

--- Edge Case: No Elevators Available ---
  Error: no elevators available

--- Edge Case: Remove Busy Elevator ---
  Error: cannot remove elevator 1: has pending requests
```

## Thread Safety

All `ElevatorController` methods acquire `sync.Mutex` before accessing shared state. The singleton itself is initialized via `sync.Once`. This makes the system safe for concurrent goroutine access without external synchronization.
