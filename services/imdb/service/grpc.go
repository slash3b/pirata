package service

import (
	"common/proto"
	"context"
	"fmt"
)

type Grpc struct {
	i *IMDB
	proto.UnimplementedIMDBServer
}

func NewIMDBGrpc(i *IMDB) Grpc {
	return Grpc{i: i}
}

func (i Grpc) GetFilm(ctx context.Context, title *proto.FilmTitle) (*proto.Film, error) {

	for filmData := range i.i.FindFilms(ctx, []string{title.GetTitle()}) {

		film := &proto.Film{
			Poster:  filmData.Poster,
			Plot:    filmData.Plot,
			Runtime: filmData.Runtime,
			Genres:  filmData.Genres,
		}

		return film, nil
	}

	return nil, fmt.Errorf("unable to find film by title")
}
