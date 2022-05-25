package metrics

import "github.com/prometheus/client_golang/prometheus"

var AllMetrics = []prometheus.Collector{
	ScraperHeartbeat,
	ScraperLatency,
	ScraperErrors,
	ScraperCache,
	ScraperCacheEvent,
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
