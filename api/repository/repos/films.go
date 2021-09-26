package repos

import (
	"context"

	"github.com/slash3b/pirata/api/model"
	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context) []model.Film
}

type FilmsRepository struct {
	conn *gorm.DB
}

func NewFilmsRepository(c *gorm.DB) *FilmsRepository {
	return &FilmsRepository{conn: c}
}

func (repo *FilmsRepository) GetAll(ctx context.Context) []model.Film {
	var films []model.Film

	select {
	case <-ctx.Done():
		return nil
	default:
		result := repo.conn.Order("register_date desc").Find(&films)
		if result.Error != nil {
			panic(result.Error)
		}
	}

	return films
}
