package fee

import "time"

type FeeStrategy interface {
	CalculateFee(duration time.Duration) float64
	GetRatePerHour() float64
}
