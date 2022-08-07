package repository

import "github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"

type CollectionGetter interface {
	GetMetrics() map[string]metrics.Metrics
}

type Getter interface {
	Get(metric metrics.Metrics) (*metrics.Metrics, error)
}

type Storer interface {
	Store(metric metrics.Metrics) error
}
