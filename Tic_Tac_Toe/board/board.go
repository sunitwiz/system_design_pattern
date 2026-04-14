package board

import (
	"fmt"
	"strings"
	"sync"
	"tic_tac_toe/player"
)

type CellState int

const (
	Empty CellState = iota
	SymbolX
	SymbolO
)

func (c CellState) String() string {
	switch c {
	case Empty:
		return " "
	case SymbolX:
		return "X"
	case SymbolO:
		return "O"
	default:
		return "?"
	}
}

type Board struct {
	Size       int
	Grid       [][]CellState
	MovesCount int
	mu         sync.Mutex
}

func NewBoard(size int) *Board {
	if size <= 0 {
		size = 3
	}
	grid := make([][]CellState, size)
	for i := range grid {
		grid[i] = make([]CellState, size)
	}
	return &Board{
		Size: size,
		Grid: grid,
	}
}

func symbolToCellState(symbol player.Symbol) CellState {
	switch symbol {
	case player.X:
		return SymbolX
	case player.O:
		return SymbolO
	default:
		return Empty
	}
}

func (b *Board) PlaceMove(row, col int, symbol player.Symbol) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if row < 0 || row >= b.Size || col < 0 || col >= b.Size {
		return fmt.Errorf("position (%d, %d) is out of bounds", row, col)
	}
	if b.Grid[row][col] != Empty {
		return fmt.Errorf("position (%d, %d) is already occupied", row, col)
	}

	b.Grid[row][col] = symbolToCellState(symbol)
	b.MovesCount++
	return nil
}

func (b *Board) IsFull() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.MovesCount >= b.Size*b.Size
}

func (b *Board) CheckWin(symbol player.Symbol) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	target := symbolToCellState(symbol)

	for i := 0; i < b.Size; i++ {
		rowWin := true
		for j := 0; j < b.Size; j++ {
			if b.Grid[i][j] != target {
				rowWin = false
				break
			}
		}
		if rowWin {
			return true
		}
	}

	for j := 0; j < b.Size; j++ {
		colWin := true
		for i := 0; i < b.Size; i++ {
			if b.Grid[i][j] != target {
				colWin = false
				break
			}
		}
		if colWin {
			return true
		}
	}

	diagWin := true
	for i := 0; i < b.Size; i++ {
		if b.Grid[i][i] != target {
			diagWin = false
			break
		}
	}
	if diagWin {
		return true
	}

	antiDiagWin := true
	for i := 0; i < b.Size; i++ {
		if b.Grid[i][b.Size-1-i] != target {
			antiDiagWin = false
			break
		}
	}
	return antiDiagWin
}

func (b *Board) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i := range b.Grid {
		for j := range b.Grid[i] {
			b.Grid[i][j] = Empty
		}
	}
	b.MovesCount = 0
}

func (b *Board) Display() {
	b.mu.Lock()
	defer b.mu.Unlock()

	separator := "  " + strings.Repeat("-", b.Size*4-1)
	for i := 0; i < b.Size; i++ {
		row := "   "
		for j := 0; j < b.Size; j++ {
			if j > 0 {
				row += " | "
			}
			row += b.Grid[i][j].String()
		}
		fmt.Println(row)
		if i < b.Size-1 {
			fmt.Println(separator)
		}
	}
	fmt.Println()
}
