package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"scraper/service"
	"scraper/storage/repository"
	"time"

	"github.com/mailjet/mailjet-apiv3-go/v3"

	_ "github.com/mattn/go-sqlite3"

	_ "embed"
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

type runtimeConfig struct {
	mailjetPubKey     string
	mailJetPrivateKey string
	sqlitePath        string
}

var runtimeConf runtimeConfig

func init() {
	runtimeConf = runtimeConfig{
		mailjetPubKey:     os.Getenv("MJ_APIKEY_PUBLIC"),
		mailJetPrivateKey: os.Getenv("MJ_APIKEY_PRIVATE"),
		sqlitePath:        os.Getenv("SQLITE_PATH"),
	}
}

// add cli help flag that describes exit codes
func main() {

	db, err := sql.Open("sqlite3", "file:../../pirata.db")
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	filmRepo := repository.NewFilmStorageRepository(db)
	emailRepo := repository.NewSubscriberRepository(db)

	scraperClient := getHttpClientWithCookies()

	scraperService, err := service.NewCineplexScraper(scraperClient, filmRepo)
	if err != nil {
		log.Fatalln(err)
	}

	imdbService := service.NewIMDB(scraperClient)

	mailjetClient := mailjet.NewMailjetClient(runtimeConf.mailjetPubKey, runtimeConf.mailJetPrivateKey)

	mailerService := service.NewMailer(mailjetClient, service.MailerConfig{
		FromEmail: os.Getenv("FROM_EMAIL"),
		FromName:  os.Getenv("FROM_NAME"),
	}, emailRepo)

	// EXECUTION below

	process(scraperService, imdbService, mailerService)

	//ticker := time.NewTicker(time.Second * 1)
	//defer ticker.Stop()
	//
	//for {
	//	select {
	//	// case to stop loop in case sigkill received
	//	case <-ticker.C:
	//		fmt.Println("processing....")
	//		process(scraperService, imdbService, mailerService)
	//
	//	}
	//}

}

func process(scraper *service.CineplexScraper, imdb *service.IMDB, mailer service.Sender) {

	allNewFilms, err := scraper.GetAllFilms()
	if err != nil {
		log.Fatalln(err)
	}

	emailFilms := imdb.EnrichFilms(allNewFilms)

	if len(emailFilms) > 0 {
		err = mailer.Send(emailFilms)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getHttpClientWithCookies() *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	return &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           jar,
		Timeout:       time.Second * 30,
	}
}
