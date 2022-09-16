package storage

import (
	"context"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
)

type Storager interface {
	CollectionGetter
	CollectionStorer
	Getter
	Storer
	CleanUper
	Pinger
	Closer
}

type CollectionGetter interface {
	GetCollection() (map[string]metrics.Metrics, error)
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

type CleanUper interface {
	CleanUp() error
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Closer interface {
	Close()
}
