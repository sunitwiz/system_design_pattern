package pricing

import "book_my_show/theatre"

type PricingStrategy interface {
	CalculatePrice(seatType theatre.SeatType) float64
	GetMultiplier() float64
}
