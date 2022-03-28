package repository

import (
	"database/sql"
	"fmt"
	"scraper/storage/model"
)

/*
use Driver interface for testing
*/

type StorageRepository interface {
	IsExists(film model.Film) bool
	Insert(film model.Film) (model.Film, error)
	//GetBy(dto dto.Film) model.Film
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) StorageRepository {
	return &Repository{db: db}
}

func (r *Repository) IsExists(film model.Film) bool {

	q := fmt.Sprintf("select exists(select 1 from films where title='%s' and lang='%s' and dimension='%s')", film.Title, film.Lang, film.Dimension)
	row := r.db.QueryRow(q)

	var exists int
	err := row.Scan(&exists)
	if err != nil {
		return true
	}

	if exists > 0 {
		return true
	}

	return false
}

func (r *Repository) Insert(film model.Film) (model.Film, error) {

	q := fmt.Sprintf(`insert into films (title , dimension , lang , release_date ) values("%s", "%s", "%s", "%s")`, film.Title, film.Dimension, film.Lang, film.StartDate.String())

	res, err := r.db.Exec(q)
	if err != nil {
		return model.Film{}, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		return model.Film{}, nil
	}
	film.ID = int(id)

	return film, nil
}

func (r *Repository) GetBy(fm model.Film) (model.Film, error) {

	var fi model.Film

	row := r.db.QueryRow(fmt.Sprintf("select * from films where title='%s' and dimension='%s' and lang='%s'", fm.Title, fm.Dimension, fm.Lang))

	var dbDateString string

	err := row.Scan(&fi.ID, &fi.Title, &fi.Lang, &fi.Dimension, &dbDateString)
	if err != nil {
		return fi, err

	}

	// todo: parse dbDateString

	fmt.Println("2", fi)

	return fi, nil
}
