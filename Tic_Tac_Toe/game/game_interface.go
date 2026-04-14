package game

import (
	"tic_tac_toe/player"
)

type GameStatus int

const (
	InProgress GameStatus = iota
	Won
	Draw
)

func (s GameStatus) String() string {
	switch s {
	case InProgress:
		return "In Progress"
	case Won:
		return "Won"
	case Draw:
		return "Draw"
	default:
		return "Unknown"
	}
}

type GameOperations interface {
	StartGame(player1Name, player2Name string)
	MakeMove(row, col int) (GameStatus, error)
	GetCurrentPlayer() *player.Player
	GetWinner() *player.Player
	GetStatus() GameStatus
	ResetGame()
	DisplayBoard()
}

var _ GameOperations = (*Game)(nil)
