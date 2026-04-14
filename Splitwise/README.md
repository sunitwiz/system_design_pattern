# Splitwise — Low-Level Design (Go)

## Problem Statement

Design an expense-sharing system (like Splitwise) that allows users to split expenses among groups using different strategies — equal, exact, or percentage-based. The system should track balances between users, support settlements, and be thread-safe for concurrent access.

---

## Core Entities and Relationships

```
SplitwiseService (Singleton)
│
├── manages ──► Users (map[id]*User)
│
├── manages ──► Groups (map[id]*Group)
│                   │
│                   ├── has many ──► Members (user IDs)
│                   └── has many ──► Expenses
│
└── creates ──► Expense
                    │
                    ├── paid by ──► User (payer)
                    ├── split via ──► Split (interface / strategy)
                    │                  ├── EqualSplit
                    │                  ├── ExactSplit
                    │                  └── PercentSplit
                    └── among ──► Participants (user IDs)
```

---

## Class Diagram

### User

```
User
├── ID        string
├── Name      string
├── Email     string
├── Balances  map[string]float64    (userID → net amount)
│
├── GetBalance(userID) → float64
└── UpdateBalance(userID, amount)
```

Positive balance means the other user owes this user; negative means this user owes the other.

---

### Split (Interface + Strategy)

```
<<interface>> Split
├── Calculate(totalAmount, participants, details) → (map[string]float64, error)
└── GetType() → SplitType

SplitType enum: EqualSplit | ExactSplit | PercentSplit

Concrete types (unexported — created via factory):
├── equalSplit     — divides totalAmount equally
├── exactSplit     — uses provided amounts (validates sum = total)
└── percentSplit   — uses provided percentages (validates sum = 100)

Factory: NewSplit(splitType) → Split
```

**Why a strategy?** Different expense-splitting logic varies by type. Adding a new split type means adding one struct — no changes to service logic.

---

### Expense

```
Expense
├── ID            string
├── PayerID       string
├── Amount        float64
├── SplitType     SplitType
├── Participants  []string
└── SplitDetails  map[string]float64   (computed: userID → amount owed)
```

---

### Group

```
Group
├── ID        string
├── Name      string
├── Members   []string
└── Expenses  []*Expense
│
├── AddMember(userID)
└── AddExpense(expense)
```

---

### SplitwiseService (Singleton + Orchestrator)

```
SplitwiseService
├── mu              sync.Mutex
├── users           map[string]*User
├── groups          map[string]*Group
├── expenses        []*Expense
├── userCounter     int
├── groupCounter    int
├── expenseCounter  int

Methods (SplitwiseOperations interface):
├── AddUser(name, email) → string
├── CreateGroup(name, memberIDs) → (string, error)
├── AddExpenseToGroup(groupID, payerID, amount, splitType, participants, details) → error
├── GetBalances(userID) → map[string]float64
├── GetGroupExpenses(groupID) → []*Expense
├── SettleUp(fromUserID, toUserID, amount) → error
└── ViewStatus()
```

---

## Flows

### Add Expense Flow

```
Client calls AddExpenseToGroup(groupID, payerID, amount, splitType, participants, details)
  │
  ├── Lock mutex
  ├── Validate group exists
  ├── Create Split strategy via NewSplit(splitType)
  ├── Call split.Calculate(amount, participants, details)
  │     │
  │     ├── EqualSplit:   amount / len(participants) each
  │     ├── ExactSplit:   validate sum = amount, use provided values
  │     └── PercentSplit: validate sum = 100%, compute amount × %
  │
  ├── Create Expense with computed splits
  ├── Update balances: each participant (except payer) owes payer their share
  ├── Add expense to group
  └── Unlock mutex
```

### Settlement Flow

```
Client calls SettleUp(fromUser, toUser, amount)
  │
  ├── Lock mutex
  ├── Validate both users exist
  ├── Update fromUser balance: reduce debt to toUser by amount
  ├── Update toUser balance: reduce credit from fromUser by amount
  └── Unlock mutex
```

---

## Split Algorithm Comparison

| Feature | Equal Split | Exact Split | Percent Split |
|---------|:-----------:|:-----------:|:-------------:|
| Input required | None (auto-calculated) | Exact amounts per user | Percentages per user |
| Validation | None | Sum must equal total | Sum must equal 100% |
| Use case | Dinner bills, shared rides | Hotel rooms, unequal usage | Income-proportional splits |
| Complexity | Simple | Simple | Simple |

---

## Design Patterns Used

| Pattern | Where | Why |
|---------|-------|-----|
| **Singleton** | `SplitwiseService` via `sync.Once` | One central service instance; safe concurrent initialization |
| **Strategy** | `Split` interface with 3 implementations | Expense-splitting logic varies by type; new types don't affect service |
| **Factory** | `NewSplit()` | Centralized creation; hides concrete types behind the `Split` interface |
| **Interface Segregation** | `SplitwiseOperations` interface | Clean contract for the service; compile-time type safety |

---

## Project Structure

```
Splitwise/
├── main.go                          # Demo driver
├── user/user.go                     # User struct with balance tracking
├── split/split.go                   # Split interface, SplitType enum, factory + strategies
├── expense/expense.go               # Expense struct
├── group/group.go                   # Group struct (members + expenses)
├── service/service_interface.go     # SplitwiseOperations interface
└── service/service.go               # Singleton orchestrator, thread-safe operations
```

---

## How to Run

```bash
cd Splitwise
go build ./...
go run main.go
```

## Sample Output

```
=== Splitwise System Demo ===

--- Adding Users ---
  Created: Alice(U1), Bob(U2), Charlie(U3), Diana(U4)

--- Creating Group ---
  Created group: Trip to Goa (G1)

--- Expense 1: Equal Split ---
  Alice pays ₹2000 for dinner (split equally among all 4)
  Expense added successfully

--- Expense 2: Exact Split ---
  Bob pays ₹3000 for hotel (Alice:500, Bob:1000, Charlie:800, Diana:700)
  Expense added successfully

--- Expense 3: Percent Split ---
  Charlie pays ₹1000 for taxi (Alice:40%, Bob:30%, Charlie:20%, Diana:10%)
  Expense added successfully

--- Balances After All Expenses ---
  Alice:
    gets back from U3: ₹100.00
    gets back from U4: ₹500.00
  Bob:
    gets back from U3: ₹500.00
    gets back from U4: ₹700.00
  Charlie:
    owes U1: ₹100.00
    owes U2: ₹500.00
    gets back from U4: ₹100.00
  Diana:
    owes U3: ₹100.00
    owes U1: ₹500.00
    owes U2: ₹700.00

--- Settle Up: Bob pays Alice ₹500 ---
  Settlement recorded successfully

--- Group Expenses ---
  [E1] Payer: U1 | Amount: ₹2000.00 | Type: Equal
         Splits: U1=₹500.00, U2=₹500.00, U3=₹500.00, U4=₹500.00
  [E2] Payer: U2 | Amount: ₹3000.00 | Type: Exact
         Splits: U3=₹800.00, U4=₹700.00, U1=₹500.00, U2=₹1000.00
  [E3] Payer: U3 | Amount: ₹1000.00 | Type: Percent
         Splits: U2=₹300.00, U3=₹200.00, U4=₹100.00, U1=₹400.00

--- Full Status ---
========================================
       SPLITWISE STATUS
========================================
  Users: 4
  Groups: 1
  [G1] Trip to Goa → Members: 4 | Expenses: 3
========================================

--- Edge Case: Invalid Percent Split (sum ≠ 100) ---
  Error: percentages sum to 110.00, expected 100

--- Edge Case: Non-existent Group ---
  Error: group G999 not found
```

## Thread Safety

All `SplitwiseService` methods acquire `sync.Mutex` before accessing shared state (users, groups, expenses, balances). The singleton itself is initialized via `sync.Once`. User balance updates happen within the service lock, ensuring consistent balance calculations across concurrent expense additions. This makes the system safe for concurrent goroutine access without external synchronization.
