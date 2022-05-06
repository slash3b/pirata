package dto

import (
	"fmt"
	"scraper/storage/model"
	"strings"
)

// EmailFilm info represent film information ready to be sent by email
type EmailFilm struct {
	CompositeTitle string
	Lang           string
	ReleaseDate    string
	Duration       string
	Description    string
	PosterUrl      string
	Genres         string
}

func EmailFromModel(film model.Film, data IMDBData) EmailFilm {

	var langEmoji string
	switch strings.ToLower(film.Lang) {
	case "ru":
		langEmoji = `🇷🇺`
	case "ro":
		langEmoji = `🇲🇩`
	case "en":
		langEmoji = `🇺🇸`
	default:
		langEmoji = strings.ToLower(film.Lang)
	}

	return EmailFilm{
		CompositeTitle: fmt.Sprintf("%s, %s", film.Title, film.Dimension),
		Lang:           langEmoji,
		ReleaseDate:    fmt.Sprintf("%s %d, %s", film.StartDate.Month(), film.StartDate.Day(), film.StartDate.Weekday()),
		Duration:       data.Runtime,
		Description:    data.Plot,
		PosterUrl:      data.Poster,
		Genres:         data.Genres,
	}
}
