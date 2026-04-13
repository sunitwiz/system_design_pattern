package parkinglot

import (
	"fmt"
	"parking_lot/level"
	"parking_lot/slot"
	"parking_lot/ticket"
	"parking_lot/vehicle"
	"sync"
)

// ParkingOperations defines the interface for parking/unparking vehicles.
type ParkingOperations interface {
	ParkVehicle(v vehicle.Vehicle) (*ticket.ParkingTicket, error)
	UnparkVehicle(ticketID string) (*ticket.ParkingTicket, error)
}

// AdminOperations defines the interface for admin management features.
type AdminOperations interface {
	ViewStatus()
	AddLevel(motorcycleSlots, carSlots, busSlots int)
	RemoveLevel(levelNumber int) error
	AddSlot(levelNumber int, slotType slot.SlotType) error
	RemoveSlot(levelNumber int, slotID int) error
}

// ParkingLot is the core orchestrator implementing both interfaces.
type ParkingLot struct {
	mu            sync.Mutex
	Levels        []*level.ParkingLevel
	ActiveTickets map[string]*ticket.ParkingTicket
	ticketCounter int
}

// Compile-time interface checks.
var _ ParkingOperations = (*ParkingLot)(nil)
var _ AdminOperations = (*ParkingLot)(nil)

// singleton instance
var (
	instance *ParkingLot
	once     sync.Once
)

// GetInstance returns the singleton ParkingLot instance.
func GetInstance() *ParkingLot {
	once.Do(func() {
		instance = &ParkingLot{
			ActiveTickets: make(map[string]*ticket.ParkingTicket),
		}
	})
	return instance
}

// ResetInstance resets the singleton (useful for testing).
func ResetInstance() {
	once = sync.Once{}
	instance = nil
}

// --- ParkingOperations Implementation ---

// ParkVehicle finds an available slot and parks the vehicle.
func (pl *ParkingLot) ParkVehicle(v vehicle.Vehicle) (*ticket.ParkingTicket, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	for _, lvl := range pl.Levels {
		s := lvl.FindAvailableSlot(v.GetType())
		if s != nil {
			if err := s.Park(v); err != nil {
				return nil, err
			}
			pl.ticketCounter++
			ticketID := fmt.Sprintf("T-%04d", pl.ticketCounter)
			t := ticket.NewParkingTicket(ticketID, v, s, lvl.LevelNumber)
			pl.ActiveTickets[ticketID] = t
			return t, nil
		}
	}

	return nil, fmt.Errorf("no available slot for %s", v)
}

// UnparkVehicle removes a vehicle and calculates the fee.
func (pl *ParkingLot) UnparkVehicle(ticketID string) (*ticket.ParkingTicket, error) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	t, exists := pl.ActiveTickets[ticketID]
	if !exists {
		return nil, fmt.Errorf("ticket %s not found", ticketID)
	}

	if _, err := t.Slot.Unpark(); err != nil {
		return nil, err
	}

	if _, err := t.CalculateFee(); err != nil {
		return nil, err
	}

	delete(pl.ActiveTickets, ticketID)
	return t, nil
}

// --- AdminOperations Implementation ---

// ViewStatus prints the current status of the entire parking lot.
func (pl *ParkingLot) ViewStatus() {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	fmt.Println("========================================")
	fmt.Println("        PARKING LOT STATUS")
	fmt.Println("========================================")

	if len(pl.Levels) == 0 {
		fmt.Println("  No levels configured.")
		return
	}

	for _, lvl := range pl.Levels {
		fmt.Printf("\n  Level %d:\n", lvl.LevelNumber)
		fmt.Printf("  %-20s %6s %8s\n", "Slot Type", "Free", "Occupied")
		fmt.Printf("  %-20s %6s %8s\n", "--------------------", "------", "--------")

		statuses := lvl.GetStatus()
		for _, s := range statuses {
			if s.Free+s.Occupied > 0 {
				fmt.Printf("  %-20s %6d %8d\n", s.SlotType, s.Free, s.Occupied)
			}
		}
	}

	fmt.Printf("\n  Active Tickets: %d\n", len(pl.ActiveTickets))
	fmt.Println("========================================")
}

// AddLevel adds a new level to the parking lot.
func (pl *ParkingLot) AddLevel(motorcycleSlots, carSlots, busSlots int) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	levelNum := len(pl.Levels) + 1
	lvl := level.NewParkingLevel(levelNum, motorcycleSlots, carSlots, busSlots)
	pl.Levels = append(pl.Levels, lvl)
	fmt.Printf("  Added Level %d (%d motorcycle, %d car, %d bus slots)\n",
		levelNum, motorcycleSlots, carSlots, busSlots)
}

// RemoveLevel removes an empty level by its number.
func (pl *ParkingLot) RemoveLevel(levelNumber int) error {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	for i, lvl := range pl.Levels {
		if lvl.LevelNumber == levelNumber {
			// Check if any slot is occupied
			for _, s := range lvl.Slots {
				if s.IsOccupied {
					return fmt.Errorf("cannot remove level %d: has occupied slots", levelNumber)
				}
			}
			pl.Levels = append(pl.Levels[:i], pl.Levels[i+1:]...)
			fmt.Printf("  Removed Level %d\n", levelNumber)
			return nil
		}
	}
	return fmt.Errorf("level %d not found", levelNumber)
}

// AddSlot adds a slot to a specific level.
func (pl *ParkingLot) AddSlot(levelNumber int, slotType slot.SlotType) error {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	for _, lvl := range pl.Levels {
		if lvl.LevelNumber == levelNumber {
			s := lvl.AddSlot(slotType)
			fmt.Printf("  Added %s (ID: %d) to Level %d\n", s.Type, s.ID, levelNumber)
			return nil
		}
	}
	return fmt.Errorf("level %d not found", levelNumber)
}

// RemoveSlot removes a slot from a specific level.
func (pl *ParkingLot) RemoveSlot(levelNumber int, slotID int) error {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	for _, lvl := range pl.Levels {
		if lvl.LevelNumber == levelNumber {
			return lvl.RemoveSlot(slotID)
		}
	}
	return fmt.Errorf("level %d not found", levelNumber)
}
