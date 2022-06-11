package service

import (
	"context"
	"log"
	"scraper/converter"
	"scraper/dto"
	"scraper/metrics"
	"scraper/storage/model"
	"sync"
)

type FilmStorageRepository interface {
	IsExists(film model.Film) bool
	Insert(film model.Film) (model.Film, error)
}

type Soup interface {
	GetMovies() ([]dto.RawFilmData, error)
}

type Scraper struct {
	r FilmStorageRepository
	s Soup
}

func NewScraper(repo FilmStorageRepository, sp Soup) *Scraper {
	return &Scraper{
		r: repo,
		s: sp,
	}
}

func (c *Scraper) GetFilms(ctx context.Context) (context.Context, <-chan model.Film) {
	response := make(chan model.Film)

	// todo: decide on logger and make it injectable
	logger := log.Default()

	rawMovies, err := c.s.GetMovies()
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_scrape_movies").Inc()
		logger.Println("unable to scrape movies", err)

		close(response)
		return ctx, response
	}

	var wg sync.WaitGroup
	for _, rawMovie := range rawMovies {

		wg.Add(1)
		go func(movieData dto.RawFilmData) {
			defer wg.Done()
			movieModel, err := converter.FromDTO(movieData)
			if err != nil {
				metrics.ScraperErrors.WithLabelValues("unable_to_convert_raw_movie_dto").Inc()
				logger.Println("Unable to convert raw movie DTO", err)

				return
			}

			if !c.r.IsExists(movieModel) {
				_, err = c.r.Insert(movieModel)
				if err != nil {
					logger.Println("unable to insert new rawMovie", movieData, err)
					return
				}

				select {
				case <-ctx.Done():
				case response <- movieModel:
				}
			}
		}(rawMovie)

	}

	// goroutine that waits for all others to complete and safely close the channel
	go func() {
		wg.Wait()
		close(response)
	}()

	return ctx, response
}
