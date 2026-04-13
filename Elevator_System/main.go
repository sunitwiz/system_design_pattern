package main

import (
	"elevator_system/controller"
	"elevator_system/scheduler"
	"fmt"
)

func main() {
	fmt.Println("=== Elevator System Demo ===\n")

	ctrl := controller.GetInstance()

	fmt.Println("--- Admin: Adding Elevators ---")
	ctrl.AddElevator(1)
	ctrl.AddElevator(2)
	ctrl.AddElevator(3)
	fmt.Println()

	fmt.Println("--- Setting Nearest Scheduler ---")
	nearestScheduler, _ := scheduler.NewScheduler(scheduler.Nearest)
	ctrl.SetScheduler(nearestScheduler)
	fmt.Println()

	ctrl.ViewStatus()

	fmt.Println("\n--- Requesting Elevators (Nearest Scheduler) ---")

	e1, err := ctrl.RequestElevator(3, 7)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Person at floor 3 → 7: assigned to Elevator %d\n", e1.ID)
	}

	e2, err := ctrl.RequestElevator(1, 5)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Person at floor 1 → 5: assigned to Elevator %d\n", e2.ID)
	}

	e3, err := ctrl.RequestElevator(8, 2)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Person at floor 8 → 2: assigned to Elevator %d\n", e3.ID)
	}

	fmt.Println("\n--- Simulation: Stepping Elevators ---")
	for step := 1; step <= 10; step++ {
		fmt.Printf("\n  [Step %d]\n", step)
		ctrl.StepAll()
		ctrl.ViewStatus()
	}

	fmt.Println("\n--- Switching to Round Robin Scheduler ---")
	rrScheduler, _ := scheduler.NewScheduler(scheduler.RoundRobin)
	ctrl.SetScheduler(rrScheduler)
	fmt.Println()

	controller.ResetInstance()
	ctrl = controller.GetInstance()
	ctrl.AddElevator(1)
	ctrl.AddElevator(2)
	ctrl.AddElevator(3)
	ctrl.SetScheduler(rrScheduler)
	fmt.Println()

	fmt.Println("--- Requesting Elevators (Round Robin Scheduler) ---")

	e1, err = ctrl.RequestElevator(3, 7)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Person at floor 3 → 7: assigned to Elevator %d\n", e1.ID)
	}

	e2, err = ctrl.RequestElevator(1, 5)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Person at floor 1 → 5: assigned to Elevator %d\n", e2.ID)
	}

	e3, err = ctrl.RequestElevator(8, 2)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Person at floor 8 → 2: assigned to Elevator %d\n", e3.ID)
	}

	fmt.Println()
	ctrl.ViewStatus()

	fmt.Println("\n--- Edge Case: No Elevators Available ---")
	controller.ResetInstance()
	ctrl = controller.GetInstance()
	_, err = ctrl.RequestElevator(5, 10)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n--- Edge Case: Remove Busy Elevator ---")
	ctrl.AddElevator(1)
	ctrl.RequestElevator(3, 7)
	err = ctrl.RemoveElevator(1)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n--- Final Status ---")
	ctrl.ViewStatus()
}
