package main

import (
	"chess/game"
	"chess/piece"
	"fmt"
)

func main() {
	fmt.Println("=== Chess Game System Demo ===")
	fmt.Println()

	scenario1()
	scenario2()
	scenario3()
	scenario4()
}

func scenario1() {
	fmt.Println("--- Scenario 1: Opening Moves ---")
	fmt.Println()

	game.ResetInstance()
	g := game.GetInstance()
	g.Initialize("Alice", "Bob")

	g.DisplayBoard()
	fmt.Println()

	moves := []struct {
		from, to piece.Position
		desc     string
	}{
		{piece.Position{Row: 1, Col: 4}, piece.Position{Row: 3, Col: 4}, "e2→e4"},
		{piece.Position{Row: 6, Col: 4}, piece.Position{Row: 4, Col: 4}, "e7→e5"},
		{piece.Position{Row: 0, Col: 6}, piece.Position{Row: 2, Col: 5}, "Nf3 (g1→f3)"},
		{piece.Position{Row: 7, Col: 1}, piece.Position{Row: 5, Col: 2}, "Nc6 (b8→c6)"},
	}

	for _, m := range moves {
		fmt.Printf("  Move: %s\n", m.desc)
		if err := g.MakeMove(m.from, m.to); err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}
		g.DisplayBoard()
		fmt.Println()
	}

	fmt.Println("  Move History:")
	for _, m := range g.GetMoveHistory() {
		fmt.Printf("    %s\n", m)
	}

	g.ViewStatus()
	fmt.Println()
}

func scenario2() {
	fmt.Println("--- Scenario 2: Capture ---")
	fmt.Println()

	game.ResetInstance()
	g := game.GetInstance()
	g.Initialize("Alice", "Bob")

	setupMoves := []struct {
		from, to piece.Position
		desc     string
	}{
		{piece.Position{Row: 1, Col: 4}, piece.Position{Row: 3, Col: 4}, "e2→e4"},
		{piece.Position{Row: 6, Col: 3}, piece.Position{Row: 4, Col: 3}, "d7→d5"},
	}

	for _, m := range setupMoves {
		fmt.Printf("  Move: %s\n", m.desc)
		if err := g.MakeMove(m.from, m.to); err != nil {
			fmt.Printf("  Error: %v\n", err)
		}
	}

	g.DisplayBoard()
	fmt.Println()

	fmt.Println("  Move: e4×d5 (White pawn captures Black pawn)")
	if err := g.MakeMove(piece.Position{Row: 3, Col: 4}, piece.Position{Row: 4, Col: 3}); err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	g.DisplayBoard()
	fmt.Println()

	fmt.Println("  Move History:")
	for _, m := range g.GetMoveHistory() {
		fmt.Printf("    %s\n", m)
	}

	g.ViewStatus()
	fmt.Println()
}

func scenario3() {
	fmt.Println("--- Scenario 3: Invalid Moves ---")
	fmt.Println()

	game.ResetInstance()
	g := game.GetInstance()
	g.Initialize("Alice", "Bob")

	fmt.Println("  Trying to move opponent's piece (Black pawn on White's turn):")
	err := g.MakeMove(piece.Position{Row: 6, Col: 4}, piece.Position{Row: 5, Col: 4})
	fmt.Printf("  Error: %v\n", err)
	fmt.Println()

	fmt.Println("  Trying invalid rook move (diagonal):")
	err = g.MakeMove(piece.Position{Row: 0, Col: 0}, piece.Position{Row: 2, Col: 2})
	fmt.Printf("  Error: %v\n", err)
	fmt.Println()

	fmt.Println("  Trying to move pawn 3 squares forward:")
	err = g.MakeMove(piece.Position{Row: 1, Col: 4}, piece.Position{Row: 4, Col: 4})
	fmt.Printf("  Error: %v\n", err)
	fmt.Println()

	fmt.Println("  Trying to move from an empty square:")
	err = g.MakeMove(piece.Position{Row: 3, Col: 3}, piece.Position{Row: 4, Col: 3})
	fmt.Printf("  Error: %v\n", err)
	fmt.Println()
}

func scenario4() {
	fmt.Println("--- Scenario 4: Scholar's Mate (4-move Checkmate) ---")
	fmt.Println()

	game.ResetInstance()
	g := game.GetInstance()
	g.Initialize("Alice", "Bob")

	moves := []struct {
		from, to piece.Position
		desc     string
	}{
		{piece.Position{Row: 1, Col: 4}, piece.Position{Row: 3, Col: 4}, "1. e4"},
		{piece.Position{Row: 6, Col: 4}, piece.Position{Row: 4, Col: 4}, "1... e5"},
		{piece.Position{Row: 0, Col: 5}, piece.Position{Row: 3, Col: 2}, "2. Bc4 (f1→c4)"},
		{piece.Position{Row: 7, Col: 1}, piece.Position{Row: 5, Col: 2}, "2... Nc6 (b8→c6)"},
		{piece.Position{Row: 0, Col: 3}, piece.Position{Row: 4, Col: 7}, "3. Qh5 (d1→h5)"},
		{piece.Position{Row: 7, Col: 6}, piece.Position{Row: 5, Col: 5}, "3... Nf6 (g8→f6)"},
		{piece.Position{Row: 4, Col: 7}, piece.Position{Row: 6, Col: 5}, "4. Qxf7# (h5→f7)"},
	}

	for _, m := range moves {
		fmt.Printf("  Move: %s\n", m.desc)
		if err := g.MakeMove(m.from, m.to); err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}
	}

	fmt.Println()
	g.DisplayBoard()
	fmt.Println()

	fmt.Println("  Move History:")
	for _, m := range g.GetMoveHistory() {
		fmt.Printf("    %s\n", m)
	}

	fmt.Println()
	g.ViewStatus()

	status := g.GetStatus()
	if status == game.Checkmate {
		fmt.Println()
		fmt.Println("  ♚ CHECKMATE! Alice (White) wins! ♚")
	}
	fmt.Println()
}
