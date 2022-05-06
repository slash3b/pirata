package service

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"scraper/dto"
	"scraper/storage/model"
	"strings"
	"sync"

	"github.com/anaskhan96/soup"
)

type IMDB struct {
	c *http.Client
}

func NewIMDB(httpClient *http.Client) *IMDB {
	imdb := IMDB{
		c: httpClient,
	}

	return &imdb
}

/*
	Ideas:
	1. so do this getting concurrently
	2. implement LRU cache in case movie was already searched for !

*/

func (c *IMDB) search(title string) (soup.Root, error) {
	vals := url.Values{}
	vals.Add("q", title)

	path := fmt.Sprintf("https://www.imdb.com/find?%s", vals.Encode())

	out, err := soup.GetWithClient(path, c.c)
	if err != nil {

		return soup.Root{}, err
	}

	root := soup.HTMLParse(out)
	allMovies := root.FindAll("tr", "class", "findResult")

	if len(allMovies) == 0 {
		// if movie info is missing
		// log it here to be checked later on
		// log.With("path", path) ... attempted to search for film but could not
		return soup.Root{}, err
	}

	firstRow := allMovies[0]

	linkNode := firstRow.Find("a")
	if linkNode.Error != nil {
		return soup.Root{}, err
	}

	attributes := linkNode.Attrs()

	movieTitleLink, ok := attributes["href"]
	if !ok {
		return soup.Root{}, errors.New("could not find link to movie")
	}

	path = fmt.Sprintf("https://www.imdb.com/%s", movieTitleLink)

	res, err := soup.GetWithClient(path, c.c)
	if err != nil {
		return soup.Root{}, err
	}

	return soup.HTMLParse(res), nil
}

func (c *IMDB) EnrichFilms(films []model.Film) []dto.EmailFilm {
	var wg sync.WaitGroup

	var emailFilms []dto.EmailFilm

	for _, f := range films {
		wg.Add(1)
		go func(film model.Film) {

			defer wg.Done()

			data := c.getFilmData(film)
			emailFilms = append(emailFilms, dto.EmailFromModel(film, data))
		}(f)
	}

	wg.Wait()

	return emailFilms
}

func (c *IMDB) getFilmData(film model.Film) dto.IMDBData {

	var data dto.IMDBData

	root, err := c.search(film.Title)
	if err != nil {
		return data
	}

	// Poster
	poster := root.Find("div", "class", "ipc-poster")
	if poster.Error == nil {
		posterImgEl := poster.Find("img")
		if posterImgEl.Error == nil {
			bb := posterImgEl.Attrs()

			posterURL, ok := bb["src"]
			if ok {
				data.Poster = posterURL
			}
		}
	}

	// Plot
	plotElem := root.Find("p", "data-testid", "plot")
	if plotElem.Error == nil {
		plot := plotElem.Find("span")
		if plot.Error == nil {
			data.Plot = plot.FullText()
		}
	}

	// Runtime
	techElem := root.Find("div", "data-testid", "title-techspecs-section")
	if techElem.Error == nil {
		techSpecs := techElem.FindAll("li")

		for _, v := range techSpecs {
			if v.Error == nil {
				techSpecName := v.Find("span")
				if techSpecName.Error == nil {
					if strings.TrimSpace(strings.ToLower(techSpecName.Text())) == "runtime" {
						content := v.Find("div")
						if content.Error == nil {
							data.Runtime = content.FullText()
							break
						}
					}
				}
			}
		}
	}

	// Genres
	genreTags := []string{}
	genres := root.Find("li", "data-testid", "storyline-genres")
	if genres.Error == nil {
		results := genres.FindAll("a")

		for _, v := range results {
			if v.Error == nil {
				// check that v is okay ?
				genreTags = append(genreTags, v.FullText())
			}
		}

		data.Genres = strings.Join(genreTags, ", ")
	}

	return data
}
