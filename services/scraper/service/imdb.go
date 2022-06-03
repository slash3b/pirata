package service

import (
	"common/dto"
	"common/proto"
	"context"
	"fmt"
	"scraper/converter"
	"scraper/storage/model"
	"sync"
)

type IMDB struct {
	cl proto.IMDBClient
}

func NewIMDB(c proto.IMDBClient) *IMDB {
	return &IMDB{cl: c}
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
				fmt.Println(err) // how to deal with error here ?
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
