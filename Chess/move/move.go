package move

import (
	"chess/piece"
	"fmt"
)

type Move struct {
	From          piece.Position
	To            piece.Position
	PieceMoved    piece.Piece
	PieceCaptured piece.Piece
	MoveNumber    int
}

func (m *Move) String() string {
	base := fmt.Sprintf("%d. %s→%s (%s)", m.MoveNumber, m.From, m.To, m.PieceMoved.GetType())
	if m.PieceCaptured != nil {
		return fmt.Sprintf("%s captures %s", base, m.PieceCaptured.GetType())
	}
	return base
}
