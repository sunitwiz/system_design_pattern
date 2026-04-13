# BookMyShow — Low-Level Design (Go)

## Problem Statement

Design an online movie ticket booking system (like BookMyShow) that allows users to search movies, browse showtimes, select seats, and book tickets. Admins can manage theatres, screens, and shows. The system should handle concurrent bookings safely and calculate prices based on seat type and day-of-week pricing strategy.

---

## Core Entities and Relationships

```
BookMyShow (Singleton)
│
├── manages many ──► Theatre
│                     │
│                     └── has many ──► Screen
│                                       │
│                                       └── has many ──► Seat
│
├── manages many ──► Movie
│
├── manages many ──► Show
│                     ├── references ──► Movie
│                     ├── references ──► Screen
│                     └── tracks ──► SeatAvailability
│
└── tracks many ──► Booking
                      ├── references ──► Show
                      ├── references ──► []Seat
                      └── uses ──► PricingStrategy (interface)
```

---

## Class Diagram

### Movie

```
Movie
├── ID        string
├── Title     string
├── Duration  int (minutes)
├── Genre     Genre (Action | Comedy | Drama | Horror | SciFi)
├── Rating    float64
└── String()  → string

Factory: NewMovie(id, title, duration, genre, rating) → *Movie
```

`Genre` is an enum: `Action | Comedy | Drama | Horror | SciFi`

---

### Theatre / Screen / Seat

```
Theatre
├── ID       string
├── Name     string
├── City     string
├── Screens  []*Screen
├── AddScreen(screen)
└── String() → string

Screen
├── ID            int
├── ScreenNumber  int
├── Seats         []*Seat
└── GetSeatsByType(SeatType) → []*Seat

Seat
├── ID      int
├── Row     string
├── Number  int
├── Type    SeatType (Regular | Premium | VIP)
└── String() → string
```

`SeatType` pricing tiers:

| Seat Type | Base Price |
|-----------|-----------|
| Regular   | $10.00    |
| Premium   | $15.00    |
| VIP       | $25.00    |

---

### Show

```
Show
├── ID                string
├── Movie             *Movie
├── Screen            *Screen
├── TheatreID         string
├── StartTime         time.Time
├── SeatAvailability  map[int]bool
├── GetAvailableSeats()          → []*Seat
├── BookSeats(seatIDs []int)     → error
├── CancelSeats(seatIDs []int)   → error
└── String()                     → string
```

`SeatAvailability` maps seat IDs to availability (`true` = available). Initialized from the Screen's seat list when a Show is created. `BookSeats` validates all requested seats are available before marking any as booked (atomic check-then-book to prevent partial bookings).

---

### Booking

```
Booking
├── ID           string (e.g. "BK-0001")
├── UserName     string
├── Show         *Show
├── Seats        []*Seat
├── TotalAmount  float64
├── Status       BookingStatus (Pending | Confirmed | Cancelled)
├── BookingTime  time.Time
├── CalculateTotal()  → float64
├── Cancel()          → error
└── String()          → string
```

A booking is created with `Status = Confirmed` and `BookingTime = now`. `CalculateTotal()` picks the pricing strategy based on the show's day-of-week, then sums `strategy.CalculatePrice(seat.Type)` for each seat. `Cancel()` flips the status and releases the seats back to the show.

---

### PricingStrategy (Strategy Pattern)

```
<<interface>> PricingStrategy
├── CalculatePrice(SeatType) → float64
└── GetMultiplier()          → float64

Concrete strategies:
├── regularPricing  → 1.0× multiplier (weekday)
├── weekendPricing  → 1.5× multiplier (Saturday/Sunday)
└── premiumPricing  → 2.0× multiplier (special events)

Factory: GetPricingStrategy(isWeekend bool) → PricingStrategy
```

Pricing formula: `basePrice(seatType) × multiplier`

**Why a strategy?** Pricing logic is decoupled from the booking. Adding a new pricing tier (e.g. holiday pricing) means adding a strategy — no changes to `Booking`.

---

### BookMyShow (Singleton + Orchestrator)

```
BookMyShow
├── mu              sync.Mutex
├── Theatres        map[string]*Theatre
├── Movies          map[string]*Movie
├── Shows           map[string]*Show
├── Bookings        map[string]*Booking
├── bookingCounter  int

Interfaces it implements:
├── BookingService
│   ├── SearchMovies(title)                          → []*Movie
│   ├── GetShows(movieID, city)                      → []*Show
│   ├── BookTickets(userName, showID, seatIDs)       → *Booking, error
│   └── CancelBooking(bookingID)                     → error
│
└── AdminService
    ├── AddTheatre(theatre)
    ├── AddScreen(theatreID, screen)                 → error
    ├── AddMovie(movie)
    ├── AddShow(id, movieID, theatreID, screen, time) → *Show, error
    └── RemoveShow(showID)                           → error
```

**Singleton** — `sync.Once` guarantees exactly one instance across all goroutines.

**Mutex** — every public method locks `mu` before reading/writing shared state (`Theatres`, `Movies`, `Shows`, `Bookings`, `bookingCounter`). This prevents race conditions like duplicate booking IDs or double-booking a seat.

**Interface Segregation** — booking operations and admin operations are defined as separate interfaces (`BookingService`, `AdminService`). A user-facing app only needs `BookingService`; an admin dashboard only needs `AdminService`.

---

## Flows

### Book Ticket

```
Client calls BookTickets(userName, showID, seatIDs)
  │
  ├── Lock mutex
  ├── Lookup show by showID
  │     └── Not found? ──► return error
  │
  ├── show.BookSeats(seatIDs)
  │     ├── Validate all seats exist
  │     ├── Validate all seats are available
  │     │     └── Any unavailable? ──► return error (no seats modified)
  │     └── Mark all seats as booked
  │
  ├── Resolve Seat objects from Screen
  ├── Generate booking ID           → "BK-0001", "BK-0002", ...
  ├── Create Booking
  │     └── CalculateTotal()
  │           ├── Determine pricing (weekend vs weekday from show.StartTime)
  │           ├── Sum: strategy.CalculatePrice(seat.Type) for each seat
  │           └── Set TotalAmount
  │
  ├── Store in Bookings map
  ├── Unlock mutex
  └── Return booking
```

### Cancel Ticket

```
Client calls CancelBooking(bookingID)
  │
  ├── Lock mutex
  ├── Lookup booking by bookingID
  │     └── Not found? ──► return error
  │
  ├── booking.Cancel()
  │     ├── Already cancelled? ──► return error
  │     ├── Set status = Cancelled
  │     └── show.CancelSeats(seatIDs) → marks seats available
  │
  ├── Unlock mutex
  └── Return success
```

---

## Design Patterns Used

| Pattern | Where | Why |
|---------|-------|-----|
| **Singleton** | `BookMyShow` via `sync.Once` | One booking system instance; safe concurrent initialization |
| **Factory** | `NewMovie()`, `NewTheatre()`, `NewScreen()`, `GetPricingStrategy()` | Centralized object creation; consistent initialization |
| **Strategy** | `PricingStrategy` interface | Pricing varies by day type; new rules = new strategy, no booking changes |
| **Interface Segregation** | `BookingService` / `AdminService` | Users and admins depend only on the methods they need |

---

## Project Structure

```
Book_My_Show/
├── main.go                      # Demo driver
├── movie/movie.go               # Movie, Genre enum
├── theatre/theatre.go           # Theatre, Screen, Seat, SeatType enum
├── show/show.go                 # Show, seat availability management
├── booking/booking.go           # Booking, BookingStatus, total calculation
├── pricing/pricing.go           # PricingStrategy interface, concrete strategies
└── bookmyshow/bookmyshow.go     # Singleton orchestrator, thread-safe operations
```

---

## How to Run

```bash
cd Book_My_Show
go build ./...
go run main.go
```

## Sample Output

```
=== BookMyShow System Demo ===

--- Admin: Adding Theatre ---
  Added theatre: PVR Cinemas (Mumbai) - 0 screens

--- Admin: Adding Screens ---
  Added Screen 1 to PVR Cinemas (10 seats)
  Added Screen 2 to PVR Cinemas (8 seats)

--- Admin: Adding Movies ---
  Added movie: Avengers: Endgame (Action, 181min, 8.4★)
  Added movie: The Hangover (Comedy, 100min, 7.7★)

--- Admin: Scheduling Shows ---
  Added show: Show[SH-001] Avengers: Endgame | Screen 1 | Sat 12-Jul 18:30 | 10/10 seats available
  Added show: Show[SH-002] The Hangover | Screen 2 | Mon 14-Jul 11:00 | 8/8 seats available
  Added show: Show[SH-003] Avengers: Endgame | Screen 2 | Sun 13-Jul 14:00 | 8/8 seats available

--- User: Searching for "Avengers" ---
  Found: Avengers: Endgame (Action, 181min, 8.4★)

--- User: Shows for Avengers in Mumbai ---
  Show[SH-001] Avengers: Endgame | Screen 1 | Sat 12-Jul 18:30 | 10/10 seats available
  Show[SH-003] Avengers: Endgame | Screen 2 | Sun 13-Jul 14:00 | 8/8 seats available

--- User: Booking Weekend Show (SH-001) ---
  Booking[BK-0001] CONFIRMED | John | Avengers: Endgame | Seats: [A1(Regular), B1(Premium), C1(VIP)] | $75.00

  Pricing Breakdown (Weekend 1.5x):
    A1(Regular)    $10.00 × 1.5 = $15.00
    B1(Premium)    $15.00 × 1.5 = $22.50
    C1(VIP)        $25.00 × 1.5 = $37.50
    Total:                        $75.00

--- Edge Case: Double Booking Seat B1 ---
  Error: seat 7 is already booked

--- User: Booking Weekday Show (SH-002) ---
  Booking[BK-0002] CONFIRMED | Jane | The Hangover | Seats: [A1(Regular), B1(Premium)] | $25.00

  Pricing Breakdown (Weekday 1.0x):
    A1(Regular)    $10.00 × 1.0 = $10.00
    B1(Premium)    $15.00 × 1.0 = $15.00
    Total:                        $25.00

--- Seat Availability After Bookings ---
  Show[SH-001] Avengers: Endgame | Screen 1 | Sat 12-Jul 18:30 | 7/10 seats available
  Show[SH-003] Avengers: Endgame | Screen 2 | Sun 13-Jul 14:00 | 8/8 seats available
  Show[SH-002] The Hangover | Screen 2 | Mon 14-Jul 11:00 | 6/8 seats available

--- User: Cancelling Booking BK-0001 ---
  Cancelled: Booking[BK-0001] CANCELLED | John | Avengers: Endgame | Seats: [A1(Regular), B1(Premium), C1(VIP)] | $75.00

--- Seat Availability After Cancellation ---
  Show[SH-001] Avengers: Endgame | Screen 1 | Sat 12-Jul 18:30 | 10/10 seats available
  Show[SH-003] Avengers: Endgame | Screen 2 | Sun 13-Jul 14:00 | 8/8 seats available

--- Edge Case: Remove Show With Active Bookings ---
  Error: cannot remove show SH-002: has active bookings

--- Admin: Removing Unbooked Show ---
  Removed show: SH-003

--- Final: All Shows in Mumbai ---
  Show[SH-001] Avengers: Endgame | Screen 1 | Sat 12-Jul 18:30 | 10/10 seats available
  Show[SH-002] The Hangover | Screen 2 | Mon 14-Jul 11:00 | 6/8 seats available
```

## Thread Safety

All `BookMyShow` methods acquire `sync.Mutex` before accessing shared state. The singleton itself is initialized via `sync.Once`. Seat booking uses a validate-then-book approach — all requested seats are checked for availability before any are marked as booked, preventing partial bookings on failure. This makes the system safe for concurrent goroutine access without external synchronization.
