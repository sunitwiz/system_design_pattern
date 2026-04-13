package pricing

import (
	"book_my_show/theatre"
	"time"
)

type PricingStrategy interface {
	CalculatePrice(seatType theatre.SeatType) float64
	GetMultiplier() float64
}

// --- Concrete Strategies ---

type regularPricing struct{}

func (r *regularPricing) CalculatePrice(seatType theatre.SeatType) float64 {
	return GetBasePrice(seatType) * 1.0
}

func (r *regularPricing) GetMultiplier() float64 { return 1.0 }

type weekendPricing struct{}

func (w *weekendPricing) CalculatePrice(seatType theatre.SeatType) float64 {
	return GetBasePrice(seatType) * 1.5
}

func (w *weekendPricing) GetMultiplier() float64 { return 1.5 }

type premiumPricing struct{}

func (p *premiumPricing) CalculatePrice(seatType theatre.SeatType) float64 {
	return GetBasePrice(seatType) * 2.0
}

func (p *premiumPricing) GetMultiplier() float64 { return 2.0 }

func GetBasePrice(seatType theatre.SeatType) float64 {
	switch seatType {
	case theatre.Regular:
		return 10.0
	case theatre.Premium:
		return 15.0
	case theatre.VIP:
		return 25.0
	default:
		return 10.0
	}
}

func IsWeekend(t time.Time) bool {
	day := t.Weekday()
	return day == time.Saturday || day == time.Sunday
}

func GetPricingStrategy(isWeekend bool) PricingStrategy {
	if isWeekend {
		return &weekendPricing{}
	}
	return &regularPricing{}
}

func GetPremiumPricingStrategy() PricingStrategy {
	return &premiumPricing{}
}
