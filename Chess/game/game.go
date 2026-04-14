package game

import (
	"chess/board"
	"chess/move"
	"chess/piece"
	"chess/player"
	"fmt"
	"sync"
)

type Game struct {
	mu          sync.Mutex
	Board       *board.Board
	Players     [2]*player.Player
	CurrentTurn int
	Status      GameStatus
	MoveHistory []*move.Move
	moveCounter int
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

func (g *Game) Initialize(p1Name, p2Name string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.Board = board.NewBoard()
	g.Players = [2]*player.Player{
		player.NewPlayer(p1Name, piece.White),
		player.NewPlayer(p2Name, piece.Black),
	}
	g.CurrentTurn = 0
	g.Status = InProgress
	g.MoveHistory = nil
	g.moveCounter = 0
}

func (g *Game) MakeMove(from, to piece.Position) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status == Checkmate || g.Status == Stalemate {
		return fmt.Errorf("game is over: %s", g.Status)
	}

	currentPlayer := g.Players[g.CurrentTurn]

	p := g.Board.GetPiece(from)
	if p == nil {
		return fmt.Errorf("no piece at %s", from)
	}

	if p.GetColor() != currentPlayer.Color {
		return fmt.Errorf("cannot move opponent's piece at %s", from)
	}

	captured, err := g.Board.MovePiece(from, to)
	if err != nil {
		return err
	}

	ownColor := currentPlayer.Color
	opponentColor := piece.Black
	if ownColor == piece.Black {
		opponentColor = piece.White
	}

	if g.isInCheck(ownColor) {
		g.Board.Grid[from.Row][from.Col] = p
		g.Board.Grid[to.Row][to.Col] = captured
		p.SetPosition(from)
		return fmt.Errorf("move would leave your king in check")
	}

	g.moveCounter++
	m := &move.Move{
		From:          from,
		To:            to,
		PieceMoved:    p,
		PieceCaptured: captured,
		MoveNumber:    g.moveCounter,
	}
	g.MoveHistory = append(g.MoveHistory, m)

	if g.isInCheck(opponentColor) {
		if !g.hasLegalMoves(opponentColor) {
			g.Status = Checkmate
		} else {
			g.Status = Check
		}
	} else {
		if !g.hasLegalMoves(opponentColor) {
			g.Status = Stalemate
		} else {
			g.Status = InProgress
		}
	}

	g.CurrentTurn = 1 - g.CurrentTurn

	return nil
}

func (g *Game) GetCurrentPlayer() *player.Player {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.Players[g.CurrentTurn]
}

func (g *Game) GetStatus() GameStatus {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.Status
}

func (g *Game) GetMoveHistory() []*move.Move {
	g.mu.Lock()
	defer g.mu.Unlock()
	result := make([]*move.Move, len(g.MoveHistory))
	copy(result, g.MoveHistory)
	return result
}

func (g *Game) DisplayBoard() {
	g.Board.Display()
}

func (g *Game) ViewStatus() {
	g.mu.Lock()
	defer g.mu.Unlock()

	fmt.Println("========================================")
	fmt.Println("         CHESS GAME STATUS")
	fmt.Println("========================================")
	fmt.Printf("  Status:         %s\n", g.Status)
	fmt.Printf("  Current Turn:   %s (%s)\n", g.Players[g.CurrentTurn].Name, g.Players[g.CurrentTurn].Color)
	fmt.Printf("  Total Moves:    %d\n", g.moveCounter)
	if len(g.MoveHistory) > 0 {
		fmt.Printf("  Last Move:      %s\n", g.MoveHistory[len(g.MoveHistory)-1])
	}
	fmt.Println("========================================")
}

func (g *Game) isInCheck(color piece.Color) bool {
	opponentColor := piece.White
	if color == piece.White {
		opponentColor = piece.Black
	}
	kingPos := g.Board.FindKing(color)
	return g.Board.IsSquareUnderAttack(kingPos, opponentColor)
}

func (g *Game) hasLegalMoves(color piece.Color) bool {
	opponentColor := piece.White
	if color == piece.White {
		opponentColor = piece.Black
	}

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			p := g.Board.Grid[row][col]
			if p == nil || p.GetColor() != color {
				continue
			}

			from := piece.Position{Row: row, Col: col}

			for toRow := 0; toRow < 8; toRow++ {
				for toCol := 0; toCol < 8; toCol++ {
					to := piece.Position{Row: toRow, Col: toCol}

					if !p.IsValidMove(from, to, g.Board.Grid) {
						continue
					}

					captured := g.Board.Grid[to.Row][to.Col]
					g.Board.Grid[to.Row][to.Col] = p
					g.Board.Grid[from.Row][from.Col] = nil

					kingPos := g.Board.FindKing(color)
					safe := !g.Board.IsSquareUnderAttack(kingPos, opponentColor)

					g.Board.Grid[from.Row][from.Col] = p
					g.Board.Grid[to.Row][to.Col] = captured

					if safe {
						return true
					}
				}
			}
		}
	}
	return false
}
