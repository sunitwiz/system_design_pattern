package main

import (
	"book_my_show/bookmyshow"
	"book_my_show/movie"
	"book_my_show/pricing"
	"book_my_show/theatre"
	"fmt"
	"time"
)

func main() {
	fmt.Println("=== BookMyShow System Demo ===")
	fmt.Println()

	bms := bookmyshow.GetInstance()

	fmt.Println("--- Admin: Adding Theatre ---")
	pvr := theatre.NewTheatre("TH-001", "PVR Cinemas", "Mumbai")
	bms.AddTheatre(pvr)

	fmt.Println("\n--- Admin: Adding Screens ---")
	screen1 := theatre.NewScreen(1, 1, buildScreenSeats(6, 3, 1))
	screen2 := theatre.NewScreen(2, 2, buildScreenSeats(4, 3, 1))
	bms.AddScreen("TH-001", screen1)
	bms.AddScreen("TH-001", screen2)

	fmt.Println("\n--- Admin: Adding Movies ---")
	avengers := movie.NewMovie("MOV-001", "Avengers: Endgame", 181, movie.Action, 8.4)
	hangover := movie.NewMovie("MOV-002", "The Hangover", 100, movie.Comedy, 7.7)
	bms.AddMovie(avengers)
	bms.AddMovie(hangover)

	fmt.Println("\n--- Admin: Scheduling Shows ---")
	saturdayEvening := time.Date(2025, 7, 12, 18, 30, 0, 0, time.Local)
	mondayMatinee := time.Date(2025, 7, 14, 11, 0, 0, 0, time.Local)
	sundayMatinee := time.Date(2025, 7, 13, 14, 0, 0, 0, time.Local)

	bms.AddShow("SH-001", "MOV-001", "TH-001", 1, saturdayEvening)
	bms.AddShow("SH-002", "MOV-002", "TH-001", 2, mondayMatinee)
	bms.AddShow("SH-003", "MOV-001", "TH-001", 2, sundayMatinee)

	fmt.Println("\n--- User: Searching for \"Avengers\" ---")
	results := bms.SearchMovies("Avengers")
	for _, m := range results {
		fmt.Printf("  Found: %s\n", m)
	}

	fmt.Println("\n--- User: Shows for Avengers in Mumbai ---")
	shows := bms.GetShows("MOV-001", "Mumbai")
	for _, s := range shows {
		fmt.Printf("  %s\n", s)
	}

	fmt.Println("\n--- User: Booking Weekend Show (SH-001) ---")
	b1, err := bms.BookTickets("John", "SH-001", []int{1, 7, 10})
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  %s\n", b1)
	}

	fmt.Println("\n  Pricing Breakdown (Weekend 1.5x):")
	strategy := pricing.GetPricingStrategy(true)
	var total float64
	for _, seat := range b1.Seats {
		base := pricing.GetBasePrice(seat.Type)
		price := strategy.CalculatePrice(seat.Type)
		total += price
		fmt.Printf("    %-14s $%.2f × %.1f = $%.2f\n", seat, base, strategy.GetMultiplier(), price)
	}
	fmt.Printf("    %-14s              $%.2f\n", "Total:", total)

	fmt.Println("\n--- Edge Case: Double Booking Seat B1 ---")
	_, err = bms.BookTickets("Bob", "SH-001", []int{7})
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n--- User: Booking Weekday Show (SH-002) ---")
	b2, err := bms.BookTickets("Jane", "SH-002", []int{1, 5})
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  %s\n", b2)
	}

	fmt.Println("\n  Pricing Breakdown (Weekday 1.0x):")
	strategy = pricing.GetPricingStrategy(false)
	total = 0
	for _, seat := range b2.Seats {
		base := pricing.GetBasePrice(seat.Type)
		price := strategy.CalculatePrice(seat.Type)
		total += price
		fmt.Printf("    %-14s $%.2f × %.1f = $%.2f\n", seat, base, strategy.GetMultiplier(), price)
	}
	fmt.Printf("    %-14s              $%.2f\n", "Total:", total)

	fmt.Println("\n--- Seat Availability After Bookings ---")
	for _, s := range bms.GetShows("MOV-001", "Mumbai") {
		fmt.Printf("  %s\n", s)
	}
	for _, s := range bms.GetShows("MOV-002", "Mumbai") {
		fmt.Printf("  %s\n", s)
	}

	fmt.Println("\n--- User: Cancelling Booking " + b1.ID + " ---")
	err = bms.CancelBooking(b1.ID)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Cancelled: %s\n", b1)
	}

	fmt.Println("\n--- Seat Availability After Cancellation ---")
	for _, s := range bms.GetShows("MOV-001", "Mumbai") {
		fmt.Printf("  %s\n", s)
	}

	fmt.Println("\n--- Edge Case: Remove Show With Active Bookings ---")
	err = bms.RemoveShow("SH-002")
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n--- Admin: Removing Unbooked Show ---")
	err = bms.RemoveShow("SH-003")
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	fmt.Println("\n--- Final: All Shows in Mumbai ---")
	for _, s := range bms.GetShows("MOV-001", "Mumbai") {
		fmt.Printf("  %s\n", s)
	}
	for _, s := range bms.GetShows("MOV-002", "Mumbai") {
		fmt.Printf("  %s\n", s)
	}
}

func buildScreenSeats(regular, premium, vip int) []*theatre.Seat {
	var seats []*theatre.Seat
	id := 1
	for i := 1; i <= regular; i++ {
		seats = append(seats, theatre.NewSeat(id, "A", i, theatre.Regular))
		id++
	}
	for i := 1; i <= premium; i++ {
		seats = append(seats, theatre.NewSeat(id, "B", i, theatre.Premium))
		id++
	}
	for i := 1; i <= vip; i++ {
		seats = append(seats, theatre.NewSeat(id, "C", i, theatre.VIP))
		id++
	}
	return seats
}
