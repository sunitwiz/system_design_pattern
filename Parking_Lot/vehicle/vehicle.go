package vehicle

import "fmt"

// VehicleType represents the type of vehicle.
type VehicleType int

const (
	Motorcycle VehicleType = iota
	Car
	Bus
)

func (v VehicleType) String() string {
	switch v {
	case Motorcycle:
		return "Motorcycle"
	case Car:
		return "Car"
	case Bus:
		return "Bus"
	default:
		return "Unknown"
	}
}

// Vehicle is the interface that all vehicle types must implement.
type Vehicle interface {
	GetType() VehicleType
	GetLicensePlate() string
	String() string
}

// --- Concrete Vehicle Types ---

type motorcycleVehicle struct {
	licensePlate string
}

func (m *motorcycleVehicle) GetType() VehicleType    { return Motorcycle }
func (m *motorcycleVehicle) GetLicensePlate() string  { return m.licensePlate }
func (m *motorcycleVehicle) String() string {
	return fmt.Sprintf("Motorcycle [%s]", m.licensePlate)
}

type carVehicle struct {
	licensePlate string
}

func (c *carVehicle) GetType() VehicleType    { return Car }
func (c *carVehicle) GetLicensePlate() string  { return c.licensePlate }
func (c *carVehicle) String() string {
	return fmt.Sprintf("Car [%s]", c.licensePlate)
}

type busVehicle struct {
	licensePlate string
}

func (b *busVehicle) GetType() VehicleType    { return Bus }
func (b *busVehicle) GetLicensePlate() string  { return b.licensePlate }
func (b *busVehicle) String() string {
	return fmt.Sprintf("Bus [%s]", b.licensePlate)
}

// NewVehicle is a factory function that creates a Vehicle based on its type.
func NewVehicle(vType VehicleType, licensePlate string) (Vehicle, error) {
	switch vType {
	case Motorcycle:
		return &motorcycleVehicle{licensePlate: licensePlate}, nil
	case Car:
		return &carVehicle{licensePlate: licensePlate}, nil
	case Bus:
		return &busVehicle{licensePlate: licensePlate}, nil
	default:
		return nil, fmt.Errorf("unknown vehicle type: %d", vType)
	}
}
