package booking

import (
	"book_my_show/pricing"
	"book_my_show/show"
	"book_my_show/theatre"
	"fmt"
	"time"
)

type BookingStatus int

const (
	Pending BookingStatus = iota
	Confirmed
	Cancelled
)

func (s BookingStatus) String() string {
	switch s {
	case Pending:
		return "PENDING"
	case Confirmed:
		return "CONFIRMED"
	case Cancelled:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}

type Booking struct {
	ID          string
	UserName    string
	Show        *show.Show
	Seats       []*theatre.Seat
	TotalAmount float64
	Status      BookingStatus
	BookingTime time.Time
}

func NewBooking(id, userName string, s *show.Show, seats []*theatre.Seat) *Booking {
	b := &Booking{
		ID:          id,
		UserName:    userName,
		Show:        s,
		Seats:       seats,
		Status:      Confirmed,
		BookingTime: time.Now(),
	}
	b.CalculateTotal()
	return b
}

func (b *Booking) CalculateTotal() float64 {
	strategy := pricing.GetPricingStrategy(pricing.IsWeekend(b.Show.StartTime))
	var total float64
	for _, seat := range b.Seats {
		total += strategy.CalculatePrice(seat.Type)
	}
	b.TotalAmount = total
	return total
}

func (b *Booking) Cancel() error {
	if b.Status == Cancelled {
		return fmt.Errorf("booking %s is already cancelled", b.ID)
	}
	b.Status = Cancelled
	seatIDs := make([]int, len(b.Seats))
	for i, seat := range b.Seats {
		seatIDs[i] = seat.ID
	}
	return b.Show.CancelSeats(seatIDs)
}

func (b *Booking) String() string {
	seatNames := ""
	for i, seat := range b.Seats {
		if i > 0 {
			seatNames += ", "
		}
		seatNames += seat.String()
	}
	return fmt.Sprintf("Booking[%s] %s | %s | %s | Seats: [%s] | $%.2f",
		b.ID, b.Status, b.UserName, b.Show.Movie.Title, seatNames, b.TotalAmount)
}
