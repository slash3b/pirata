package repository

import (
	"database/sql"
	"scraper/storage/model"

	"github.com/sirupsen/logrus"
)

type SubscriberRepository struct {
	db *sql.DB
	l  logrus.FieldLogger
}

func NewSubscriberRepository(log logrus.FieldLogger, db *sql.DB) *SubscriberRepository {
	return &SubscriberRepository{
		db: db,
		l:  log,
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
			r.l.Errorf("could not close *sql.Rows properly: %v", err)
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
