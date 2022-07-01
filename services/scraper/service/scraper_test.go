package service_test

import (
	"context"
	"scraper/dto"
	"scraper/service"
	mock_service "scraper/service/mock"
	"scraper/storage/model"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -source ./scraper.go -destination ./mock/scraper_mock.go

func TestScraper_GetAllFilms_WithNewResults(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockSoup := mock_service.NewMockMovieFetcher(mockCtrl)
	mockRepo := mock_service.NewMockFilmStorageRepository(mockCtrl)

	mockSoup.EXPECT().GetMovies().Return([]dto.RawFilmData{
		{Title: "Foo Bar ",
			Lang: "KO",
			Date: "(10.03.2022)",
		},
		{Title: "Moo Bar ",
			Lang: "EN",
			Date: "(10.03.2022)",
		},
	}, nil).Times(1)

	mockRepo.EXPECT().IsExists(gomock.Any()).Return(false).Times(2)

	mockRepo.EXPECT().Insert(gomock.Any()).Return(model.Film{}, nil).Times(2)

	testLogger, _ := test.NewNullLogger() // not sure I need second param for now
	scraper := service.NewScraper(mockRepo, mockSoup, testLogger)

	_, films := scraper.GetFilms(context.Background())

	var filmsCollection []model.Film
	for f := range films {
		filmsCollection = append(filmsCollection, f)
	}

	location, err := time.LoadLocation("Europe/Chisinau")
	assert.NoError(t, err)

	expectedFilms := []model.Film{
		{
			ID:        0,
			Title:     "Foo Bar",
			Lang:      "KO",
			Dimension: "2D",
			StartDate: time.Date(2022, time.March, 10, 0, 0, 0, 0, location),
		},
		{
			ID:        0,
			Title:     "Moo Bar",
			Lang:      "EN",
			Dimension: "2D",
			StartDate: time.Date(2022, time.March, 10, 0, 0, 0, 0, location),
		},
	}

	assert.Len(t, filmsCollection, 2)

	for _, c := range expectedFilms {
		assert.Contains(t, filmsCollection, c)
	}
}

func TestScraper_GetAllFilms_WithAlreadyExistingResults(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockSoup := mock_service.NewMockMovieFetcher(mockCtrl)
	mockRepo := mock_service.NewMockFilmStorageRepository(mockCtrl)

	mockSoup.EXPECT().GetMovies().Return([]dto.RawFilmData{
		{Title: "Foo Bar ",
			Lang: "KO",
			Date: "(10.03.2022)",
		},
		{Title: "Moo Bar ",
			Lang: "EN",
			Date: "(10.03.2022)",
		},
	}, nil).Times(1)

	mockRepo.EXPECT().IsExists(gomock.Any()).Return(true).Times(2)

	testLogger, _ := test.NewNullLogger()
	scraper := service.NewScraper(mockRepo, mockSoup, testLogger)

	_, films := scraper.GetFilms(context.Background())

	var filmsCollection []model.Film
	for f := range films {
		filmsCollection = append(filmsCollection, f)
	}

	assert.Empty(t, filmsCollection)
}
