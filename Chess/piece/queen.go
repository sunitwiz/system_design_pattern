package piece

type queen struct {
	basePiece
}

func (q *queen) IsValidMove(from, to Position, board [8][8]Piece) bool {
	if !isInBounds(to) {
		return false
	}

	rowDiff := abs(to.Row - from.Row)
	colDiff := abs(to.Col - from.Col)

	isStraight := from.Row == to.Row || from.Col == to.Col
	isDiagonal := rowDiff == colDiff && rowDiff > 0

	if !isStraight && !isDiagonal {
		return false
	}

	if from.Row == to.Row && from.Col == to.Col {
		return false
	}

	if !isPathClear(from, to, board) {
		return false
	}

	if board[to.Row][to.Col] != nil && board[to.Row][to.Col].GetColor() == q.color {
		return false
	}

	return true
}
