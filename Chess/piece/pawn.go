package piece

type pawn struct {
	basePiece
}

func (p *pawn) IsValidMove(from, to Position, board [8][8]Piece) bool {
	if !isInBounds(to) {
		return false
	}

	direction := 1
	startRow := 1
	if p.color == Black {
		direction = -1
		startRow = 6
	}

	rowDiff := to.Row - from.Row
	colDiff := abs(to.Col - from.Col)

	if colDiff == 0 && rowDiff == direction {
		return board[to.Row][to.Col] == nil
	}

	if colDiff == 0 && rowDiff == 2*direction && from.Row == startRow {
		middleRow := from.Row + direction
		return board[middleRow][from.Col] == nil && board[to.Row][to.Col] == nil
	}

	if colDiff == 1 && rowDiff == direction {
		return board[to.Row][to.Col] != nil && board[to.Row][to.Col].GetColor() != p.color
	}

	return false
}
