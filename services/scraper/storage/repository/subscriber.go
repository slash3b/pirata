package repository

import (
	"database/sql"
	"log"
	"scraper/storage/model"
)

type SubscriberRepository struct {
	db *sql.DB
}

func NewSubscriberRepository(db *sql.DB) *SubscriberRepository {
	return &SubscriberRepository{
		db: db,
	}
}

func (r *SubscriberRepository) GetAllActive() ([]model.Subscriber, error) {

	var subs []model.Subscriber

	rows, err := r.db.Query("select * from subscribers where subscribed=1")

	if err != nil {
		return subs, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
	}(rows)

	for rows.Next() {
		var email string
		var name string
		var subscribed bool
		err = rows.Scan(&email, &name, &subscribed)
		if err != nil {
			return subs, err
		}

		subs = append(subs, model.Subscriber{
			Email:      email,
			Name:       name,
			Subscribed: subscribed,
		})
	}

	return subs, nil
}
