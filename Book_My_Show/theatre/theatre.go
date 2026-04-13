package theatre

import "fmt"

type SeatType int

const (
	Regular SeatType = iota
	Premium
	VIP
)

func (s SeatType) String() string {
	switch s {
	case Regular:
		return "Regular"
	case Premium:
		return "Premium"
	case VIP:
		return "VIP"
	default:
		return "Unknown"
	}
}

type Seat struct {
	ID     int
	Row    string
	Number int
	Type   SeatType
}

func NewSeat(id int, row string, number int, seatType SeatType) *Seat {
	return &Seat{
		ID:     id,
		Row:    row,
		Number: number,
		Type:   seatType,
	}
}

func (s *Seat) String() string {
	return fmt.Sprintf("%s%d(%s)", s.Row, s.Number, s.Type)
}

type Screen struct {
	ID           int
	ScreenNumber int
	Seats        []*Seat
}

func NewScreen(id, screenNumber int, seats []*Seat) *Screen {
	return &Screen{
		ID:           id,
		ScreenNumber: screenNumber,
		Seats:        seats,
	}
}

func (s *Screen) GetSeatsByType(seatType SeatType) []*Seat {
	var result []*Seat
	for _, seat := range s.Seats {
		if seat.Type == seatType {
			result = append(result, seat)
		}
	}
	return result
}

type Theatre struct {
	ID      string
	Name    string
	City    string
	Screens []*Screen
}

func NewTheatre(id, name, city string) *Theatre {
	return &Theatre{
		ID:   id,
		Name: name,
		City: city,
	}
}

func (t *Theatre) AddScreen(screen *Screen) {
	t.Screens = append(t.Screens, screen)
}

func (t *Theatre) String() string {
	return fmt.Sprintf("%s (%s) - %d screens", t.Name, t.City, len(t.Screens))
}
