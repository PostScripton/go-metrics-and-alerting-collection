package repository

import "github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"

type CollectionGetter interface {
	GetCounterMetrics() map[string]metrics.Counter
	GetGaugeMetrics() map[string]metrics.Gauge
}

type Storer interface {
	Store(name string, value metrics.MetricType)
}
