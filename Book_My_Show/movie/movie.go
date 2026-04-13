package movie

import "fmt"

type Genre int

const (
	Action Genre = iota
	Comedy
	Drama
	Horror
	SciFi
)

func (g Genre) String() string {
	switch g {
	case Action:
		return "Action"
	case Comedy:
		return "Comedy"
	case Drama:
		return "Drama"
	case Horror:
		return "Horror"
	case SciFi:
		return "Sci-Fi"
	default:
		return "Unknown"
	}
}

type Movie struct {
	ID       string
	Title    string
	Duration int
	Genre    Genre
	Rating   float64
}

func NewMovie(id, title string, duration int, genre Genre, rating float64) *Movie {
	return &Movie{
		ID:       id,
		Title:    title,
		Duration: duration,
		Genre:    genre,
		Rating:   rating,
	}
}

func (m *Movie) String() string {
	return fmt.Sprintf("%s (%s, %dmin, %.1f★)", m.Title, m.Genre, m.Duration, m.Rating)
}
