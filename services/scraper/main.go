package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"scraper/metrics"
	"scraper/service"

	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	_ "github.com/mattn/go-sqlite3"

	_ "embed"
)

//go:embed migrations/*
var migrationFiles embed.FS

var (
	scraper *service.Scraper
	imdb    *service.IMDB
	mailer  *service.Mailer
)

func main() {
	db, err := initAndMaintainDB()
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			metrics.ScraperErrors.WithLabelValues("could_not_close_db_connection").Inc()
			log.Println(err)
		}
	}()

	metrics.Start()

	initServices(db)

	ticker := time.NewTicker(time.Minute * 60)
	defer ticker.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	log.Println("Scraper started!")

	err = process(scraper, imdb, mailer)
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("scraper_error").Inc()
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
					metrics.ScraperErrors.WithLabelValues("scraper_error").Inc()
					log.Println(err)
				}
			}
		}
	}()

	<-done

	fmt.Println("Done!")
}

// process works using kind of "pipeline" pattern, or smth close
func process(scraper *service.Scraper, imdb *service.IMDB, mailer service.Sender) error {
	timer := prometheus.NewTimer(metrics.ScraperLatency)
	defer timer.ObserveDuration()

	ctx, release := context.WithTimeout(context.Background(), time.Second*30)
	defer release()

	return mailer.Send(imdb.FindFilms(scraper.GetFilms(ctx)))
}
