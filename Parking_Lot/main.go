package main

import (
	"fmt"
	"parking_lot/parkinglot"
	"parking_lot/slot"
	"parking_lot/vehicle"
	"time"
)

func main() {
	fmt.Println("=== Parking Lot System Demo ===\n")

	// 1. Initialize the Parking Lot (Singleton)
	lot := parkinglot.GetInstance()

	// 2. Admin: Add levels
	fmt.Println("--- Admin: Adding Levels ---")
	lot.AddLevel(5, 10, 2) // Level 1: 5 motorcycle, 10 car, 2 bus slots
	lot.AddLevel(3, 8, 1)  // Level 2: 3 motorcycle, 8 car, 1 bus slot
	fmt.Println()

	// 3. View initial status
	lot.ViewStatus()

	// 4. Create vehicles using the Factory
	fmt.Println("\n--- Parking Vehicles ---")
	motorcycle, _ := vehicle.NewVehicle(vehicle.Motorcycle, "MC-1001")
	car1, _ := vehicle.NewVehicle(vehicle.Car, "CAR-2001")
	car2, _ := vehicle.NewVehicle(vehicle.Car, "CAR-2002")
	bus, _ := vehicle.NewVehicle(vehicle.Bus, "BUS-3001")

	// 5. Park vehicles
	t1, err := lot.ParkVehicle(motorcycle)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Parked: %s\n", t1)
	}

	t2, err := lot.ParkVehicle(car1)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Parked: %s\n", t2)
	}

	t3, err := lot.ParkVehicle(car2)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Parked: %s\n", t3)
	}

	t4, err := lot.ParkVehicle(bus)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Parked: %s\n", t4)
	}

	// 6. View status after parking
	fmt.Println("\n--- Status After Parking ---")
	lot.ViewStatus()

	// 7. Simulate time passing and unpark
	fmt.Println("\n--- Unparking Vehicles (simulating 3-hour stay) ---")

	// Simulate a 3-hour stay by modifying the entry time
	t1.EntryTime = time.Now().Add(-3 * time.Hour)
	t2.EntryTime = time.Now().Add(-3 * time.Hour)
	t4.EntryTime = time.Now().Add(-3 * time.Hour)

	closedT1, err := lot.UnparkVehicle(t1.TicketID)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Unparked: %s\n", closedT1)
	}

	closedT2, err := lot.UnparkVehicle(t2.TicketID)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Unparked: %s\n", closedT2)
	}

	closedT4, err := lot.UnparkVehicle(t4.TicketID)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Unparked: %s\n", closedT4)
	}

	// 8. Fee summary
	fmt.Println("\n--- Fee Summary ---")
	fmt.Printf("  Motorcycle (3 hrs @ $1/hr): $%.2f\n", closedT1.Fee)
	fmt.Printf("  Car        (3 hrs @ $2/hr): $%.2f\n", closedT2.Fee)
	fmt.Printf("  Bus        (3 hrs @ $5/hr): $%.2f\n", closedT4.Fee)

	// 9. Try parking when slot is not available
	fmt.Println("\n--- Edge Case: Parking Full Scenario ---")
	// Fill all bus slots
	for i := 0; i < 5; i++ {
		b, _ := vehicle.NewVehicle(vehicle.Bus, fmt.Sprintf("BUS-400%d", i))
		_, err := lot.ParkVehicle(b)
		if err != nil {
			fmt.Printf("  Bus %s: %v\n", b.GetLicensePlate(), err)
		} else {
			fmt.Printf("  Bus %s: parked successfully\n", b.GetLicensePlate())
		}
	}

	// 10. Admin: Add a slot dynamically
	fmt.Println("\n--- Admin: Adding Extra Bus Slot to Level 1 ---")
	lot.AddSlot(1, slot.BusSlot)

	// 11. Final status
	fmt.Println("\n--- Final Status ---")
	lot.ViewStatus()
}
