package parkinglot

import (
	"parking_lot/slot"
	"parking_lot/ticket"
	"parking_lot/vehicle"
)

type ParkingLotOperations interface {
	ParkVehicle(v vehicle.Vehicle) (*ticket.ParkingTicket, error)
	UnparkVehicle(ticketID string) (*ticket.ParkingTicket, error)
	ViewStatus()
	AddLevel(motorcycleSlots, carSlots, busSlots int)
	RemoveLevel(levelNumber int) error
	AddSlot(levelNumber int, slotType slot.SlotType) error
	RemoveSlot(levelNumber int, slotID int) error
}
