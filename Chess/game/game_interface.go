package game

import (
	"chess/move"
	"chess/piece"
	"chess/player"
)

type GameStatus int

const (
	InProgress GameStatus = iota
	Check
	Checkmate
	Stalemate
)

func (gs GameStatus) String() string {
	switch gs {
	case InProgress:
		return "In Progress"
	case Check:
		return "Check"
	case Checkmate:
		return "Checkmate"
	case Stalemate:
		return "Stalemate"
	default:
		return "Unknown"
	}
}

type GameOperations interface {
	MakeMove(from, to piece.Position) error
	GetCurrentPlayer() *player.Player
	GetStatus() GameStatus
	GetMoveHistory() []*move.Move
	DisplayBoard()
	ViewStatus()
}

var _ GameOperations = (*Game)(nil)
