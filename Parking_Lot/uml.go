package main

import "fmt"

func printClassDiagram() {
	fmt.Println(`classDiagram

    class VehicleType {
        <<enumeration>>
        Motorcycle
        Car
        Bus
        func (v VehicleType) String() string
    }

    class Vehicle {
        <<interface>>
        GetType() VehicleType
        GetLicensePlate() string
        String() string
    }

    class motorcycleVehicle {
        licensePlate string
        func (m *motorcycleVehicle) GetType() VehicleType
        func (m *motorcycleVehicle) GetLicensePlate() string
        func (m *motorcycleVehicle) String() string
    }

    class carVehicle {
        licensePlate string
        func (c *carVehicle) GetType() VehicleType
        func (c *carVehicle) GetLicensePlate() string
        func (c *carVehicle) String() string
    }

    class busVehicle {
        licensePlate string
        func (b *busVehicle) GetType() VehicleType
        func (b *busVehicle) GetLicensePlate() string
        func (b *busVehicle) String() string
    }

    class SlotType {
        <<enumeration>>
        MotorcycleSlot
        CarSlot
        BusSlot
        func (s SlotType) String() string
    }

    class ParkingSlot {
        ID            int
        Type          SlotType
        IsOccupied    bool
        ParkedVehicle vehicle.Vehicle
        func NewParkingSlot(id int, slotType SlotType) *ParkingSlot
        func (ps *ParkingSlot) CanFit(vType vehicle.VehicleType) bool
        func (ps *ParkingSlot) Park(v vehicle.Vehicle) error
        func (ps *ParkingSlot) Unpark() (vehicle.Vehicle, error)
    }

    class SlotStatus {
        SlotType string
        Free     int
        Occupied int
    }

    class ParkingLevel {
        LevelNumber int
        Slots       []*slot.ParkingSlot
        nextSlotID  int
        func NewParkingLevel(levelNumber int, motorcycleSlots, carSlots, busSlots int) *ParkingLevel
        func (pl *ParkingLevel) FindAvailableSlot(vType vehicle.VehicleType) *slot.ParkingSlot
        func (pl *ParkingLevel) GetStatus() []SlotStatus
        func (pl *ParkingLevel) AddSlot(slotType slot.SlotType) *slot.ParkingSlot
        func (pl *ParkingLevel) RemoveSlot(slotID int) error
    }

    class FeeStrategy {
        <<interface>>
        CalculateFee(duration time.Duration) float64
        GetRatePerHour() float64
    }

    class motorcycleFee {
        func (m *motorcycleFee) CalculateFee(duration time.Duration) float64
        func (m *motorcycleFee) GetRatePerHour() float64
    }

    class carFee {
        func (c *carFee) CalculateFee(duration time.Duration) float64
        func (c *carFee) GetRatePerHour() float64
    }

    class busFee {
        func (b *busFee) CalculateFee(duration time.Duration) float64
        func (b *busFee) GetRatePerHour() float64
    }

    class ParkingTicket {
        TicketID  string
        Vehicle   vehicle.Vehicle
        Slot      *slot.ParkingSlot
        LevelNum  int
        EntryTime time.Time
        ExitTime  time.Time
        Fee       float64
        IsActive  bool
        func NewParkingTicket(id string, v vehicle.Vehicle, s *slot.ParkingSlot, levelNum int) *ParkingTicket
        func (t *ParkingTicket) CalculateFee() (float64, error)
        func (t *ParkingTicket) CalculateFeeWithTime(exitTime time.Time) (float64, error)
        func (t *ParkingTicket) String() string
    }

    class ParkingLotOperations {
        <<interface>>
        ParkVehicle(v vehicle.Vehicle) (*ticket.ParkingTicket, error)
        UnparkVehicle(ticketID string) (*ticket.ParkingTicket, error)
        ViewStatus()
        AddLevel(motorcycleSlots, carSlots, busSlots int)
        RemoveLevel(levelNumber int) error
        AddSlot(levelNumber int, slotType slot.SlotType) error
        RemoveSlot(levelNumber int, slotID int) error
    }

    class ParkingLot {
        mu            sync.Mutex
        Levels        []*level.ParkingLevel
        ActiveTickets map[string]*ticket.ParkingTicket
        ticketCounter int
        func GetInstance() *ParkingLot
        func (pl *ParkingLot) ParkVehicle(v vehicle.Vehicle) (*ticket.ParkingTicket, error)
        func (pl *ParkingLot) UnparkVehicle(ticketID string) (*ticket.ParkingTicket, error)
        func (pl *ParkingLot) ViewStatus()
        func (pl *ParkingLot) AddLevel(motorcycleSlots, carSlots, busSlots int)
        func (pl *ParkingLot) RemoveLevel(levelNumber int) error
        func (pl *ParkingLot) AddSlot(levelNumber int, slotType slot.SlotType) error
        func (pl *ParkingLot) RemoveSlot(levelNumber int, slotID int) error
    }

    motorcycleVehicle ..|> Vehicle : implements
    carVehicle ..|> Vehicle : implements
    busVehicle ..|> Vehicle : implements
    motorcycleFee ..|> FeeStrategy : implements
    carFee ..|> FeeStrategy : implements
    busFee ..|> FeeStrategy : implements
    ParkingSlot --> SlotType
    ParkingSlot --> Vehicle : holds when occupied
    ParkingLevel *-- ParkingSlot : contains
    ParkingLevel --> SlotStatus
    ParkingTicket --> Vehicle : references
    ParkingTicket --> ParkingSlot : references
    ParkingTicket ..> FeeStrategy : uses
    ParkingLot ..|> ParkingLotOperations : implements
    ParkingLot *-- ParkingLevel : contains
    ParkingLot o-- ParkingTicket : tracks`)
}
