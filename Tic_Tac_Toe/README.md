# Tic Tac Toe — Low-Level Design (Go)

## Problem Statement

Design a Tic Tac Toe game that supports two players on a configurable board (default 3×3). The system should handle move validation, win detection (rows, columns, diagonals), draw detection, and turn management. It must be thread-safe for concurrent access.

---

## Core Entities and Relationships

```
Game (Singleton)
│
├── has ──► Board
│            └── Grid [size][size]CellState
│
├── has ──► Players [2]*Player
│            ├── Player 1 (Symbol X)
│            └── Player 2 (Symbol O)
│
└── tracks ──► GameStatus (InProgress | Won | Draw)
               └── Winner *Player (if Won)
```

---

## Class Diagram

### Player

```
Symbol enum: X | O

Player
├── Name    string
└── Symbol  Symbol
```

---

### Board

```
CellState enum: Empty | SymbolX | SymbolO

Board
├── Size        int
├── Grid        [][]CellState
├── MovesCount  int
├── mu          sync.Mutex

Methods:
├── PlaceMove(row, col, symbol)  → error
├── IsFull()                     → bool
├── CheckWin(symbol)             → bool
├── Reset()
└── Display()
```

---

### Game (Singleton)

```
GameStatus enum: InProgress | Won | Draw

Game
├── mu          sync.Mutex
├── Board       *Board
├── Players     [2]*Player
├── CurrentTurn int
├── Status      GameStatus
├── Winner      *Player

Methods (GameOperations interface):
├── StartGame(player1Name, player2Name)
├── MakeMove(row, col)        → (GameStatus, error)
├── GetCurrentPlayer()        → *Player
├── GetWinner()               → *Player
├── GetStatus()               → GameStatus
├── ResetGame()
└── DisplayBoard()
```

---

## Flows

### Move Flow

```
Player calls MakeMove(row, col)
  │
  ├── Lock mutex
  ├── Validate game is InProgress
  ├── Get current player's symbol
  │
  ├── Board.PlaceMove(row, col, symbol)
  │     ├── Validate bounds (0 ≤ row,col < size)
  │     ├── Validate cell is Empty
  │     └── Place symbol, increment MovesCount
  │
  ├── Check for win: Board.CheckWin(symbol)
  │     ├── Check all rows
  │     ├── Check all columns
  │     ├── Check main diagonal
  │     └── Check anti-diagonal
  │
  ├── Win found? → Status = Won, Winner = current player
  ├── Board full? → Status = Draw
  ├── Otherwise → switch turn
  │
  └── Unlock mutex, return status
```

### Win Detection Algorithm

```
For a 3×3 board with symbol X:

Row check:      Column check:     Diagonal checks:
X X X  ← win   X . .             X . .     . . X
. . .           X . .             . X .     . X .
. . .           X . .  ← win     . . X ←   X . . ←

For each row:    count matching cells == size → win
For each col:    count matching cells == size → win
Main diagonal:   (0,0)(1,1)(2,2) all match → win
Anti-diagonal:   (0,2)(1,1)(2,0) all match → win
```

---

## Design Patterns Used

| Pattern | Where | Why |
|---------|-------|-----|
| **Singleton** | `Game` via `sync.Once` | One game instance system-wide; safe concurrent initialization |
| **Interface Segregation** | `GameOperations` interface | Clean contract for game operations; compile-time type safety |
| **Enum** | `Symbol`, `CellState`, `GameStatus` | Type-safe constants with `String()` methods for readable output |

---

## Project Structure

```
Tic_Tac_Toe/
├── main.go                      # Demo driver with 4 scenarios
├── player/player.go             # Symbol enum, Player struct
├── board/board.go               # Board with move/win/draw logic
├── game/game_interface.go       # GameStatus enum, GameOperations interface
└── game/game.go                 # Singleton game, turn management
```

---

## How to Run

```bash
cd Tic_Tac_Toe
go build ./...
go run main.go
```

## Sample Output

```
=== Tic Tac Toe — Low-Level Design Demo ===

--- Scenario 1: Game with a Winner ---

  Starting game: Alice (X) vs Bob (O)

  Alice (X) plays at (0, 0)
   X |   |  
  -----------
     |   |  
  -----------
     |   |  

  Bob (O) plays at (1, 0)
   X |   |  
  -----------
   O |   |  
  -----------
     |   |  

  Alice (X) plays at (0, 1)
   X | X |  
  -----------
   O |   |  
  -----------
     |   |  

  Bob (O) plays at (1, 1)
   X | X |  
  -----------
   O | O |  
  -----------
     |   |  

  Alice (X) plays at (0, 2)
   X | X | X
  -----------
   O | O |  
  -----------
     |   |  

  >> Alice (X) wins!

--- Scenario 2: Game ending in a Draw ---

  Starting game: Alice (X) vs Bob (O)

  [9 moves leading to full board with no winner]

  >> It's a draw!

--- Scenario 3: Invalid Moves ---

  Attempting to place on occupied cell (1, 1):
  Error: position (1, 1) is already occupied

  Attempting to place out of bounds (5, 5):
  Error: position (5, 5) is out of bounds

--- Scenario 4: Diagonal Win ---

  Bob (O) plays main diagonal: (0,0) → (1,1) → (2,2)

  >> Bob (O) wins!
```

## Thread Safety

Both `Board` and `Game` protect shared state with `sync.Mutex`. The `Board` locks during move placement and win/draw checks. The `Game` locks during the entire move lifecycle (place → check win → check draw → switch turn) to prevent race conditions between concurrent move attempts. The singleton is initialized via `sync.Once`.
