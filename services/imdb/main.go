package main

import (
	"common/client"
	"common/proto"
	"fmt"
	"imdb/metrics"
	"imdb/service"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	metricsStart()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcSrv := grpc.NewServer()

	//c := cache.NewLRU[string, dto.IMDBData](20)

	cln, err := client.NewHttpClientWithCookies()
	if err != nil {
		panic(err)
	}
	imdbService := service.NewIMDB(cln)

	proto.RegisterIMDBServer(grpcSrv, service.NewIMDBGrpc(imdbService))

	fmt.Println("started GRPC server on port 50051")
	err = grpcSrv.Serve(listener)
	if err != nil {
		panic(err)
	}

}

func metricsStart() {

	prometheus.MustRegister(
		metrics.HitMissCache,
		metrics.CacheEvent,
	)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("metrics available on :2112/metrics")
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			log.Println(fmt.Errorf("unable to start metrics %v", err))
		}
	}()
}
