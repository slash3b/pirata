package publisher

import "common/dto"

type Mail struct {
}

func NewMailPublisher() *Mail {
	return &Mail{}
}

func (m *Mail) Send(mailFilms <-chan dto.EmailFilm) error {

	var allFilms []dto.EmailFilm
	for i := range mailFilms {
		allFilms = append(allFilms, i)
	}

	// here should just send it to kafka

	return nil
}
