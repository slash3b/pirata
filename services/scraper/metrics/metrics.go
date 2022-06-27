package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func Start(logger logrus.FieldLogger) {
	prometheus.MustRegister(
		ScraperHeartbeat,
		ScraperLatency,
		ScraperErrors,
		ScraperCache,
		ScraperCacheEvent,
	)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Println("metrics available on :2112/metrics")
		err := http.ListenAndServe(":2112", nil)
		if err != nil {
			ScraperErrors.WithLabelValues("unable_to_start_metrics").Inc()
			logger.Println(fmt.Errorf("unable to start metrics %v", err))
		}
	}()

}

var ScraperHeartbeat = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "scraper_heartbeat",
		Help: "Scraper heartbeat",
	})

var ScraperLatency = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name: "scraper_latency",
		Help: "Scraper latency",
	})

var ScraperErrors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "scraper_errors",
		Help: "Scraper Errors",
	}, []string{"err_type"})

var ScraperCache = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "scraper_cache",
		Help: "Scraper cache hit and miss gauge",
	})

var ScraperCacheEvent = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "scraper_cache_events",
		Help: "Scraper events",
	}, []string{"event"})
