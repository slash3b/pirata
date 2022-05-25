package service_test

import (
	"context"
	"scraper/dto"
	"scraper/service"
	mock_service "scraper/service/mock"
	"scraper/storage/model"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -source ./scraper.go -destination ./mock/scraper_mock.go

func TestScraper_GetAllFilms_WithMixedResults(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockSoup := mock_service.NewMockSoup(mockCtrl)
	mockRepo := mock_service.NewMockFilmStorageRepository(mockCtrl)

	mockSoup.EXPECT().GetMovies().Return([]dto.RawFilmData{
		{Title: "Foo Bar ",
			Lang: "KO",
			Date: "(10.03.2022)",
		},
		{Title: "Foo Bar ",
			Lang: "EN",
			Date: "(10.03.2022)",
		},
	}, nil).Times(1)

	gomock.InOrder(
		mockRepo.EXPECT().IsExists(gomock.Any()).Return(false).Times(1),
		mockRepo.EXPECT().IsExists(gomock.Any()).Return(true).Times(1),
	)

	mockRepo.EXPECT().Insert(gomock.Any()).Return(model.Film{}, nil).Times(1)

	scraper := service.NewScraper(mockRepo, mockSoup)

	_, films := scraper.GetFilms(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		assert.Equal(t, "Foo Bar", (<-films).Title)
		// todo: mock time to be able to compare complete DTO struct
		// fixme: apparently this test became flaky, needs fixing
	}()

	wg.Wait()

}
