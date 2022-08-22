package repository

import "github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"

type Storager interface {
	CollectionGetter
	CollectionStorer
	Getter
	Storer
}

type CollectionGetter interface {
	GetCollection() map[string]metrics.Metrics
}
type CollectionStorer interface {
	StoreCollection(map[string]metrics.Metrics) error
}

type Getter interface {
	Get(metric metrics.Metrics) (*metrics.Metrics, error)
}

type Storer interface {
	Store(metric metrics.Metrics) error
}
