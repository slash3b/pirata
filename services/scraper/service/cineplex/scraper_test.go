package cineplex_test

import (
	"scraper/dto"
	"scraper/service/cineplex"
	mock_cineplex "scraper/service/cineplex/mock"
	"scraper/storage/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -source ./scraper.go -destination ./mock/scraper_mock.go

func TestScraper_GetAllFilms_WithMixedResults(t *testing.T) {
	mockCtrl := gomock.NewController(t)

	mockSoup := mock_cineplex.NewMockSoup(mockCtrl)
	mockRepo := mock_cineplex.NewMockFilmStorageRepository(mockCtrl)

	mockSoup.EXPECT().GetMovies().Return([]dto.Film{
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

	scraper, err := cineplex.NewScraper(mockRepo, mockSoup)
	assert.NoError(t, err)

	films, err := scraper.GetAllFilms()
	assert.NoError(t, err)
	assert.Len(t, films, 1)
}
