package piece

type bishop struct {
	basePiece
}

func (b *bishop) IsValidMove(from, to Position, board [8][8]Piece) bool {
	if !isInBounds(to) {
		return false
	}

	rowDiff := abs(to.Row - from.Row)
	colDiff := abs(to.Col - from.Col)

	if rowDiff != colDiff || rowDiff == 0 {
		return false
	}

	if !isPathClear(from, to, board) {
		return false
	}

	if board[to.Row][to.Col] != nil && board[to.Row][to.Col].GetColor() == b.color {
		return false
	}

	return true
}
