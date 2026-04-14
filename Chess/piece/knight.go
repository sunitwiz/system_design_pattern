package piece

type knight struct {
	basePiece
}

func (k *knight) IsValidMove(from, to Position, board [8][8]Piece) bool {
	if !isInBounds(to) {
		return false
	}

	rowDiff := abs(to.Row - from.Row)
	colDiff := abs(to.Col - from.Col)

	if !((rowDiff == 2 && colDiff == 1) || (rowDiff == 1 && colDiff == 2)) {
		return false
	}

	if board[to.Row][to.Col] != nil && board[to.Row][to.Col].GetColor() == k.color {
		return false
	}

	return true
}
