package bookmyshow

import (
	"book_my_show/booking"
	"book_my_show/movie"
	"book_my_show/show"
	"book_my_show/theatre"
	"time"
)

type BookingService interface {
	SearchMovies(title string) []*movie.Movie
	GetShows(movieID, city string) []*show.Show
	BookTickets(userName, showID string, seatIDs []int) (*booking.Booking, error)
	CancelBooking(bookingID string) error
}

type AdminService interface {
	AddTheatre(theatre *theatre.Theatre)
	AddScreen(theatreID string, screen *theatre.Screen) error
	AddMovie(movie *movie.Movie)
	AddShow(id string, movieID, theatreID string, screenNumber int, startTime time.Time) (*show.Show, error)
	RemoveShow(showID string) error
}
