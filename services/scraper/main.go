package main

import (
	"context"
	"database/sql"
	"embed"
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

	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"

	"github.com/mailjet/mailjet-apiv3-go/v3"

	_ "github.com/mattn/go-sqlite3"

	_ "embed"
)

type runtimeConfig struct {
	mailjetPubKey     string
	mailJetPrivateKey string
	sqlitePath        string
}

var runtimeConf runtimeConfig

//go:embed migrations/*
var migrationFiles embed.FS

func init() {
	runtimeConf = runtimeConfig{
		mailjetPubKey:     os.Getenv("MJ_APIKEY_PUBLIC"),
		mailJetPrivateKey: os.Getenv("MJ_APIKEY_PRIVATE"),
		sqlitePath:        os.Getenv("SQLITE_PATH"),
	}
}

// add cli help flag that describes exit codes (!). Does it necessary?
func main() {

	db, err := sql.Open("sqlite3", "file:./pirata.db")
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		panic(err)
	}

	migrationDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		panic(err)
	}

	migration, err := migrate.NewWithInstance("migrations", sourceDriver, "sqlite3", migrationDriver)

	err = migration.Up()
	if err != nil {
		panic(err)
	}

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
	process(scraperService, imdbService, mailerService)

	ticker := time.NewTicker(time.Hour * 3)
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
