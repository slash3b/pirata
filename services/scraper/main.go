package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"scraper/service"
	"scraper/storage/repository"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/*
	converter.go converter -- gets a strings parsed and transforms them to a Film Model ?
	cineplex.go service -- knows how and gets raw films and uses converter to get those

	final product of service is a batch of Film models


	I want to add those to the database
	I want to send email with films, possibly updated with more info, trailers and so on.
*/

/*
	for now
	lets have just two services:
	- this with everything inside
	- prometheus
	- grafana

	todo: figure what to do with logging
*/

func main() {

	/**
	DB stuff
	file:test.db?cache=shared&mode=memory
	*/

	db, err := sql.Open("sqlite3", "file:../../pirata.db") // todo update: pass through ENV variables
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	filmRepo := repository.NewRepository(db)

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	scraperClient := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           jar,
		Timeout:       time.Second * 30,
	}

	scraperService, err := service.NewCineplexScraper(scraperClient, filmRepo)
	if err != nil {
		log.Fatalln(err)
	}

	allNewFilms, err := scraperService.GetAllFilms()
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range allNewFilms {
		fmt.Printf("%#v \n", v)
	}

	/*
		def register(cn, title: str):
		    cursor = cn.cursor()
		    meta = {}
		    imdb = IMDb()
		    trimmed_title = _clean_title(title)

		    imdb_search_result = imdb.search_movie(trimmed_title)

		    # some films could not be found on IMDB, for instance some new Romanian film
		    if imdb_search_result:

		        movie = imdb_search_result[0]

		        imdb.update(movie)
		        infoset = ['cover url', 'rating', 'title', 'plot', 'long imdb title', 'genres', 'runtimes']

		        # update take specific object instances
		        # find out how to do it properly
		        for info in infoset:
		            meta[info] = movie.get(info)

		        # now fill in youtube data
		        meta['trailer'] = ""
	*/

	/*
		def _get_yt_trailer(title: str) -> str:
		    query_string = urllib.parse.urlencode({"search_query" : f'{title} trailer english'})
		    html_content = urllib.request.urlopen(f'http://www.youtube.com/results?{query_string}')
		    soup = BeautifulSoup(html_content.read().decode(), 'html.parser')

		    trailer = 'http://www.youtube.com' + soup.find('div', 'yt-lockup yt-lockup-tile yt-lockup-video vve-check clearfix').find('a')['href']

		    return trailer
	*/
}
