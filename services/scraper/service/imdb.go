package service

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"common/dto"
	"common/proto"
	"scraper/converter"
	scraper_dto "scraper/dto"

	"github.com/sirupsen/logrus"
)

type IMDB struct {
	cl proto.IMDBClient
	l  logrus.FieldLogger
}

func NewIMDB(c proto.IMDBClient, log logrus.FieldLogger) *IMDB {
	return &IMDB{
		cl: c,
		l:  log,
	}
}

func (im *IMDB) GetFilms(ctx context.Context, films <-chan scraper_dto.FilmResponse) <-chan dto.EmailFilm {
	res := make(chan dto.EmailFilm)

	var wg sync.WaitGroup

	for f := range films {
		wg.Add(1)
		go func(f scraper_dto.FilmResponse) {
			defer wg.Done()

			// so normally one should make more intelligent decision about what to do with error
			// but here you can do not much with it. Are you?
			if f.Error != nil {
				return
			}

			imdbFilmData, err := im.getDataWithRetry(ctx, f.Film.Title)
			if err != nil {
				im.l.Errorf("could not get IMDB info for %s film, error : %v", f.Film.Title, err)
			}

			data := dto.IMDBData{
				Poster:  imdbFilmData.Poster,
				Plot:    imdbFilmData.Plot,
				Runtime: imdbFilmData.Runtime,
				Genres:  imdbFilmData.Genres,
			}

			res <- dto.NewEmailFilm(converter.FromModel(f.Film), data)
		}(f)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}

// exponential backoff with jitter
func (im *IMDB) getDataWithRetry(ctx context.Context, t string) (*proto.Film, error) {
	var sleepTime int
	var err error

	protoReq := &proto.FilmTitle{Title: t}

	for i := 0; i < 3; i++ {
		film, err := im.cl.GetFilm(ctx, protoReq)
		if err == nil {
			return film, err
		}

		jitter := time.Millisecond * time.Duration(rand.Intn(900))
		time.Sleep(time.Second*1<<sleepTime + jitter)
		sleepTime++
	}

	return nil, err
}
