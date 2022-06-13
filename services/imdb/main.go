package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"common/client"
	"common/proto"
	"imdb/metrics"
	"imdb/service"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {

	viper.SetDefault("GRPC_PORT", "50052")

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
	grpcPort := viper.GetString("GRPC_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		panic(err)
	}

	grpcSrv := grpc.NewServer()

	cln, err := client.NewHttpClientWithCookies()
	if err != nil {
		errCh <- err
		return
	}

	/*
		Question: how the heck can I turn cache in some sort of an interface ?
	*/
	imdbService := service.NewIMDB(cln)

	proto.RegisterIMDBServer(grpcSrv, service.NewGRPCServer(imdbService))

	log.Printf("started GRPC server on port %s", grpcPort)
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
