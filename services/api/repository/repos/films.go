package repos

import (
	"context"
	"fmt"

	"github.com/slash3b/pirata/api/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context) []model.Film
}

type FilmsRepository struct {
	db *gorm.DB
}

func NewFilmsRepository(c *gorm.DB) *FilmsRepository {
	return &FilmsRepository{db: c}
}

func (fr *FilmsRepository) GetAll(ctx context.Context) ([]model.Film, error) {
	var films []model.Film

	select {
	case <-ctx.Done():
		return films, fmt.Errorf("cancelled context")
	default:
		result := fr.db.Limit(10).Order("id desc").Find(&films)
		if result.Error != nil {
			return films, result.Error
		}
	}

	return films, nil
}
