package dto

import "scraper/storage/model"

// RawFilmData is film info scraped from scraper
type RawFilmData struct {
	Title string
	Lang  string
	Date  string
}

// FilmResponse basically a container to hold converted
// FilmModel and error
type FilmResponse struct {
	Film  model.Film
	Error error
}
