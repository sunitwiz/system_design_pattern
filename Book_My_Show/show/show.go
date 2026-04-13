package show

import (
	"book_my_show/movie"
	"book_my_show/theatre"
	"fmt"
	"time"
)

type Show struct {
	ID               string
	Movie            *movie.Movie
	Screen           *theatre.Screen
	TheatreID        string
	StartTime        time.Time
	SeatAvailability map[int]bool
}

func NewShow(id string, m *movie.Movie, screen *theatre.Screen, theatreID string, startTime time.Time) *Show {
	availability := make(map[int]bool)
	for _, seat := range screen.Seats {
		availability[seat.ID] = true
	}
	return &Show{
		ID:               id,
		Movie:            m,
		Screen:           screen,
		TheatreID:        theatreID,
		StartTime:        startTime,
		SeatAvailability: availability,
	}
}

func (s *Show) GetAvailableSeats() []*theatre.Seat {
	var available []*theatre.Seat
	for _, seat := range s.Screen.Seats {
		if s.SeatAvailability[seat.ID] {
			available = append(available, seat)
		}
	}
	return available
}

// Validates all seats before booking any — prevents partial bookings.
func (s *Show) BookSeats(seatIDs []int) error {
	for _, id := range seatIDs {
		avail, exists := s.SeatAvailability[id]
		if !exists {
			return fmt.Errorf("seat %d does not exist in this show", id)
		}
		if !avail {
			return fmt.Errorf("seat %d is already booked", id)
		}
	}
	for _, id := range seatIDs {
		s.SeatAvailability[id] = false
	}
	return nil
}

func (s *Show) CancelSeats(seatIDs []int) error {
	for _, id := range seatIDs {
		if _, exists := s.SeatAvailability[id]; !exists {
			return fmt.Errorf("seat %d does not exist in this show", id)
		}
	}
	for _, id := range seatIDs {
		s.SeatAvailability[id] = true
	}
	return nil
}

func (s *Show) String() string {
	available := len(s.GetAvailableSeats())
	total := len(s.Screen.Seats)
	return fmt.Sprintf("Show[%s] %s | Screen %d | %s | %d/%d seats available",
		s.ID, s.Movie.Title, s.Screen.ScreenNumber,
		s.StartTime.Format("Mon 02-Jan 15:04"), available, total)
}
