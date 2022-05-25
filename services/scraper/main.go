package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"scraper/client"
	"scraper/config"
	"scraper/metrics"
	"scraper/service"
	"scraper/service/cineplex"
	"scraper/service/cineplex/decorator"
	"scraper/storage/repository"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"

	"github.com/mailjet/mailjet-apiv3-go/v3"

	_ "github.com/mattn/go-sqlite3"

	_ "embed"
)

//go:embed migrations/*
var migrationFiles embed.FS

var (
	scraper *cineplex.Scraper
	imdb    *service.IMDB
	mailer  *service.Mailer
)

func main() {
	metrics.Start()

	db, err := sql.Open("sqlite3", "file:./pirata.db")
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_establish_db_connection").Inc()
		log.Fatalln(err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			metrics.ScraperErrors.WithLabelValues("could_not_close_db_connection").Inc()
			log.Fatalln(err)
		}
	}()

	err = dbMigrationUp(db)
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("migration_failed").Inc()
		log.Fatalln(err)
	}

	initServices(db)

	ticker := time.NewTicker(time.Minute * 60)
	defer ticker.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	log.Println("Scraper started!")
	err = process(scraper, imdb, mailer) // this process should be another service so I can do myService.Scrape() // or not ?
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("could_not_run_scraper_process").Inc()
		log.Println(err)
	}

	metrics.ScraperHeartbeat.Inc()
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
				metrics.ScraperHeartbeat.Inc()
				err = process(scraper, imdb, mailer)
				if err != nil {
					metrics.ScraperErrors.WithLabelValues("could_not_run_scraper_process").Inc()
					log.Println(err)
				}
			}
		}
	}()

	<-done

	fmt.Println("Done!")
}

func initServices(db *sql.DB) {

	filmRepo := repository.NewFilmStorageRepository(db)
	emailRepo := repository.NewSubscriberRepository(db)

	httpClient, err := client.NewHttpClientWithCookies()
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_create_http_client_with_cookies").Inc()
		log.Fatalln(err)
	}

	soupService := decorator.NewSoupDecorator(httpClient)
	scraper = cineplex.NewScraper(filmRepo, soupService)
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_create_scraper").Inc()
		log.Fatalln(err)
	}

	imdb = service.NewIMDB(httpClient)

	env, err := config.GetEnv()
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("incomplete_environment").Inc()
		log.Fatalln(err)
	}

	mailjetClient := mailjet.NewMailjetClient(env.MailjetPubKey, env.MailJetPrivateKey)

	mailer = service.NewMailer(mailjetClient, service.MailerConfig{
		FromEmail: env.FromEmail,
		FromName:  env.FromName,
	}, emailRepo)
}

// pipeline pattern may be ?
func process(scraper *cineplex.Scraper, imdb *service.IMDB, mailer service.Sender) error {
	timer := prometheus.NewTimer(metrics.ScraperLatency)

	defer timer.ObserveDuration()

	newFilms, err := scraper.GetAllFilms()
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("scraper_could_not_get_films").Inc()
		return err
	}

	if len(newFilms) > 0 {
		emailFilms := imdb.FindFilms(newFilms)

		err = mailer.Send(emailFilms)
		if err != nil {
			metrics.ScraperErrors.WithLabelValues("could_not_send_email").Inc()
			return err
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
