package game

import (
	"fmt"
	"sync"
	"tic_tac_toe/board"
	"tic_tac_toe/player"
)

type Game struct {
	mu          sync.Mutex
	Board       *board.Board
	Players     [2]*player.Player
	CurrentTurn int
	Status      GameStatus
	Winner      *player.Player
}

var (
	instance *Game
	once     sync.Once
)

func GetInstance() *Game {
	once.Do(func() {
		instance = &Game{}
	})
	return instance
}

func ResetInstance() {
	once = sync.Once{}
	instance = nil
}

func (g *Game) StartGame(player1Name, player2Name string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Board = board.NewBoard(3)
	g.Players = [2]*player.Player{
		player.NewPlayer(player1Name, player.X),
		player.NewPlayer(player2Name, player.O),
	}
	g.CurrentTurn = 0
	g.Status = InProgress
	g.Winner = nil
}

func (g *Game) MakeMove(row, col int) (GameStatus, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status != InProgress {
		return g.Status, fmt.Errorf("game is not in progress (status: %s)", g.Status)
	}

	currentPlayer := g.Players[g.CurrentTurn]

	if err := g.Board.PlaceMove(row, col, currentPlayer.Symbol); err != nil {
		return g.Status, err
	}

	if g.Board.CheckWin(currentPlayer.Symbol) {
		g.Status = Won
		g.Winner = currentPlayer
		return g.Status, nil
	}

	if g.Board.IsFull() {
		g.Status = Draw
		return g.Status, nil
	}

	g.CurrentTurn = 1 - g.CurrentTurn
	return g.Status, nil
}

func (g *Game) GetCurrentPlayer() *player.Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.Players[g.CurrentTurn]
}

func (g *Game) GetWinner() *player.Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.Winner
}

func (g *Game) GetStatus() GameStatus {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.Status
}

func (g *Game) ResetGame() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Board = nil
	g.Players = [2]*player.Player{}
	g.CurrentTurn = 0
	g.Status = InProgress
	g.Winner = nil
}

func (g *Game) DisplayBoard() {
	g.Board.Display()
}
