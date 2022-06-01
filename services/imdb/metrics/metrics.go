package metrics

import "github.com/prometheus/client_golang/prometheus"

var HitMissCache = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "hit_miss_cache",
		Help: "Cache hit and miss gauge",
	})

var CacheEvent = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "cache_events",
		Help: "Cache gauge events",
	}, []string{"event"})
