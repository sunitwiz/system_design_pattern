package main

import (
	"fmt"
	"splitwise/service"
	"splitwise/split"
)

func main() {
	fmt.Println("=== Splitwise System Demo ===")
	fmt.Println()

	svc := service.GetInstance()

	fmt.Println("--- Adding Users ---")
	aliceID := svc.AddUser("Alice", "alice@example.com")
	bobID := svc.AddUser("Bob", "bob@example.com")
	charlieID := svc.AddUser("Charlie", "charlie@example.com")
	dianaID := svc.AddUser("Diana", "diana@example.com")
	fmt.Printf("  Created: Alice(%s), Bob(%s), Charlie(%s), Diana(%s)\n", aliceID, bobID, charlieID, dianaID)

	fmt.Println("\n--- Creating Group ---")
	groupID, err := svc.CreateGroup("Trip to Goa", []string{aliceID, bobID, charlieID, dianaID})
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
		return
	}
	fmt.Printf("  Created group: Trip to Goa (%s)\n", groupID)

	fmt.Println("\n--- Expense 1: Equal Split ---")
	fmt.Println("  Alice pays ₹2000 for dinner (split equally among all 4)")
	err = svc.AddExpenseToGroup(groupID, aliceID, 2000, split.EqualSplit,
		[]string{aliceID, bobID, charlieID, dianaID}, nil)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Println("  Expense added successfully")
	}

	fmt.Println("\n--- Expense 2: Exact Split ---")
	fmt.Println("  Bob pays ₹3000 for hotel (Alice:500, Bob:1000, Charlie:800, Diana:700)")
	err = svc.AddExpenseToGroup(groupID, bobID, 3000, split.ExactSplit,
		[]string{aliceID, bobID, charlieID, dianaID},
		map[string]float64{aliceID: 500, bobID: 1000, charlieID: 800, dianaID: 700})
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Println("  Expense added successfully")
	}

	fmt.Println("\n--- Expense 3: Percent Split ---")
	fmt.Println("  Charlie pays ₹1000 for taxi (Alice:40%, Bob:30%, Charlie:20%, Diana:10%)")
	err = svc.AddExpenseToGroup(groupID, charlieID, 1000, split.PercentSplit,
		[]string{aliceID, bobID, charlieID, dianaID},
		map[string]float64{aliceID: 40, bobID: 30, charlieID: 20, dianaID: 10})
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Println("  Expense added successfully")
	}

	fmt.Println("\n--- Balances After All Expenses ---")
	printUserBalances(svc, aliceID, "Alice")
	printUserBalances(svc, bobID, "Bob")
	printUserBalances(svc, charlieID, "Charlie")
	printUserBalances(svc, dianaID, "Diana")

	fmt.Println("\n--- Settle Up: Bob pays Alice ₹500 ---")
	err = svc.SettleUp(bobID, aliceID, 500)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Println("  Settlement recorded successfully")
	}

	fmt.Println("\n--- Balances After Settlement ---")
	printUserBalances(svc, aliceID, "Alice")
	printUserBalances(svc, bobID, "Bob")
	printUserBalances(svc, charlieID, "Charlie")
	printUserBalances(svc, dianaID, "Diana")

	fmt.Println("\n--- Group Expenses ---")
	expenses := svc.GetGroupExpenses(groupID)
	for _, exp := range expenses {
		fmt.Printf("  [%s] Payer: %s | Amount: ₹%.2f | Type: %s\n",
			exp.ID, exp.PayerID, exp.Amount, exp.SplitType)
		fmt.Printf("         Splits: ")
		first := true
		for uid, amt := range exp.SplitDetails {
			if !first {
				fmt.Printf(", ")
			}
			fmt.Printf("%s=₹%.2f", uid, amt)
			first = false
		}
		fmt.Println()
	}

	fmt.Println("\n--- Full Status ---")
	svc.ViewStatus()

	fmt.Println("\n--- Edge Case: Invalid Percent Split (sum ≠ 100) ---")
	err = svc.AddExpenseToGroup(groupID, aliceID, 500, split.PercentSplit,
		[]string{aliceID, bobID},
		map[string]float64{aliceID: 60, bobID: 50})
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n--- Edge Case: Non-existent Group ---")
	err = svc.AddExpenseToGroup("G999", aliceID, 100, split.EqualSplit,
		[]string{aliceID, bobID}, nil)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}
}

func printUserBalances(svc *service.SplitwiseService, userID, name string) {
	balances := svc.GetBalances(userID)
	if len(balances) == 0 {
		fmt.Printf("  %s: all settled up\n", name)
		return
	}
	fmt.Printf("  %s:\n", name)
	for otherID, amt := range balances {
		if amt > 0 {
			fmt.Printf("    owes %s: ₹%.2f\n", otherID, amt)
		} else if amt < 0 {
			fmt.Printf("    gets back from %s: ₹%.2f\n", otherID, -amt)
		}
	}
}
