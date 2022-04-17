package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"scraper/service"
	"scraper/storage/repository"
	"syscall"
	"time"

	"github.com/mailjet/mailjet-apiv3-go/v3"

	_ "github.com/mattn/go-sqlite3"
)

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

// add cli help flag that describes exit codes (!). Does it necessary?
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

	log.Println("Scraper started!")

	ticker := time.NewTicker(time.Hour * 6)
	defer ticker.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("program execution interrupted, exiting")
				done <- struct{}{}
				return
			case <-ticker.C:
				process(scraperService, imdbService, mailerService)
			}
		}
	}()

	<-done

	fmt.Println("Done!")
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
