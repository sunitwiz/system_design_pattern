package piece

import "fmt"

type Color int

const (
	White Color = iota
	Black
)

func (c Color) String() string {
	switch c {
	case White:
		return "White"
	case Black:
		return "Black"
	default:
		return "Unknown"
	}
}

type PieceType int

const (
	King PieceType = iota
	Queen
	Rook
	Bishop
	Knight
	Pawn
)

func (pt PieceType) String() string {
	switch pt {
	case King:
		return "King"
	case Queen:
		return "Queen"
	case Rook:
		return "Rook"
	case Bishop:
		return "Bishop"
	case Knight:
		return "Knight"
	case Pawn:
		return "Pawn"
	default:
		return "Unknown"
	}
}

type Position struct {
	Row int
	Col int
}

func (p Position) String() string {
	return fmt.Sprintf("%c%d", 'a'+rune(p.Col), p.Row+1)
}

type Piece interface {
	GetType() PieceType
	GetColor() Color
	GetPosition() Position
	SetPosition(Position)
	IsValidMove(from, to Position, board [8][8]Piece) bool
	String() string
}

type basePiece struct {
	pieceType PieceType
	color     Color
	position  Position
}

func (b *basePiece) GetType() PieceType    { return b.pieceType }
func (b *basePiece) GetColor() Color       { return b.color }
func (b *basePiece) GetPosition() Position { return b.position }
func (b *basePiece) SetPosition(p Position) { b.position = p }

func (b *basePiece) String() string {
	return fmt.Sprintf("%s %s at %s", b.color, b.pieceType, b.position)
}

func NewPiece(pieceType PieceType, color Color, position Position) (Piece, error) {
	switch pieceType {
	case King:
		return &king{basePiece{pieceType: pieceType, color: color, position: position}}, nil
	case Queen:
		return &queen{basePiece{pieceType: pieceType, color: color, position: position}}, nil
	case Rook:
		return &rook{basePiece{pieceType: pieceType, color: color, position: position}}, nil
	case Bishop:
		return &bishop{basePiece{pieceType: pieceType, color: color, position: position}}, nil
	case Knight:
		return &knight{basePiece{pieceType: pieceType, color: color, position: position}}, nil
	case Pawn:
		return &pawn{basePiece{pieceType: pieceType, color: color, position: position}}, nil
	default:
		return nil, fmt.Errorf("unknown piece type: %d", pieceType)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func isPathClear(from, to Position, board [8][8]Piece) bool {
	rowDir := sign(to.Row - from.Row)
	colDir := sign(to.Col - from.Col)

	row := from.Row + rowDir
	col := from.Col + colDir

	for row != to.Row || col != to.Col {
		if board[row][col] != nil {
			return false
		}
		row += rowDir
		col += colDir
	}

	return true
}

func isInBounds(pos Position) bool {
	return pos.Row >= 0 && pos.Row <= 7 && pos.Col >= 0 && pos.Col <= 7
}
