package service

import (
	"net/http"
	"scraper/converter"
	"scraper/dto"
	"scraper/storage/model"
	"scraper/storage/repository"
	"strings"
	"sync"

	"github.com/anaskhan96/soup"
)

type CineplexScraper struct {
	c *http.Client
	r repository.FilmStorageRepository
}

func NewCineplexScraper(httpClient *http.Client, repo repository.FilmStorageRepository) (*CineplexScraper, error) {

	scraper := CineplexScraper{
		c: httpClient,
		r: repo,
	}

	return &scraper, nil
}

func (c *CineplexScraper) GetAllFilms() ([]model.Film, error) {
	var response []model.Film

	out, err := soup.GetWithClient("https://cineplex.md/lang/en", c.c)
	if err != nil {
		return response, err
	}

	out, err = soup.GetWithClient("https://cineplex.md/films", c.c)
	if err != nil {
		return response, err
	}

	root := soup.HTMLParse(out)

	allMovies := root.FindAll("div", "class", "movies_blcks")

	var wg sync.WaitGroup

	for _, movie := range allMovies {
		wg.Add(1)

		go func(m soup.Root) {
			defer wg.Done()

			title := strings.TrimSpace(m.Find("h3", "class", "overlay__title").Text()) // title
			lang := strings.TrimSpace(m.Find("span", "class", "overlay__lang").Text()) // lang
			startDate := strings.TrimSpace(m.Find("div", "class", "startdate").Text()) // startdate

			filmDTO := dto.Film{
				Title: title,
				Lang:  lang,
				Date:  startDate,
			}

			filmModel, err := converter.FromDTO(filmDTO)
			if err != nil {
				// todo log with DTO
				return
			}

			// check if film exists
			if c.r.IsExists(filmModel) {
				return
			}

			c.r.Insert(filmModel)

			//c.r.IsExists()

			response = append(response, filmModel)
		}(movie)
	}

	wg.Wait()

	return response, nil
}
