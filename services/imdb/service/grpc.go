package service

import (
	"common/proto"
	"context"
	"strings"
)

type Grpc struct {
	i *IMDB
	proto.UnimplementedIMDBServer
}

func NewIMDBGrpc(i *IMDB) Grpc {
	return Grpc{i: i}
}

func (i Grpc) GetFilms(titles *proto.FilmTitles, srv proto.IMDB_GetFilmsServer) error {

	grpcTitles := strings.Split(titles.GetTitles(), ",")

	for v := range i.i.FindFilms(context.Background(), grpcTitles) {
		film := &proto.Film{
			Poster:  v.Poster,
			Plot:    v.Plot,
			Runtime: v.Runtime,
			Genres:  v.Genres,
		}

		err := srv.Send(film)
		if err != nil {
			return err
		}
	}

	return nil
}
