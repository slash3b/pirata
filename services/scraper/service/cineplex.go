package service

import (
	"net/http"
	"scraper/converter"
	"scraper/dto"
	"scraper/storage/model"
	"scraper/storage/repository"
	"strings"

	"github.com/anaskhan96/soup"
)

type CineplexScraper struct {
	c *http.Client
	r repository.StorageRepository
}

func NewCineplexScraper(httpClient *http.Client, repo repository.StorageRepository) (*CineplexScraper, error) {

	scraper := CineplexScraper{
		c: httpClient,
		r: repo,
	}

	return &scraper, nil
}

func (c *CineplexScraper) GetAllFilms() ([]model.Film, error) {
	var response []model.Film

	out, err := soup.GetWithClient("https://cineplex.md/lang/en", c.c)
	if err != nil {
		return response, err
	}

	out, err = soup.GetWithClient("https://cineplex.md/films", c.c)
	if err != nil {
		return response, err
	}

	root := soup.HTMLParse(out)

	allMovies := root.FindAll("div", "class", "movies_blcks")
	for _, movie := range allMovies {

		title := strings.TrimSpace(movie.Find("h3", "class", "overlay__title").Text()) // title
		lang := strings.TrimSpace(movie.Find("span", "class", "overlay__lang").Text()) // lang
		startDate := strings.TrimSpace(movie.Find("div", "class", "startdate").Text()) // startdate

		filmDTO := dto.Film{
			Title: title,
			Lang:  lang,
			Date:  startDate,
		}

		filmModel, err := converter.FromDTO(filmDTO)
		if err != nil {
			// todo log with DTO
			continue
		}

		// check if film exists
		if c.r.IsExists(filmModel) {
			//fmt.Println(c.r.Insert(filmModel))
			//return response, nil
			continue
		}

		//c.r.IsExists()

		response = append(response, filmModel)
	}

	return response, nil
}
