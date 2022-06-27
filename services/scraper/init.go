package main

import (
	"common/client"
	"common/proto"
	"database/sql"
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
	"github.com/sirupsen/logrus"
)

func initServices(db *sql.DB, grpcConn grpc.ClientConnInterface, logger logrus.FieldLogger) {

	filmRepo := repository.NewFilmStorageRepository(db)

	httpClient, err := client.NewHttpClientWithCookies()
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_create_http_client_with_cookies").Inc()
		logger.Fatalln(err)
	}

	soupService := decorator.NewSoupDecorator(httpClient)
	scraper = service.NewScraper(filmRepo, soupService, logger)
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("unable_to_create_scraper").Inc()
		logger.Fatalln(err)
	}

	imdb = service.NewIMDB(proto.NewIMDBClient(grpcConn), logger)

	emailRepo := repository.NewSubscriberRepository(logger, db)

	env, err := config.GetEnv()
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("incomplete_environment").Inc()
		logger.Fatalln(err)
	}

	mailjetClient := mailjet.NewMailjetClient(env.MailjetPubKey, env.MailJetPrivateKey)

	mailer = service.NewMailer(mailjetClient, service.MailerConfig{
		FromEmail: env.FromEmail,
		FromName:  env.FromName,
	}, emailRepo, logger)
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

func initLogger(serviceName string) *logrus.Entry {
	log := logrus.New()

	return log.WithFields(logrus.Fields{"service_name": serviceName})
}
