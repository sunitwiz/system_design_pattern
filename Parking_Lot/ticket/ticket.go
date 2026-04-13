package ticket

import (
	"fmt"
	"parking_lot/fee"
	"parking_lot/slot"
	"parking_lot/vehicle"
	"time"
)

// ParkingTicket represents a ticket issued when a vehicle is parked.
type ParkingTicket struct {
	TicketID  string
	Vehicle   vehicle.Vehicle
	Slot      *slot.ParkingSlot
	LevelNum  int
	EntryTime time.Time
	ExitTime  time.Time
	Fee       float64
	IsActive  bool
}

// NewParkingTicket creates a new parking ticket.
func NewParkingTicket(id string, v vehicle.Vehicle, s *slot.ParkingSlot, levelNum int) *ParkingTicket {
	return &ParkingTicket{
		TicketID:  id,
		Vehicle:   v,
		Slot:      s,
		LevelNum:  levelNum,
		EntryTime: time.Now(),
		IsActive:  true,
	}
}

// CalculateFee calculates the parking fee using the Strategy pattern.
func (t *ParkingTicket) CalculateFee() (float64, error) {
	strategy, err := fee.GetFeeStrategy(t.Vehicle.GetType())
	if err != nil {
		return 0, err
	}
	t.ExitTime = time.Now()
	duration := t.ExitTime.Sub(t.EntryTime)
	t.Fee = strategy.CalculateFee(duration)
	t.IsActive = false
	return t.Fee, nil
}

// CalculateFeeWithTime calculates the fee for a specific exit time (useful for testing).
func (t *ParkingTicket) CalculateFeeWithTime(exitTime time.Time) (float64, error) {
	strategy, err := fee.GetFeeStrategy(t.Vehicle.GetType())
	if err != nil {
		return 0, err
	}
	t.ExitTime = exitTime
	duration := t.ExitTime.Sub(t.EntryTime)
	t.Fee = strategy.CalculateFee(duration)
	t.IsActive = false
	return t.Fee, nil
}

func (t *ParkingTicket) String() string {
	status := "ACTIVE"
	if !t.IsActive {
		status = "CLOSED"
	}
	return fmt.Sprintf("Ticket[%s] %s | %s | Level %d, Slot %d (%s) | Entry: %s | Fee: $%.2f",
		t.TicketID,
		status,
		t.Vehicle,
		t.LevelNum,
		t.Slot.ID,
		t.Slot.Type,
		t.EntryTime.Format("15:04:05"),
		t.Fee,
	)
}
