package main

import "github.com/prometheus/client_golang/prometheus"

var metrics = []prometheus.Collector{
	ScraperHeartbeatMetric,
	ScraperLatencyMetric,
	ScraperErrorsMetric,
}

var f = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "scraper_heartbeat",
	Help: "Scraper heartbeat",
}, []string{})

var ScraperHeartbeatMetric = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "scraper_heartbeat",
		Help: "Scraper heartbeat",
	})

var ScraperLatencyMetric = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name: "scraper_latency",
		Help: "Scraper latency",
	})

var ScraperErrorsMetric = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "scraper_errors",
		Help: "Scraper Errors",
	}, []string{"err_type"})
