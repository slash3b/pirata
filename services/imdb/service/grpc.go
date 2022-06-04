package service

import (
	"context"
	"fmt"

	"common/proto"
)

type GRPCServer struct {
	imdbService *IMDB
	proto.UnimplementedIMDBServer
}

func NewGRPCServer(i *IMDB) GRPCServer {
	return GRPCServer{imdbService: i}
}

func (i GRPCServer) GetFilm(ctx context.Context, title *proto.FilmTitle) (*proto.Film, error) {

	for filmData := range i.imdbService.FindFilms(ctx, []string{title.GetTitle()}) {

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
