package main

import "fmt"

func printClassDiagram() {
	fmt.Println(`classDiagram

    class Genre {
        <<enumeration>>
        Action
        Comedy
        Drama
        Horror
        SciFi
        func (g Genre) String() string
    }

    class Movie {
        ID       string
        Title    string
        Duration int
        Genre    Genre
        Rating   float64
        func NewMovie(id, title string, duration int, genre Genre, rating float64) *Movie
        func (m *Movie) String() string
    }

    class SeatType {
        <<enumeration>>
        Regular
        Premium
        VIP
        func (s SeatType) String() string
    }

    class Seat {
        ID     int
        Row    string
        Number int
        Type   SeatType
        func NewSeat(id int, row string, number int, seatType SeatType) *Seat
        func (s *Seat) String() string
    }

    class Screen {
        ID           int
        ScreenNumber int
        Seats        []*Seat
        func NewScreen(id, screenNumber int, seats []*Seat) *Screen
        func (s *Screen) GetSeatsByType(seatType SeatType) []*Seat
    }

    class Theatre {
        ID      string
        Name    string
        City    string
        Screens []*Screen
        func NewTheatre(id, name, city string) *Theatre
        func (t *Theatre) AddScreen(screen *Screen)
        func (t *Theatre) String() string
    }

    class Show {
        ID               string
        Movie            *movie.Movie
        Screen           *theatre.Screen
        TheatreID        string
        StartTime        time.Time
        SeatAvailability map[int]bool
        func NewShow(id string, m *movie.Movie, screen *theatre.Screen, theatreID string, startTime time.Time) *Show
        func (s *Show) GetAvailableSeats() []*theatre.Seat
        func (s *Show) BookSeats(seatIDs []int) error
        func (s *Show) CancelSeats(seatIDs []int) error
        func (s *Show) String() string
    }

    class BookingStatus {
        <<enumeration>>
        Pending
        Confirmed
        Cancelled
        func (s BookingStatus) String() string
    }

    class Booking {
        ID          string
        UserName    string
        Show        *show.Show
        Seats       []*theatre.Seat
        TotalAmount float64
        Status      BookingStatus
        BookingTime time.Time
        func NewBooking(id, userName string, s *show.Show, seats []*theatre.Seat) *Booking
        func (b *Booking) CalculateTotal() float64
        func (b *Booking) Cancel() error
        func (b *Booking) String() string
    }

    class PricingStrategy {
        <<interface>>
        CalculatePrice(seatType theatre.SeatType) float64
        GetMultiplier() float64
    }

    class regularPricing {
        func (r *regularPricing) CalculatePrice(seatType theatre.SeatType) float64
        func (r *regularPricing) GetMultiplier() float64
    }

    class weekendPricing {
        func (w *weekendPricing) CalculatePrice(seatType theatre.SeatType) float64
        func (w *weekendPricing) GetMultiplier() float64
    }

    class premiumPricing {
        func (p *premiumPricing) CalculatePrice(seatType theatre.SeatType) float64
        func (p *premiumPricing) GetMultiplier() float64
    }

    class BookingService {
        <<interface>>
        SearchMovies(title string) []*movie.Movie
        GetShows(movieID, city string) []*show.Show
        BookTickets(userName, showID string, seatIDs []int) (*booking.Booking, error)
        CancelBooking(bookingID string) error
    }

    class AdminService {
        <<interface>>
        AddTheatre(theatre *theatre.Theatre)
        AddScreen(theatreID string, screen *theatre.Screen) error
        AddMovie(movie *movie.Movie)
        AddShow(id string, movieID, theatreID string, screenNumber int, startTime time.Time) (*show.Show, error)
        RemoveShow(showID string) error
    }

    class BookMyShow {
        mu             sync.Mutex
        Theatres       map[string]*theatre.Theatre
        Movies         map[string]*movie.Movie
        Shows          map[string]*show.Show
        Bookings       map[string]*booking.Booking
        bookingCounter int
        func GetInstance() *BookMyShow
        func (bms *BookMyShow) SearchMovies(title string) []*movie.Movie
        func (bms *BookMyShow) GetShows(movieID, city string) []*show.Show
        func (bms *BookMyShow) BookTickets(userName, showID string, seatIDs []int) (*booking.Booking, error)
        func (bms *BookMyShow) CancelBooking(bookingID string) error
        func (bms *BookMyShow) AddTheatre(t *theatre.Theatre)
        func (bms *BookMyShow) AddScreen(theatreID string, screen *theatre.Screen) error
        func (bms *BookMyShow) AddMovie(m *movie.Movie)
        func (bms *BookMyShow) AddShow(id string, movieID, theatreID string, screenNumber int, startTime time.Time) (*show.Show, error)
        func (bms *BookMyShow) RemoveShow(showID string) error
    }

    Movie --> Genre
    Seat --> SeatType
    Screen *-- Seat : contains
    Theatre *-- Screen : contains
    Show --> Movie : references
    Show --> Screen : references
    Show o-- Seat : tracks availability
    Booking --> Show : references
    Booking --> Seat : references
    Booking --> BookingStatus
    Booking ..> PricingStrategy : uses
    regularPricing ..|> PricingStrategy : implements
    weekendPricing ..|> PricingStrategy : implements
    premiumPricing ..|> PricingStrategy : implements
    BookMyShow ..|> BookingService : implements
    BookMyShow ..|> AdminService : implements
    BookMyShow o-- Theatre : manages
    BookMyShow o-- Movie : manages
    BookMyShow o-- Show : manages
    BookMyShow o-- Booking : tracks`)
}
