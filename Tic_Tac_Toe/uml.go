package main

import "fmt"

func printClassDiagram() {
	fmt.Println(`classDiagram

    class Symbol {
        <<enumeration>>
        X
        O
        func (s Symbol) String() string
    }

    class Player {
        Name   string
        Symbol Symbol
        func NewPlayer(name string, symbol Symbol) *Player
    }

    class CellState {
        <<enumeration>>
        Empty
        SymbolX
        SymbolO
        func (c CellState) String() string
    }

    class Board {
        Size       int
        Grid       [][]CellState
        MovesCount int
        mu         sync.Mutex
        func NewBoard(size int) *Board
        func (b *Board) PlaceMove(row, col int, symbol player.Symbol) error
        func (b *Board) IsFull() bool
        func (b *Board) CheckWin(symbol player.Symbol) bool
        func (b *Board) Reset()
        func (b *Board) Display()
    }

    class GameStatus {
        <<enumeration>>
        InProgress
        Won
        Draw
        func (s GameStatus) String() string
    }

    class GameOperations {
        <<interface>>
        StartGame(player1Name, player2Name string)
        MakeMove(row, col int) (GameStatus, error)
        GetCurrentPlayer() *player.Player
        GetWinner() *player.Player
        GetStatus() GameStatus
        ResetGame()
        DisplayBoard()
    }

    class Game {
        mu          sync.Mutex
        Board       *board.Board
        Players     [2]*player.Player
        CurrentTurn int
        Status      GameStatus
        Winner      *player.Player
        func GetInstance() *Game
        func (g *Game) StartGame(player1Name, player2Name string)
        func (g *Game) MakeMove(row, col int) (GameStatus, error)
        func (g *Game) GetCurrentPlayer() *player.Player
        func (g *Game) GetWinner() *player.Player
        func (g *Game) GetStatus() GameStatus
        func (g *Game) ResetGame()
        func (g *Game) DisplayBoard()
    }

    Player --> Symbol
    Board --> CellState : Grid stores
    Board ..> Symbol : converts to CellState
    Game *-- Board : owns
    Game *-- Player : has 2
    Game --> GameStatus
    Game --> Player : Winner
    Game ..|> GameOperations : implements`)
}
