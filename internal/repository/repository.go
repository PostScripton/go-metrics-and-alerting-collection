package repository

import "github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"

type Storager interface {
	Getter
	Storer
}

type CollectionGetter interface {
	GetCounterMetrics() map[string]metrics.Counter
	GetGaugeMetrics() map[string]metrics.Gauge
}

type Getter interface {
	Get(t string, name string) (metrics.MetricType, error)
}

type Storer interface {
	Store(name string, value metrics.MetricType)
}
