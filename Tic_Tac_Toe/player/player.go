package player

type Symbol int

const (
	X Symbol = iota
	O
)

func (s Symbol) String() string {
	switch s {
	case X:
		return "X"
	case O:
		return "O"
	default:
		return "?"
	}
}

type Player struct {
	Name   string
	Symbol Symbol
}

func NewPlayer(name string, symbol Symbol) *Player {
	return &Player{
		Name:   name,
		Symbol: symbol,
	}
}
