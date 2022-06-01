package dto

// EmailFilm info represent film information ready to be sent by email
// it combines IMDBData and FilmData
type EmailFilm struct {
	CompositeTitle string
	Lang           string
	ReleaseDate    string
	Duration       string
	Description    string
	PosterUrl      string
	Genres         string
}

func NewEmailFilm(film FilmData, imdbData IMDBData) EmailFilm {
	return EmailFilm{
		CompositeTitle: film.Title,
		Lang:           film.Lang,
		ReleaseDate:    film.ReleaseDate,
		Duration:       imdbData.Runtime,
		Description:    imdbData.Plot,
		PosterUrl:      imdbData.Poster,
		Genres:         imdbData.Genres,
	}
}

type FilmData struct {
	Title       string
	Lang        string
	ReleaseDate string
}

type IMDBData struct {
	Poster  string
	Plot    string
	Runtime string
	Genres  string
}
