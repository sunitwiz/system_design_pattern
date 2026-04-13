package level

import (
	"fmt"
	"parking_lot/slot"
	"parking_lot/vehicle"
)

// ParkingLevel represents a single level in the parking lot.
type ParkingLevel struct {
	LevelNumber int
	Slots       []*slot.ParkingSlot
	nextSlotID  int
}

// NewParkingLevel creates a new parking level with the given slot configuration.
func NewParkingLevel(levelNumber int, motorcycleSlots, carSlots, busSlots int) *ParkingLevel {
	pl := &ParkingLevel{
		LevelNumber: levelNumber,
		nextSlotID:  1,
	}

	for i := 0; i < motorcycleSlots; i++ {
		pl.Slots = append(pl.Slots, slot.NewParkingSlot(pl.nextSlotID, slot.MotorcycleSlot))
		pl.nextSlotID++
	}
	for i := 0; i < carSlots; i++ {
		pl.Slots = append(pl.Slots, slot.NewParkingSlot(pl.nextSlotID, slot.CarSlot))
		pl.nextSlotID++
	}
	for i := 0; i < busSlots; i++ {
		pl.Slots = append(pl.Slots, slot.NewParkingSlot(pl.nextSlotID, slot.BusSlot))
		pl.nextSlotID++
	}

	return pl
}

// FindAvailableSlot finds the first available slot that can fit the given vehicle type.
// It prioritizes the smallest suitable slot to minimize waste.
func (pl *ParkingLevel) FindAvailableSlot(vType vehicle.VehicleType) *slot.ParkingSlot {
	// First pass: find exact match slots
	for _, s := range pl.Slots {
		if !s.IsOccupied && isExactMatch(s.Type, vType) {
			return s
		}
	}
	// Second pass: find any compatible slot
	for _, s := range pl.Slots {
		if !s.IsOccupied && s.CanFit(vType) {
			return s
		}
	}
	return nil
}

// isExactMatch returns true if the slot type is the natural fit for the vehicle type.
func isExactMatch(sType slot.SlotType, vType vehicle.VehicleType) bool {
	switch vType {
	case vehicle.Motorcycle:
		return sType == slot.MotorcycleSlot
	case vehicle.Car:
		return sType == slot.CarSlot
	case vehicle.Bus:
		return sType == slot.BusSlot
	default:
		return false
	}
}

// SlotStatus holds the count of free and occupied slots for a given slot type.
type SlotStatus struct {
	SlotType string
	Free     int
	Occupied int
}

// GetStatus returns the status of all slot types on this level.
func (pl *ParkingLevel) GetStatus() []SlotStatus {
	counts := map[slot.SlotType]*SlotStatus{
		slot.MotorcycleSlot: {SlotType: slot.MotorcycleSlot.String()},
		slot.CarSlot:        {SlotType: slot.CarSlot.String()},
		slot.BusSlot:        {SlotType: slot.BusSlot.String()},
	}

	for _, s := range pl.Slots {
		if s.IsOccupied {
			counts[s.Type].Occupied++
		} else {
			counts[s.Type].Free++
		}
	}

	return []SlotStatus{
		*counts[slot.MotorcycleSlot],
		*counts[slot.CarSlot],
		*counts[slot.BusSlot],
	}
}

// AddSlot adds a new slot of the given type to this level.
func (pl *ParkingLevel) AddSlot(slotType slot.SlotType) *slot.ParkingSlot {
	s := slot.NewParkingSlot(pl.nextSlotID, slotType)
	pl.Slots = append(pl.Slots, s)
	pl.nextSlotID++
	return s
}

// RemoveSlot removes an unoccupied slot by its ID.
func (pl *ParkingLevel) RemoveSlot(slotID int) error {
	for i, s := range pl.Slots {
		if s.ID == slotID {
			if s.IsOccupied {
				return fmt.Errorf("cannot remove occupied slot %d", slotID)
			}
			pl.Slots = append(pl.Slots[:i], pl.Slots[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("slot %d not found on level %d", slotID, pl.LevelNumber)
}
