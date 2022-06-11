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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"

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

	conn, err := grpc.Dial("imdb:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v \n", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			metrics.ScraperErrors.WithLabelValues("could_not_close_grpc_connection").Inc()
			log.Println(err)
		}
	}()

	metrics.Start()

	initServices(db, conn)

	ticker := time.NewTicker(time.Minute * 60) // todo: make the program reload when env variable changes!
	defer ticker.Stop()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	log.Println("Scraper started!")

	process(scraper, imdb, mailer)

	metrics.ScraperHeartbeat.Inc()
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
				process(scraper, imdb, mailer)
			}
		}
	}()

	<-done

	fmt.Println("Scraper service stopped!")
}

// process works using kind of "pipeline" pattern, or smth close
func process(scraper *service.Scraper, imdb *service.IMDB, mailer *service.Mailer) {
	timer := prometheus.NewTimer(metrics.ScraperLatency)
	defer timer.ObserveDuration()

	ctx, release := context.WithTimeout(context.Background(), time.Second*30)
	defer release()

	err := mailer.Send(imdb.GetFilms(scraper.GetFilms(ctx)))
	if err != nil {
		metrics.ScraperErrors.WithLabelValues("scraper_error").Inc()
		log.Println(err)
	}
}
