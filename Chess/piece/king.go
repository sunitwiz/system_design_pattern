package piece

type king struct {
	basePiece
}

func (k *king) IsValidMove(from, to Position, board [8][8]Piece) bool {
	if !isInBounds(to) {
		return false
	}

	rowDiff := abs(to.Row - from.Row)
	colDiff := abs(to.Col - from.Col)

	if rowDiff > 1 || colDiff > 1 || (rowDiff == 0 && colDiff == 0) {
		return false
	}

	if board[to.Row][to.Col] != nil && board[to.Row][to.Col].GetColor() == k.color {
		return false
	}

	return true
}
