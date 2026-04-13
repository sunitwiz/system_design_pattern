package slot

import (
	"fmt"
	"parking_lot/vehicle"
)

// SlotType represents the type of parking slot.
type SlotType int

const (
	MotorcycleSlot SlotType = iota
	CarSlot
	BusSlot
)

func (s SlotType) String() string {
	switch s {
	case MotorcycleSlot:
		return "Motorcycle Slot"
	case CarSlot:
		return "Car Slot"
	case BusSlot:
		return "Bus Slot"
	default:
		return "Unknown Slot"
	}
}

// ParkingSlot represents a single parking slot.
type ParkingSlot struct {
	ID            int
	Type          SlotType
	IsOccupied    bool
	ParkedVehicle vehicle.Vehicle
}

// NewParkingSlot creates a new parking slot.
func NewParkingSlot(id int, slotType SlotType) *ParkingSlot {
	return &ParkingSlot{
		ID:   id,
		Type: slotType,
	}
}

// CanFit checks if a vehicle type can fit into this slot.
// Rules:
//   - MotorcycleSlot: only motorcycles
//   - CarSlot: cars and motorcycles
//   - BusSlot: buses, cars, and motorcycles
func (ps *ParkingSlot) CanFit(vType vehicle.VehicleType) bool {
	switch ps.Type {
	case MotorcycleSlot:
		return vType == vehicle.Motorcycle
	case CarSlot:
		return vType == vehicle.Car || vType == vehicle.Motorcycle
	case BusSlot:
		return true // accepts all vehicle types
	default:
		return false
	}
}

// Park parks a vehicle in this slot.
func (ps *ParkingSlot) Park(v vehicle.Vehicle) error {
	if ps.IsOccupied {
		return fmt.Errorf("slot %d is already occupied", ps.ID)
	}
	if !ps.CanFit(v.GetType()) {
		return fmt.Errorf("vehicle type %s cannot fit in %s", v.GetType(), ps.Type)
	}
	ps.IsOccupied = true
	ps.ParkedVehicle = v
	return nil
}

// Unpark removes the vehicle from this slot.
func (ps *ParkingSlot) Unpark() (vehicle.Vehicle, error) {
	if !ps.IsOccupied {
		return nil, fmt.Errorf("slot %d is already empty", ps.ID)
	}
	v := ps.ParkedVehicle
	ps.IsOccupied = false
	ps.ParkedVehicle = nil
	return v, nil
}
