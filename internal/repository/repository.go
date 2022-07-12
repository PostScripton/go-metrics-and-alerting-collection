package repository

import "github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"

type Getter interface {
	GetCounterMetrics() map[string]metrics.Counter
	GetGaugeMetrics() map[string]metrics.Gauge
}

type Storer interface {
	StoreCounter(name string, value metrics.Counter)
	StoreGauge(name string, value metrics.Gauge)
}

type Storage interface {
	Getter
	Storer
}
