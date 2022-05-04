package cineplex

import (
	"log"
	"scraper/converter"
	"scraper/dto"
	"scraper/storage/model"
)

type FilmStorageRepository interface {
	IsExists(film model.Film) bool
	Insert(film model.Film) (model.Film, error)
}

type Soup interface {
	GetMovies() ([]dto.Film, error)
}

type Scraper struct {
	r FilmStorageRepository
	s Soup
}

func NewScraper(repo FilmStorageRepository, sp Soup) (*Scraper, error) {
	scraper := Scraper{
		r: repo,
		s: sp,
	}

	return &scraper, nil
}

func (c *Scraper) GetAllFilms() ([]model.Film, error) {
	var response []model.Film

	logger := log.Default()

	rawMovies, err := c.s.GetMovies()
	if err != nil {
		logger.Println("unable to scrape movies", err)
		return []model.Film{}, err
	}

	for _, mv := range rawMovies {

		filmModel, err := converter.FromDTO(mv)
		if err != nil {
			logger.Println("Unable to convert DTO to model", err)
			continue
		}

		if c.r.IsExists(filmModel) {
			continue
		}

		_, err = c.r.Insert(filmModel)
		if err != nil {
			logger.Println("unable to insert new movie", mv, err)
			continue
		}

		response = append(response, filmModel)
	}

	return response, nil
}
