package service

import (
	"common/dto"
	"common/proto"
	"context"
	"scraper/converter"
	"scraper/storage/model"
	"sync"

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

func (i *IMDB) GetFilms(ctx context.Context, films <-chan model.Film) <-chan dto.EmailFilm {
	res := make(chan dto.EmailFilm)

	var wg sync.WaitGroup

	for f := range films {

		wg.Add(1)
		go func(f model.Film) {
			defer wg.Done()

			fr, err := i.cl.GetFilm(ctx, &proto.FilmTitle{Title: f.Title})
			if err != nil {
				i.l.Errorf("could not get IMDB info for %s film, error : %v", f.Title, err)
				return
			}

			data := dto.IMDBData{
				Poster:  fr.Poster,
				Plot:    fr.Plot,
				Runtime: fr.Runtime,
				Genres:  fr.Genres,
			}

			res <- dto.NewEmailFilm(converter.FromModel(f), data)
		}(f)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
