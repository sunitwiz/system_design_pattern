package main

import "fmt"

func printClassDiagram() {
	fmt.Println(`classDiagram

    class Color {
        <<enumeration>>
        White
        Black
        func (c Color) String() string
    }

    class PieceType {
        <<enumeration>>
        King
        Queen
        Rook
        Bishop
        Knight
        Pawn
        func (pt PieceType) String() string
    }

    class Position {
        Row int
        Col int
        func (p Position) String() string
    }

    class Piece {
        <<interface>>
        GetType() PieceType
        GetColor() Color
        GetPosition() Position
        SetPosition(Position)
        IsValidMove(from, to Position, board [8][8]Piece) bool
        String() string
    }

    class basePiece {
        pieceType PieceType
        color     Color
        position  Position
        func (b *basePiece) GetType() PieceType
        func (b *basePiece) GetColor() Color
        func (b *basePiece) GetPosition() Position
        func (b *basePiece) SetPosition(p Position)
        func (b *basePiece) String() string
    }

    class pawn {
        basePiece
        func (p *pawn) IsValidMove(from, to Position, board [8][8]Piece) bool
    }

    class rook {
        basePiece
        func (r *rook) IsValidMove(from, to Position, board [8][8]Piece) bool
    }

    class knight {
        basePiece
        func (k *knight) IsValidMove(from, to Position, board [8][8]Piece) bool
    }

    class bishop {
        basePiece
        func (b *bishop) IsValidMove(from, to Position, board [8][8]Piece) bool
    }

    class queen {
        basePiece
        func (q *queen) IsValidMove(from, to Position, board [8][8]Piece) bool
    }

    class king {
        basePiece
        func (k *king) IsValidMove(from, to Position, board [8][8]Piece) bool
    }

    class Board {
        Grid [8][8]piece.Piece
        mu   sync.Mutex
        func NewBoard() *Board
        func (b *Board) MovePiece(from, to piece.Position) (piece.Piece, error)
        func (b *Board) GetPiece(pos piece.Position) piece.Piece
        func (b *Board) IsSquareUnderAttack(pos piece.Position, byColor piece.Color) bool
        func (b *Board) FindKing(color piece.Color) piece.Position
        func (b *Board) Display()
    }

    class Player {
        Name  string
        Color piece.Color
        func NewPlayer(name string, color piece.Color) *Player
    }

    class Move {
        From          piece.Position
        To            piece.Position
        PieceMoved    piece.Piece
        PieceCaptured piece.Piece
        MoveNumber    int
        func (m *Move) String() string
    }

    class GameStatus {
        <<enumeration>>
        InProgress
        Check
        Checkmate
        Stalemate
        func (gs GameStatus) String() string
    }

    class GameOperations {
        <<interface>>
        MakeMove(from, to piece.Position) error
        GetCurrentPlayer() *player.Player
        GetStatus() GameStatus
        GetMoveHistory() []*move.Move
        DisplayBoard()
        ViewStatus()
    }

    class Game {
        mu          sync.Mutex
        Board       *board.Board
        Players     [2]*player.Player
        CurrentTurn int
        Status      GameStatus
        MoveHistory []*move.Move
        moveCounter int
        func GetInstance() *Game
        func (g *Game) Initialize(p1Name, p2Name string)
        func (g *Game) MakeMove(from, to piece.Position) error
        func (g *Game) GetCurrentPlayer() *player.Player
        func (g *Game) GetStatus() GameStatus
        func (g *Game) GetMoveHistory() []*move.Move
        func (g *Game) DisplayBoard()
        func (g *Game) ViewStatus()
    }

    basePiece --> Color
    basePiece --> PieceType
    basePiece --> Position
    pawn --|> basePiece : embeds
    rook --|> basePiece : embeds
    knight --|> basePiece : embeds
    bishop --|> basePiece : embeds
    queen --|> basePiece : embeds
    king --|> basePiece : embeds
    pawn ..|> Piece : implements
    rook ..|> Piece : implements
    knight ..|> Piece : implements
    bishop ..|> Piece : implements
    queen ..|> Piece : implements
    king ..|> Piece : implements
    Board *-- Piece : Grid contains
    Move --> Piece : PieceMoved
    Move --> Piece : PieceCaptured
    Move --> Position
    Player --> Color
    Game *-- Board : owns
    Game *-- Player : has 2
    Game *-- Move : MoveHistory
    Game --> GameStatus
    Game ..|> GameOperations : implements`)
}
