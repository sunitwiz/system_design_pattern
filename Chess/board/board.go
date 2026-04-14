package board

import (
	"chess/piece"
	"fmt"
	"strings"
	"sync"
)

type Board struct {
	Grid [8][8]piece.Piece
	mu   sync.Mutex
}

func NewBoard() *Board {
	b := &Board{}
	b.initializePieces()
	return b
}

func (b *Board) initializePieces() {
	backRank := []piece.PieceType{
		piece.Rook, piece.Knight, piece.Bishop, piece.Queen,
		piece.King, piece.Bishop, piece.Knight, piece.Rook,
	}

	for col, pt := range backRank {
		b.Grid[0][col], _ = piece.NewPiece(pt, piece.White, piece.Position{Row: 0, Col: col})
	}

	for col := 0; col < 8; col++ {
		b.Grid[1][col], _ = piece.NewPiece(piece.Pawn, piece.White, piece.Position{Row: 1, Col: col})
	}

	for col := 0; col < 8; col++ {
		b.Grid[6][col], _ = piece.NewPiece(piece.Pawn, piece.Black, piece.Position{Row: 6, Col: col})
	}

	for col, pt := range backRank {
		b.Grid[7][col], _ = piece.NewPiece(pt, piece.Black, piece.Position{Row: 7, Col: col})
	}
}

func (b *Board) MovePiece(from, to piece.Position) (piece.Piece, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	p := b.Grid[from.Row][from.Col]
	if p == nil {
		return nil, fmt.Errorf("no piece at %s", from)
	}

	if !p.IsValidMove(from, to, b.Grid) {
		return nil, fmt.Errorf("invalid move for %s from %s to %s", p.GetType(), from, to)
	}

	captured := b.Grid[to.Row][to.Col]
	b.Grid[to.Row][to.Col] = p
	b.Grid[from.Row][from.Col] = nil
	p.SetPosition(to)

	return captured, nil
}

func (b *Board) GetPiece(pos piece.Position) piece.Piece {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Grid[pos.Row][pos.Col]
}

func (b *Board) IsSquareUnderAttack(pos piece.Position, byColor piece.Color) bool {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := b.Grid[row][col]
			if p == nil || p.GetColor() != byColor {
				continue
			}
			from := piece.Position{Row: row, Col: col}
			if p.IsValidMove(from, pos, b.Grid) {
				return true
			}
		}
	}
	return false
}

func (b *Board) FindKing(color piece.Color) piece.Position {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := b.Grid[row][col]
			if p != nil && p.GetType() == piece.King && p.GetColor() == color {
				return piece.Position{Row: row, Col: col}
			}
		}
	}
	return piece.Position{Row: -1, Col: -1}
}

func (b *Board) Display() {
	b.mu.Lock()
	defer b.mu.Unlock()

	fmt.Println("  a  b  c  d  e  f  g  h")
	for row := 7; row >= 0; row-- {
		fmt.Printf("%d ", row+1)
		for col := 0; col < 8; col++ {
			p := b.Grid[row][col]
			if p == nil {
				fmt.Print(".  ")
			} else {
				fmt.Printf("%s  ", pieceSymbol(p))
			}
		}
		fmt.Printf("%d\n", row+1)
	}
	fmt.Println("  a  b  c  d  e  f  g  h")
}

func pieceSymbol(p piece.Piece) string {
	symbols := map[piece.PieceType]string{
		piece.King:   "K",
		piece.Queen:  "Q",
		piece.Rook:   "R",
		piece.Bishop: "B",
		piece.Knight: "N",
		piece.Pawn:   "P",
	}

	s := symbols[p.GetType()]
	if p.GetColor() == piece.Black {
		return strings.ToLower(s)
	}
	return s
}
