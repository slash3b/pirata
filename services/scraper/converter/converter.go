package converter

import (
	"scraper/dto"
	"scraper/storage/model"
	"strconv"
	"strings"
	"time"
)

func FromDTO(dto dto.Film) (model.Film, error) {
	is3D := strings.Contains(dto.Title, "3D")
	dimension := "2D"
	if is3D {
		dimension = "3D"
	}

	title := strings.Replace(dto.Title, "\n", " ", -1)
	title = strings.TrimSpace(title)
	title = strings.TrimRight(title, "3D")
	title = strings.TrimRight(title, "2D")
	title = strings.TrimSpace(title)

	lang := strings.Replace(dto.Lang, "\n", " ", -1)
	lang = strings.TrimFunc(lang, func(r rune) bool {
		return r == '(' || r == ')' || r == ' '
	})

	date := strings.TrimFunc(dto.Date, func(r rune) bool {
		return r == '(' || r == ')'
	})

	dateChunks := strings.Split(date, ".")

	day, err := strconv.Atoi(dateChunks[0])
	if err != nil {
		return model.Film{}, err
	}
	month, err := strconv.Atoi(dateChunks[1])
	if err != nil {
		return model.Film{}, err
	}
	year, err := strconv.Atoi(dateChunks[2])
	if err != nil {
		return model.Film{}, err
	}

	location, err := time.LoadLocation("Europe/Chisinau")
	if err != nil {
		return model.Film{}, err
	}

	startTimeDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)

	film := model.Film{
		Title:     title,
		Dimension: dimension,
		StartDate: startTimeDate,
		Lang:      lang,
	}

	return film, nil
}
