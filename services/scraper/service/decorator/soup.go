package decorator

import (
	"net/http"
	"scraper/dto"
	"strings"

	"github.com/anaskhan96/soup"
)

type SoupDecorator struct {
	c *http.Client
}

func NewSoupDecorator(client *http.Client) *SoupDecorator {
	return &SoupDecorator{
		c: client,
	}
}

func (s *SoupDecorator) GetMovies() ([]dto.RawFilmData, error) {
	var response []dto.RawFilmData

	_, err := soup.GetWithClient("https://cineplex.md/lang/en", s.c)
	if err != nil {
		return response, err
	}

	out, err := soup.GetWithClient("https://cineplex.md/films", s.c)
	if err != nil {
		return response, err
	}

	root := soup.HTMLParse(out)

	allMovies := root.FindAll("div", "class", "movies_blcks")

	for _, movie := range allMovies {
		title := strings.TrimSpace(movie.Find("h3", "class", "overlay__title").Text())
		lang := strings.TrimSpace(movie.Find("span", "class", "overlay__lang").Text())
		startDate := strings.TrimSpace(movie.Find("div", "class", "startdate").Text())

		filmDTO := dto.RawFilmData{
			Title: title,
			Lang:  lang,
			Date:  startDate,
		}

		response = append(response, filmDTO)
	}

	return response, nil
}
