package main

import (
	"context"
	"embed"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"scraper/metrics"
	"scraper/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*
var migrationFiles embed.FS

var (
	scraper *service.Scraper
	imdb    *service.IMDB
	mailer  *service.Mailer
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logger := initLogger("scraper")

	db, err := initAndMaintainDB()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			metrics.ScraperErrors.WithLabelValues("could_not_close_db_connection").Inc()
			logger.Errorf("could not close database connection %v", err)
		}
	}()

	conn, err := grpc.Dial("imdb:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("cound not connect to imdb service, %v", err)
		os.Exit(1)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			logger.Errorf("could not close grpc connection properly %v", err)
			metrics.ScraperErrors.WithLabelValues("could_not_close_grpc_connection").Inc()
		}
	}()

	metrics.Start(logger)

	initServices(db, conn, logger)

	ticker := time.NewTicker(time.Minute * 60) // todo: make the program reload when env variable changes!
	defer ticker.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	logger.Info("Scraper started!")

	process(scraper, imdb, mailer, logger)

	metrics.ScraperHeartbeat.Inc()
	// syscall.SIGHUP
	// possibly make ticker change on like the do in caddy and prometheus

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("program execution interrupted, exiting")
				done <- struct{}{}
				return
			case <-ticker.C:
				metrics.ScraperHeartbeat.Inc()
				process(scraper, imdb, mailer, logger)
				logger.Info("done processing")
			}
		}
	}()

	<-done

	fmt.Println("Scraper service stopped!")
}

// process works using kind of "pipeline" pattern, or smth close
func process(scraper *service.Scraper, imdb *service.IMDB, mailer *service.Mailer, logger logrus.FieldLogger) {
	timer := prometheus.NewTimer(metrics.ScraperLatency)
	defer timer.ObserveDuration()

	ctx, release := context.WithTimeout(context.Background(), time.Second*30)
	defer release()

	err := mailer.Send(imdb.GetFilms(scraper.GetFilms(ctx)))
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("scraper_error").Inc()
		logger.Println(err)
	}
}
