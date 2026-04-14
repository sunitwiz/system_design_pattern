# Chess — Low-Level Design (Go)

## Problem Statement

Design a Chess game that supports two players on a standard 8×8 board. The system should model all six piece types with their unique movement rules, validate moves, detect captures, track check/checkmate/stalemate, and maintain a move history. It must be thread-safe for concurrent access.

---

## Core Entities and Relationships

```
Game (Singleton)
│
├── has ──► Board (8×8 grid)
│            └── Grid [8][8]Piece (interface)
│                        │
│                        ├── King     (1 square any direction)
│                        ├── Queen    (straight + diagonal, path clear)
│                        ├── Rook     (straight lines, path clear)
│                        ├── Bishop   (diagonals, path clear)
│                        ├── Knight   (L-shape, can jump)
│                        └── Pawn     (forward, capture diagonal)
│
├── has ──► Players [2]*Player
│            ├── Player 1 (White, moves first)
│            └── Player 2 (Black)
│
├── tracks ──► MoveHistory []*Move
│
└── tracks ──► GameStatus (InProgress | Check | Checkmate | Stalemate)
```

---

## Class Diagram

### Position

```
Position
├── Row  int   (0–7, where 0 = rank 1)
└── Col  int   (0–7, where 0 = file a)

String() → chess notation (e.g., "e4")
```

---

### Piece (Interface + Strategy)

```
Color enum:     White | Black
PieceType enum: King | Queen | Rook | Bishop | Knight | Pawn

<<interface>> Piece
├── GetType()     → PieceType
├── GetColor()    → Color
├── GetPosition() → Position
├── SetPosition(Position)
├── IsValidMove(from, to, board) → bool
└── String()      → string

Concrete types (unexported — created via factory):
├── king       ├── queen     ├── rook
├── bishop     ├── knight    └── pawn

Factory: NewPiece(pieceType, color, position) → (Piece, error)
```

**Why a strategy?** Each piece type has fundamentally different movement rules. The `Piece` interface lets the board and game logic work uniformly while each type encapsulates its own validation.

---

### Movement Rules

| Piece | Movement | Special Rules |
|-------|----------|---------------|
| **King** | 1 square in any direction | Cannot move into check |
| **Queen** | Any straight line or diagonal | Path must be clear |
| **Rook** | Horizontal or vertical lines | Path must be clear |
| **Bishop** | Diagonal lines | Path must be clear |
| **Knight** | L-shape (2+1 squares) | Can jump over pieces |
| **Pawn** | 1 forward, 2 from start row | Captures diagonally; direction depends on color |

```
Movement patterns:

King:    Queen:   Rook:    Bishop:  Knight:     Pawn (White):
 xxx     \|/      |       \  /     . x . x .     . x .
 xKx     -Q-     -R-      .B.     x . . . x     . P .
 xxx     /|\      |       /  \    . . . . .     . . .
                                   x . . . x
                                   . x . x .
```

---

### Board

```
Board
├── Grid  [8][8]Piece
├── mu    sync.Mutex

Methods:
├── MovePiece(from, to)           → (capturedPiece, error)
├── GetPiece(pos)                 → Piece
├── IsSquareUnderAttack(pos, byColor) → bool
├── FindKing(color)               → Position
└── Display()
```

Initial board setup:
```
Row 7 (rank 8): r  n  b  q  k  b  n  r   ← Black pieces
Row 6 (rank 7): p  p  p  p  p  p  p  p   ← Black pawns
Row 5–2:        .  .  .  .  .  .  .  .   ← Empty
Row 1 (rank 2): P  P  P  P  P  P  P  P   ← White pawns
Row 0 (rank 1): R  N  B  Q  K  B  N  R   ← White pieces
```

---

### Move

```
Move
├── From           Position
├── To             Position
├── PieceMoved     Piece
├── PieceCaptured  Piece (nil if no capture)
└── MoveNumber     int

String() → "1. e2→e4 (Pawn)" or "3. e4→d5 (Pawn) captures Pawn"
```

---

### Player

```
Player
├── Name   string
└── Color  Color
```

---

### Game (Singleton + Orchestrator)

```
GameStatus enum: InProgress | Check | Checkmate | Stalemate

Game
├── mu          sync.Mutex
├── Board       *Board
├── Players     [2]*Player
├── CurrentTurn int
├── Status      GameStatus
├── MoveHistory []*Move
├── moveCounter int

Methods (GameOperations interface):
├── Initialize(p1Name, p2Name)
├── MakeMove(from, to)           → error
├── GetCurrentPlayer()           → *Player
├── GetStatus()                  → GameStatus
├── GetMoveHistory()             → []*Move
├── DisplayBoard()
└── ViewStatus()
```

---

## Flows

### Move Flow

```
Player calls MakeMove(from, to)
  │
  ├── Lock mutex
  ├── Validate game status (InProgress or Check)
  ├── Validate piece at 'from' exists
  ├── Validate piece belongs to current player
  │
  ├── Board.MovePiece(from, to)
  │     ├── piece.IsValidMove(from, to, grid)
  │     │     │
  │     │     ├── King:   abs(Δrow) ≤ 1, abs(Δcol) ≤ 1
  │     │     ├── Queen:  straight or diagonal + path clear
  │     │     ├── Rook:   same row or same col + path clear
  │     │     ├── Bishop: abs(Δrow) == abs(Δcol) + path clear
  │     │     ├── Knight: (Δrow,Δcol) ∈ {(2,1),(1,2)} permutations
  │     │     └── Pawn:   1 fwd / 2 fwd from start / diagonal capture
  │     │
  │     ├── Handle capture (return captured piece)
  │     └── Update grid
  │
  ├── Record move in history
  │
  ├── Check opponent's king status:
  │     ├── isInCheck(opponentColor)?
  │     │     └── Board.IsSquareUnderAttack(kingPos, currentColor)
  │     │
  │     ├── In check + no legal moves → CHECKMATE
  │     ├── In check + has legal moves → CHECK
  │     ├── Not in check + no legal moves → STALEMATE
  │     └── Otherwise → IN PROGRESS
  │
  ├── Switch turn
  └── Unlock mutex
```

### Check / Checkmate Detection

```
isInCheck(color):
  1. Find king position for given color
  2. Check if any opponent piece can attack that square
  3. Return true if under attack

hasLegalMoves(color):
  For each piece of given color:
    For each square on the board:
      If piece.IsValidMove(from, to, grid):
        Simulate move
        If own king is NOT in check after move → legal move exists
        Undo simulation
  Return false (no legal moves found)

Checkmate = isInCheck(color) && !hasLegalMoves(color)
Stalemate = !isInCheck(color) && !hasLegalMoves(color)
```

---

## Design Patterns Used

| Pattern | Where | Why |
|---------|-------|-----|
| **Singleton** | `Game` via `sync.Once` | One game instance system-wide; safe concurrent initialization |
| **Strategy** | `Piece` interface with 6 implementations | Each piece has unique movement rules; board/game logic stays uniform |
| **Factory** | `NewPiece()` | Centralized creation; hides concrete types behind the `Piece` interface |
| **Interface Segregation** | `GameOperations` interface | Clean contract for game operations; compile-time type safety |

---

## Project Structure

```
Chess/
├── main.go                     # Demo driver with 4 scenarios
├── piece/
│   ├── piece.go                # Piece interface, Color/PieceType/Position, factory
│   ├── king.go                 # King movement (1 square any direction)
│   ├── queen.go                # Queen movement (straight + diagonal)
│   ├── rook.go                 # Rook movement (straight lines)
│   ├── bishop.go               # Bishop movement (diagonals)
│   ├── knight.go               # Knight movement (L-shape)
│   └── pawn.go                 # Pawn movement (forward, diagonal capture)
├── board/
│   └── board.go                # 8×8 board, move execution, attack detection
├── player/
│   └── player.go               # Player with name and color
├── move/
│   └── move.go                 # Move record (from, to, captured)
└── game/
    ├── game_interface.go       # GameStatus enum, GameOperations interface
    └── game.go                 # Singleton orchestrator, check/checkmate/stalemate
```

---

## How to Run

```bash
cd Chess
go build ./...
go run main.go
```

## Sample Output

```
=== Chess Game System Demo ===

--- Scenario 1: Opening Moves ---

  a  b  c  d  e  f  g  h
8 r  n  b  q  k  b  n  r  8
7 p  p  p  p  p  p  p  p  7
6 .  .  .  .  .  .  .  .  6
5 .  .  .  .  .  .  .  .  5
4 .  .  .  .  .  .  .  .  4
3 .  .  .  .  .  .  .  .  3
2 P  P  P  P  P  P  P  P  2
1 R  N  B  Q  K  B  N  R  1
  a  b  c  d  e  f  g  h

  Move: e2→e4
  Move: e7→e5
  Move: Nf3 (g1→f3)
  Move: Nc6 (b8→c6)

  Move History:
    1. e2→e4 (Pawn)
    2. e7→e5 (Pawn)
    3. g1→f3 (Knight)
    4. b8→c6 (Knight)

--- Scenario 2: Capture ---

  Move: e4×d5 (White pawn captures Black pawn)

  Move History:
    3. e4→d5 (Pawn) captures Pawn

--- Scenario 3: Invalid Moves ---

  Error: cannot move opponent's piece at e7
  Error: invalid move for Rook from a1 to c3
  Error: invalid move for Pawn from e2 to e5
  Error: no piece at d4

--- Scenario 4: Scholar's Mate (4-move Checkmate) ---

  1. e4   e5
  2. Bc4  Nc6
  3. Qh5  Nf6
  4. Qxf7#

  a  b  c  d  e  f  g  h
8 r  .  b  q  k  b  .  r  8
7 p  p  p  p  .  Q  p  p  7
6 .  .  n  .  .  n  .  .  6
5 .  .  .  .  p  .  .  .  5
4 .  .  B  .  P  .  .  .  4
3 .  .  .  .  .  .  .  .  3
2 P  P  P  P  .  P  P  P  2
1 R  N  B  .  K  .  N  R  1
  a  b  c  d  e  f  g  h

  ♚ CHECKMATE! Alice (White) wins! ♚
```

## Thread Safety

Both `Board` and `Game` protect shared state with `sync.Mutex`. The `Game` locks during the entire move lifecycle — piece validation, move execution, capture handling, check/checkmate detection, and turn switching. The `Board` has its own mutex for grid operations. The singleton is initialized via `sync.Once`. This makes the system safe for concurrent goroutine access without external synchronization.
