package piece

type rook struct {
	basePiece
}

func (r *rook) IsValidMove(from, to Position, board [8][8]Piece) bool {
	if !isInBounds(to) {
		return false
	}

	if from.Row != to.Row && from.Col != to.Col {
		return false
	}

	if from.Row == to.Row && from.Col == to.Col {
		return false
	}

	if !isPathClear(from, to, board) {
		return false
	}

	if board[to.Row][to.Col] != nil && board[to.Row][to.Col].GetColor() == r.color {
		return false
	}

	return true
}
