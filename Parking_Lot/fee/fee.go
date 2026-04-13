package fee

import (
	"fmt"
	"math"
	"parking_lot/vehicle"
	"time"
)

// --- Concrete Strategies ---

type motorcycleFee struct{}

func (m *motorcycleFee) CalculateFee(duration time.Duration) float64 {
	hours := math.Ceil(duration.Hours())
	if hours < 1 {
		hours = 1
	}
	return hours * 1.0 // $1 per hour
}

func (m *motorcycleFee) GetRatePerHour() float64 { return 1.0 }

type carFee struct{}

func (c *carFee) CalculateFee(duration time.Duration) float64 {
	hours := math.Ceil(duration.Hours())
	if hours < 1 {
		hours = 1
	}
	return hours * 2.0 // $2 per hour
}

func (c *carFee) GetRatePerHour() float64 { return 2.0 }

type busFee struct{}

func (b *busFee) CalculateFee(duration time.Duration) float64 {
	hours := math.Ceil(duration.Hours())
	if hours < 1 {
		hours = 1
	}
	return hours * 5.0 // $5 per hour
}

func (b *busFee) GetRatePerHour() float64 { return 5.0 }

// GetFeeStrategy returns the appropriate fee strategy for a vehicle type.
func GetFeeStrategy(vType vehicle.VehicleType) (FeeStrategy, error) {
	switch vType {
	case vehicle.Motorcycle:
		return &motorcycleFee{}, nil
	case vehicle.Car:
		return &carFee{}, nil
	case vehicle.Bus:
		return &busFee{}, nil
	default:
		return nil, fmt.Errorf("no fee strategy for vehicle type: %s", vType)
	}
}
