package player

import "chess/piece"

type Player struct {
	Name  string
	Color piece.Color
}

func NewPlayer(name string, color piece.Color) *Player {
	return &Player{
		Name:  name,
		Color: color,
	}
}
