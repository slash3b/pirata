package service_test

import (
	"context"
	"testing"

	"common/proto"

	"github.com/golang/mock/gomock"

	scraper_dto "scraper/dto"
	"scraper/service"

	//"scraper/service/dto"
	"scraper/storage/model"

	mock_service "scraper/service/mock"

	"github.com/sirupsen/logrus/hooks/test"
)

//go:generate mockgen -source ../../common/proto/imdb_grpc.pb.go  -destination ./mock/imdb_mock.go -package mock_service

func TestGetFilms(t *testing.T) {
	gomockCtrl := gomock.NewController(t)

	responses := make(chan scraper_dto.FilmResponse, 1)
	go func() {
		defer close(responses)
		responses <- scraper_dto.FilmResponse{
			Film: model.Film{Title: "fooo"},
		}
	}()

	imdbClient := mock_service.NewMockIMDBClient(gomockCtrl)
	imdbClient.EXPECT().GetFilm(gomock.Any(), gomock.Any()).Return(&proto.Film{Plot: "plot"}, nil).Times(1)
	nullLogger, _ := test.NewNullLogger()
	imdbSrv := service.NewIMDB(imdbClient, nullLogger)

	film := <-imdbSrv.GetFilms(context.Background(), responses)

	if film.Description != "plot" {
		t.Error("something terrible happened")
	}
}
