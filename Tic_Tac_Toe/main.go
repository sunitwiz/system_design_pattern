package main

import (
	"fmt"
	"tic_tac_toe/game"
)

func main() {
	fmt.Println("=== Tic Tac Toe — Low-Level Design Demo ===")
	fmt.Println()

	scenarioWin()
	scenarioDraw()
	scenarioInvalidMoves()
	scenarioDiagonalWin()
}

func playMove(g *game.Game, row, col int) {
	current := g.GetCurrentPlayer()
	fmt.Printf("  %s (%s) plays at (%d, %d)\n", current.Name, current.Symbol, row, col)

	status, err := g.MakeMove(row, col)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}

	g.DisplayBoard()

	switch status {
	case game.Won:
		fmt.Printf("  >> %s (%s) wins!\n\n", g.GetWinner().Name, g.GetWinner().Symbol)
	case game.Draw:
		fmt.Println("  >> It's a draw!")
		fmt.Println()
	}
}

func scenarioWin() {
	fmt.Println("--- Scenario 1: Game with a Winner ---")
	fmt.Println()

	g := game.GetInstance()
	g.StartGame("Alice", "Bob")

	fmt.Println("  Starting game: Alice (X) vs Bob (O)")
	fmt.Println()

	playMove(g, 0, 0)
	playMove(g, 1, 0)
	playMove(g, 0, 1)
	playMove(g, 1, 1)
	playMove(g, 0, 2)
}

func scenarioDraw() {
	fmt.Println("--- Scenario 2: Game ending in a Draw ---")
	fmt.Println()

	g := game.GetInstance()
	g.ResetGame()
	g.StartGame("Alice", "Bob")

	fmt.Println("  Starting game: Alice (X) vs Bob (O)")
	fmt.Println()

	playMove(g, 0, 0)
	playMove(g, 0, 1)
	playMove(g, 0, 2)
	playMove(g, 1, 0)
	playMove(g, 1, 1)
	playMove(g, 2, 0)
	playMove(g, 1, 2)
	playMove(g, 2, 2)
	playMove(g, 2, 1)
}

func scenarioInvalidMoves() {
	fmt.Println("--- Scenario 3: Invalid Moves ---")
	fmt.Println()

	g := game.GetInstance()
	g.ResetGame()
	g.StartGame("Alice", "Bob")

	fmt.Println("  Starting game: Alice (X) vs Bob (O)")
	fmt.Println()

	playMove(g, 1, 1)

	fmt.Println("  Attempting to place on occupied cell (1, 1):")
	current := g.GetCurrentPlayer()
	fmt.Printf("  %s (%s) plays at (%d, %d)\n", current.Name, current.Symbol, 1, 1)
	_, err := g.MakeMove(1, 1)
	if err != nil {
		fmt.Printf("  Error: %v\n\n", err)
	}

	fmt.Println("  Attempting to place out of bounds (5, 5):")
	current = g.GetCurrentPlayer()
	fmt.Printf("  %s (%s) plays at (%d, %d)\n", current.Name, current.Symbol, 5, 5)
	_, err = g.MakeMove(5, 5)
	if err != nil {
		fmt.Printf("  Error: %v\n\n", err)
	}
}

func scenarioDiagonalWin() {
	fmt.Println("--- Scenario 4: Diagonal Win ---")
	fmt.Println()

	g := game.GetInstance()
	g.ResetGame()
	g.StartGame("Alice", "Bob")

	fmt.Println("  Starting game: Alice (X) vs Bob (O)")
	fmt.Println("  Bob (O) aims for the main diagonal:")
	fmt.Println()

	playMove(g, 0, 1)
	playMove(g, 0, 0)
	playMove(g, 1, 0)
	playMove(g, 1, 1)
	playMove(g, 2, 0)
	playMove(g, 2, 2)
}
