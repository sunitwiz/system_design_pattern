# Parking Lot — Low-Level Design (Go)

## Problem Statement

Design a multi-level parking lot that can park motorcycles, cars, and buses. The system should handle vehicle entry/exit, assign the best-fit slot, calculate time-based fees, and allow admins to manage levels and slots at runtime.

---

## Core Entities and Relationships

```
ParkingLot (Singleton)
│
├── has many ──► ParkingLevel
│                 │
│                 └── has many ──► ParkingSlot
│                                   │
│                                   └── holds one ──► Vehicle (interface)
│
└── tracks many ──► ParkingTicket
                      ├── references ──► Vehicle
                      ├── references ──► ParkingSlot
                      └── uses ──► FeeStrategy (interface)
```

---

## Class Diagram

### Vehicle (Interface + Factory)

```
<<interface>> Vehicle
├── GetType()         → VehicleType
├── GetLicensePlate() → string
└── String()          → string

Concrete types (unexported — created only via factory):
├── motorcycleVehicle { licensePlate string }
├── carVehicle        { licensePlate string }
└── busVehicle        { licensePlate string }

Factory: NewVehicle(type, plate) → Vehicle
```

`VehicleType` is an enum: `Motorcycle | Car | Bus`

**Why a factory?** The concrete structs are unexported. The only way to create a vehicle is through `NewVehicle()`, which validates the type and returns the correct implementation. This keeps the creation logic centralized.

---

### ParkingSlot

```
ParkingSlot
├── ID            int
├── Type          SlotType (MotorcycleSlot | CarSlot | BusSlot)
├── IsOccupied    bool
├── ParkedVehicle Vehicle
├── CanFit(VehicleType) → bool
├── Park(Vehicle)       → error
└── Unpark()            → Vehicle, error
```

**Slot compatibility rules (which vehicle fits where):**

| Slot Type | Motorcycle | Car | Bus |
|-----------|:---:|:---:|:---:|
| Motorcycle Slot | Y | N | N |
| Car Slot | Y | Y | N |
| Bus Slot | Y | Y | Y |

A larger slot can accept smaller vehicles, but not the other way around.

---

### ParkingLevel

```
ParkingLevel
├── LevelNumber  int
├── Slots        []*ParkingSlot
├── nextSlotID   int (auto-increment)
├── FindAvailableSlot(VehicleType) → *ParkingSlot
├── GetStatus()                    → []SlotStatus
├── AddSlot(SlotType)              → *ParkingSlot
└── RemoveSlot(slotID)             → error
```

**Slot assignment strategy (two-pass):**
1. **Exact match first** — a car gets a Car Slot before being placed in a Bus Slot
2. **Compatible fallback** — if no exact match, use any slot that can fit the vehicle

This minimizes wasted capacity (a motorcycle won't consume a bus slot if motorcycle slots are free).

---

### ParkingTicket

```
ParkingTicket
├── TicketID   string (e.g. "T-0001")
├── Vehicle    Vehicle
├── Slot       *ParkingSlot
├── LevelNum   int
├── EntryTime  time.Time
├── ExitTime   time.Time
├── Fee        float64
├── IsActive   bool
├── CalculateFee()                → float64, error
└── CalculateFeeWithTime(exitTime) → float64, error
```

A ticket is created at entry (`IsActive = true`, `EntryTime = now`). At exit, `CalculateFee()` sets `ExitTime`, computes the fee via the strategy pattern, and marks the ticket `IsActive = false`.

---

### FeeStrategy (Strategy Pattern)

```
<<interface>> FeeStrategy
├── CalculateFee(duration) → float64
└── GetRatePerHour()       → float64

Concrete strategies:
├── motorcycleFee → $1/hour
├── carFee        → $2/hour
└── busFee        → $5/hour

Factory: GetFeeStrategy(VehicleType) → FeeStrategy
```

Partial hours are rounded up (e.g., 2h10m is charged as 3h).

**Why a strategy?** Fee logic is decoupled from the ticket. Adding a new vehicle type or changing rates means adding/modifying a strategy — no changes to `ParkingTicket`.

---

### ParkingLot (Singleton + Orchestrator)

```
ParkingLot
├── mu              sync.Mutex
├── Levels          []*ParkingLevel
├── ActiveTickets   map[string]*ParkingTicket
├── ticketCounter   int

Interfaces it implements:
├── ParkingOperations
│   ├── ParkVehicle(Vehicle)      → *ParkingTicket, error
│   └── UnparkVehicle(ticketID)   → *ParkingTicket, error
│
└── AdminOperations
    ├── ViewStatus()
    ├── AddLevel(motorcycle, car, bus counts)
    ├── RemoveLevel(levelNumber)  → error
    ├── AddSlot(levelNumber, SlotType) → error
    └── RemoveSlot(levelNumber, slotID) → error
```

**Singleton** — `sync.Once` guarantees exactly one instance across all goroutines.

**Mutex** — every public method locks `mu` before reading/writing shared state (`Levels`, `ActiveTickets`, `ticketCounter`). This prevents race conditions like duplicate ticket IDs or double-booking a slot.

**Interface Segregation** — parking operations and admin operations are defined as separate interfaces (`ParkingOperations`, `AdminOperations`). A valet system only needs `ParkingOperations`; an admin dashboard only needs `AdminOperations`.

---

## Flows

### Park Vehicle

```
Client calls ParkVehicle(vehicle)
  │
  ├── Lock mutex
  ├── For each level:
  │     └── FindAvailableSlot(vehicleType)
  │           ├── Pass 1: exact-match slot
  │           └── Pass 2: any compatible slot
  │
  ├── Found slot? ──No──► return error("no available slot")
  │      │
  │     Yes
  │      │
  ├── slot.Park(vehicle)         → marks slot occupied
  ├── Generate ticket ID         → "T-0001", "T-0002", ...
  ├── Create ParkingTicket       → entry time = now
  ├── Store in ActiveTickets map
  ├── Unlock mutex
  └── Return ticket
```

### Unpark Vehicle

```
Client calls UnparkVehicle(ticketID)
  │
  ├── Lock mutex
  ├── Lookup ticket in ActiveTickets
  │     └── Not found? ──► return error
  │
  ├── slot.Unpark()              → marks slot free, clears vehicle
  ├── ticket.CalculateFee()
  │     ├── GetFeeStrategy(vehicleType)  → picks strategy
  │     ├── duration = now - entryTime
  │     ├── fee = strategy.CalculateFee(duration)
  │     └── marks ticket IsActive = false
  │
  ├── Remove ticket from ActiveTickets
  ├── Unlock mutex
  └── Return closed ticket (with fee)
```

---

## Design Patterns Used

| Pattern | Where | Why |
|---------|-------|-----|
| **Singleton** | `ParkingLot` via `sync.Once` | One parking lot instance system-wide; safe concurrent initialization |
| **Factory** | `NewVehicle()`, `GetFeeStrategy()` | Centralized object creation; hides concrete types behind interfaces |
| **Strategy** | `FeeStrategy` interface | Fee calculation varies by vehicle type; new rates = new strategy, no ticket changes |
| **Interface Segregation** | `ParkingOperations` / `AdminOperations` | Clients depend only on the methods they need |

---

## Project Structure

```
Parking_Lot/
├── main.go                    # Demo driver
├── vehicle/vehicle.go         # Vehicle interface, concrete types, factory
├── slot/slot.go               # ParkingSlot, SlotType, fit rules
├── level/level.go             # ParkingLevel, slot search, status reporting
├── fee/fee.go                 # FeeStrategy interface, per-vehicle strategies
├── ticket/ticket.go           # ParkingTicket, fee calculation
└── parkinglot/parkinglot.go   # Singleton orchestrator, thread-safe operations
```

---

## How to Run

```bash
go build ./...
go run main.go
```

## Sample Output

```
--- Admin: Adding Levels ---
  Added Level 1 (5 motorcycle, 10 car, 2 bus slots)
  Added Level 2 (3 motorcycle, 8 car, 1 bus slots)

--- Parking Vehicles ---
  Parked: Ticket[T-0001] ACTIVE | Motorcycle [MC-1001] | Level 1, Slot 1 (Motorcycle Slot)
  Parked: Ticket[T-0002] ACTIVE | Car [CAR-2001] | Level 1, Slot 6 (Car Slot)
  Parked: Ticket[T-0004] ACTIVE | Bus [BUS-3001] | Level 1, Slot 16 (Bus Slot)

--- Unparking Vehicles (simulating 3-hour stay) ---
  Motorcycle (3 hrs @ $1/hr): $3.00
  Car        (3 hrs @ $2/hr): $6.00
  Bus        (3 hrs @ $5/hr): $15.00
```

## Thread Safety

All `ParkingLot` methods acquire `sync.Mutex` before accessing shared state. The singleton itself is initialized via `sync.Once`. This makes the system safe for concurrent goroutine access without external synchronization.

---

## How to Approach Any LLD Problem

The mental model for low-level design is always the same seven steps. Here's how each one maps to this parking lot:

### 1. Requirements — "What does the system do?"

List every feature before writing code. For this problem: multi-level lot, three vehicle types, slot compatibility, ticket-based entry/exit, time-based fees, admin management.

### 2. Entities — "What are the nouns?"

Pull out every noun from the requirements — these become your structs:

`Vehicle`, `ParkingSlot`, `ParkingLevel`, `ParkingTicket`, `FeeStrategy`, `ParkingLot`

### 3. Behaviors — "What does each entity do?"

Each noun gets verbs — these become methods:

- `ParkingSlot` → `CanFit()`, `Park()`, `Unpark()`
- `ParkingLevel` → `FindAvailableSlot()`, `AddSlot()`, `RemoveSlot()`
- `ParkingTicket` → `CalculateFee()`
- `ParkingLot` → `ParkVehicle()`, `UnparkVehicle()`, `AddLevel()`, `ViewStatus()`

### 4. Relationships — "Who knows about whom?"

This defines your dependency graph:

- `ParkingLot` knows about `ParkingLevel` and `ParkingTicket`
- `ParkingLevel` knows about `ParkingSlot`
- `ParkingSlot` knows about `Vehicle`
- `ParkingTicket` knows about `Vehicle`, `ParkingSlot`, and `FeeStrategy`

Notice that dependencies flow **downward** — `Vehicle` and `ParkingSlot` don't know about `ParkingLot`. This keeps lower-level entities reusable and testable in isolation.

### 5. Variations — "Where does behavior change?"

Wherever the same operation behaves differently depending on type, introduce an **interface**:

- Fee calculation varies by vehicle → `FeeStrategy` interface (Strategy pattern)
- Vehicle creation varies by type → `Vehicle` interface + `NewVehicle()` factory
- Parking vs Admin operations are different concerns → two separate interfaces (Interface Segregation)

### 6. Concurrency — "Who shares state?"

Ask: can multiple callers access this at the same time?

- `ParkingLot` is a singleton shared across all goroutines → needs `sync.Mutex` on every method that touches `Levels`, `ActiveTickets`, or `ticketCounter`
- `ParkingSlot` and `ParkingLevel` are only accessed through `ParkingLot`'s locked methods → no separate lock needed

### 7. Flows — "Trace a full request start to finish"

Walk through the complete park and unpark flows to validate your design. If any step feels awkward or requires an entity to reach into something it shouldn't know about, your relationships are wrong — go back to step 4.

---

**Apply these seven steps to any LLD problem** (elevator system, library management, ride sharing, etc.) and you'll arrive at a clean, extensible design every time.
