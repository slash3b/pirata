package main

import (
	"common/client"
	"common/proto"
	"context"
	"fmt"
	"imdb/metrics"
	"imdb/service"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	errCh := make(chan error, 1)
	go metricsStart(errCh)
	go grcpStart(errCh)

	select {
	case <-ctx.Done():
		fmt.Println("exiting due to system signal received ...")
	case e := <-errCh:
		fmt.Println("error received, exiting, ", e)
	}
}

func grcpStart(errCh chan<- error) {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcSrv := grpc.NewServer()

	cln, err := client.NewHttpClientWithCookies()
	if err != nil {
		errCh <- err
		return
	}
	imdbService := service.NewIMDB(cln)

	proto.RegisterIMDBServer(grpcSrv, service.NewIMDBGrpc(imdbService))

	log.Println("started GRPC server on port 50051")
	err = grpcSrv.Serve(listener)
	if err != nil {
		errCh <- err
		return
	}
}

func metricsStart(errCh chan<- error) {
	prometheus.MustRegister(
		metrics.HitMissCache,
		metrics.CacheEvent,
	)

	http.Handle("/metrics", promhttp.Handler())
	log.Println("metrics available on :2112/metrics")

	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Println(fmt.Errorf("unable to start metrics %v", err))
		errCh <- err
	}
}
