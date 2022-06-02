package main

import (
	"common/client"
	"common/proto"
	"database/sql"
	"log"
	"scraper/config"
	"scraper/metrics"
	"scraper/service"
	"scraper/service/decorator"

	"google.golang.org/grpc"

	"scraper/storage/repository"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/mailjet/mailjet-apiv3-go/v3"
)

func initServices(db *sql.DB, grpcConn grpc.ClientConnInterface) {

	filmRepo := repository.NewFilmStorageRepository(db)

	httpClient, err := client.NewHttpClientWithCookies()
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_create_http_client_with_cookies").Inc()
		log.Fatalln(err)
	}

	soupService := decorator.NewSoupDecorator(httpClient)
	scraper = service.NewScraper(filmRepo, soupService)
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_create_scraper").Inc()
		log.Fatalln(err)
	}

	imdb = service.NewIMDB(proto.NewIMDBClient(grpcConn))

	emailRepo := repository.NewSubscriberRepository(db)

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

func initAndMaintainDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:./pirata.db")
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_establish_db_connection").Inc()
		return nil, err
	}

	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("fs_error_migration_files").Inc()
		return nil, err
	}

	migrationDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_init_migration_driver").Inc()
		return nil, err
	}

	// NewWithInstance always returns nil error
	migration, err := migrate.NewWithInstance("migrations", sourceDriver, "sqlite3", migrationDriver)
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_init_migration").Inc()
		return nil, err
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		metrics.ScraperErrors.WithLabelValues("migration_failed").Inc()
		return nil, err
	}

	return db, nil
}
