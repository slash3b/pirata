package service

import (
	"context"
	"scraper/converter"
	"scraper/dto"
	"scraper/metrics"
	"scraper/storage/model"
	"sync"

	"github.com/sirupsen/logrus"
)

type FilmStorageRepository interface {
	IsExists(film model.Film) bool
	Insert(film model.Film) (model.Film, error)
}

type MovieFetcher interface {
	GetMovies() ([]dto.RawFilmData, error)
}

type Scraper struct {
	r FilmStorageRepository
	s MovieFetcher
	l logrus.FieldLogger
}

func NewScraper(repo FilmStorageRepository, sp MovieFetcher, log logrus.FieldLogger) *Scraper {
	return &Scraper{
		r: repo,
		s: sp,
		l: log,
	}
}

func (c *Scraper) GetFilms(ctx context.Context) (context.Context, <-chan dto.FilmResponse) {
	response := make(chan dto.FilmResponse)

	rawMovies, err := c.s.GetMovies()
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_scrape_movies").Inc()
		c.l.Println("unable to scrape movies", err)

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
				c.l.Println("Unable to convert raw movie DTO", err)
				response <- dto.FilmResponse{Film: movieModel, Error: err}
				return
			}

			if !c.r.IsExists(movieModel) {
				_, err = c.r.Insert(movieModel)
				if err != nil {
					metrics.ScraperErrors.WithLabelValues("unable_to_record_movie").Inc()
					response <- dto.FilmResponse{Film: movieModel, Error: err}
					c.l.Println("unable to insert new rawMovie", movieData, err)
					return
				}

				select {
				case <-ctx.Done():
				case response <- dto.FilmResponse{Film: movieModel}:
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
