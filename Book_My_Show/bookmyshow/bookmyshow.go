package bookmyshow

import (
	"book_my_show/booking"
	"book_my_show/movie"
	"book_my_show/show"
	"book_my_show/theatre"
	"fmt"
	"strings"
	"sync"
	"time"
)

type BookMyShow struct {
	mu             sync.Mutex
	Theatres       map[string]*theatre.Theatre
	Movies         map[string]*movie.Movie
	Shows          map[string]*show.Show
	Bookings       map[string]*booking.Booking
	bookingCounter int
}

// Compile-time interface checks.
var _ BookingService = (*BookMyShow)(nil)
var _ AdminService = (*BookMyShow)(nil)

var (
	instance *BookMyShow
	once     sync.Once
)

func GetInstance() *BookMyShow {
	once.Do(func() {
		instance = &BookMyShow{
			Theatres: make(map[string]*theatre.Theatre),
			Movies:   make(map[string]*movie.Movie),
			Shows:    make(map[string]*show.Show),
			Bookings: make(map[string]*booking.Booking),
		}
	})
	return instance
}

func ResetInstance() {
	once = sync.Once{}
	instance = nil
}

// --- BookingService Implementation ---

func (bms *BookMyShow) SearchMovies(title string) []*movie.Movie {
	bms.mu.Lock()
	defer bms.mu.Unlock()

	var results []*movie.Movie
	for _, m := range bms.Movies {
		if strings.Contains(strings.ToLower(m.Title), strings.ToLower(title)) {
			results = append(results, m)
		}
	}
	return results
}

func (bms *BookMyShow) GetShows(movieID, city string) []*show.Show {
	bms.mu.Lock()
	defer bms.mu.Unlock()

	var results []*show.Show
	for _, s := range bms.Shows {
		if s.Movie.ID == movieID {
			t, exists := bms.Theatres[s.TheatreID]
			if exists && strings.EqualFold(t.City, city) {
				results = append(results, s)
			}
		}
	}
	return results
}

func (bms *BookMyShow) BookTickets(userName, showID string, seatIDs []int) (*booking.Booking, error) {
	bms.mu.Lock()
	defer bms.mu.Unlock()

	s, exists := bms.Shows[showID]
	if !exists {
		return nil, fmt.Errorf("show %s not found", showID)
	}

	if err := s.BookSeats(seatIDs); err != nil {
		return nil, err
	}

	seats := make([]*theatre.Seat, 0, len(seatIDs))
	for _, seatID := range seatIDs {
		for _, seat := range s.Screen.Seats {
			if seat.ID == seatID {
				seats = append(seats, seat)
				break
			}
		}
	}

	bms.bookingCounter++
	bookingID := fmt.Sprintf("BK-%04d", bms.bookingCounter)
	b := booking.NewBooking(bookingID, userName, s, seats)
	bms.Bookings[bookingID] = b
	return b, nil
}

func (bms *BookMyShow) CancelBooking(bookingID string) error {
	bms.mu.Lock()
	defer bms.mu.Unlock()

	b, exists := bms.Bookings[bookingID]
	if !exists {
		return fmt.Errorf("booking %s not found", bookingID)
	}
	return b.Cancel()
}

// --- AdminService Implementation ---

func (bms *BookMyShow) AddTheatre(t *theatre.Theatre) {
	bms.mu.Lock()
	defer bms.mu.Unlock()

	bms.Theatres[t.ID] = t
	fmt.Printf("  Added theatre: %s\n", t)
}

func (bms *BookMyShow) AddScreen(theatreID string, screen *theatre.Screen) error {
	bms.mu.Lock()
	defer bms.mu.Unlock()

	t, exists := bms.Theatres[theatreID]
	if !exists {
		return fmt.Errorf("theatre %s not found", theatreID)
	}
	t.AddScreen(screen)
	fmt.Printf("  Added Screen %d to %s (%d seats)\n", screen.ScreenNumber, t.Name, len(screen.Seats))
	return nil
}

func (bms *BookMyShow) AddMovie(m *movie.Movie) {
	bms.mu.Lock()
	defer bms.mu.Unlock()

	bms.Movies[m.ID] = m
	fmt.Printf("  Added movie: %s\n", m)
}

func (bms *BookMyShow) AddShow(id string, movieID, theatreID string, screenNumber int, startTime time.Time) (*show.Show, error) {
	bms.mu.Lock()
	defer bms.mu.Unlock()

	m, exists := bms.Movies[movieID]
	if !exists {
		return nil, fmt.Errorf("movie %s not found", movieID)
	}

	t, exists := bms.Theatres[theatreID]
	if !exists {
		return nil, fmt.Errorf("theatre %s not found", theatreID)
	}

	var screen *theatre.Screen
	for _, s := range t.Screens {
		if s.ScreenNumber == screenNumber {
			screen = s
			break
		}
	}
	if screen == nil {
		return nil, fmt.Errorf("screen %d not found in theatre %s", screenNumber, theatreID)
	}

	s := show.NewShow(id, m, screen, theatreID, startTime)
	bms.Shows[id] = s
	fmt.Printf("  Added show: %s\n", s)
	return s, nil
}

func (bms *BookMyShow) RemoveShow(showID string) error {
	bms.mu.Lock()
	defer bms.mu.Unlock()

	_, exists := bms.Shows[showID]
	if !exists {
		return fmt.Errorf("show %s not found", showID)
	}

	for _, b := range bms.Bookings {
		if b.Show.ID == showID && b.Status == booking.Confirmed {
			return fmt.Errorf("cannot remove show %s: has active bookings", showID)
		}
	}

	delete(bms.Shows, showID)
	fmt.Printf("  Removed show: %s\n", showID)
	return nil
}
