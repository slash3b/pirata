package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scraper/client"
	"scraper/service"
	"scraper/service/cineplex"
	"scraper/service/cineplex/decorator"
	"scraper/storage/repository"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

var scrapeCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "pirata_scrape_counter",
		Help: "Counter for scrape times, just to see it running",
	})

var scrapeHist = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name: "pirata_scrape_hist",
		Help: "Scrape times",
	})

// add cli help flag that describes exit codes (!). Does it necessary?
func main() {

	prometheus.MustRegister(scrapeCounter, scrapeHist)

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

	err = dbMigrationUp(db)
	if err != nil {
		log.Fatalln(err)
	}

	filmRepo := repository.NewFilmStorageRepository(db)
	emailRepo := repository.NewSubscriberRepository(db)

	httpClient, err := client.NewHttpClientWithCookies()
	if err != nil {
		log.Fatalln(err)
	}

	soupService := decorator.NewSoupDecorator(httpClient)
	scraperService := cineplex.NewScraper(filmRepo, soupService)
	if err != nil {
		log.Fatalln(err)
	}

	imdbService := service.NewIMDB(httpClient)

	mailjetClient := mailjet.NewMailjetClient(runtimeConf.mailjetPubKey, runtimeConf.mailJetPrivateKey)

	mailerService := service.NewMailer(mailjetClient, service.MailerConfig{
		FromEmail: os.Getenv("FROM_EMAIL"),
		FromName:  os.Getenv("FROM_NAME"),
	}, emailRepo)

	ticker := time.NewTicker(time.Minute * 3)
	defer ticker.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	log.Println("Scraper started!")
	err = process(scraperService, imdbService, mailerService) // this process should be another service so I can do myService.Scrape() // or not ?
	if err != nil {
		log.Println(err)
	}
	scrapeCounter.Inc()
	// learn about proper context propagation
	// syscall.SIGHUP
	// possibly make ticker change on like the do in caddy and prometheus

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("program execution interrupted, exiting")
				done <- struct{}{}
				return
			case <-ticker.C:
				err = process(scraperService, imdbService, mailerService)
				if err != nil {
					log.Println(err)
				}
				scrapeCounter.Inc()
			}
		}
	}()

	<-done

	fmt.Println("Done!")
}

// pipeline pattern may be ?
func process(scraper *cineplex.Scraper, imdb *service.IMDB, mailer service.Sender) error {
	timer := prometheus.NewTimer(scrapeHist)

	defer timer.ObserveDuration()

	newFilms, err := scraper.GetAllFilms()
	if err != nil {
		return err
	}

	if len(newFilms) > 0 {
		emailFilms := imdb.EnrichFilms(newFilms)

		err = mailer.Send(emailFilms)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func dbMigrationUp(db *sql.DB) error {
	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	migrationDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	// NewWithInstance always returns nil error
	migration, err := migrate.NewWithInstance("migrations", sourceDriver, "sqlite3", migrationDriver)
	if err != nil {
		return err
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
