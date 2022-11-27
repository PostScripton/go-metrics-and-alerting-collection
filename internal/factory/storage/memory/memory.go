package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
)

type MemoryStorage struct {
	mu      sync.Mutex
	metrics map[string]metrics.Metrics
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		mu:      sync.Mutex{},
		metrics: make(map[string]metrics.Metrics),
	}
}

func (ms *MemoryStorage) GetCollection() (map[string]metrics.Metrics, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.metrics, nil
}

func (ms *MemoryStorage) StoreCollection(collection map[string]metrics.Metrics) error {
	for _, metric := range collection {
		if err := ms.Store(metric); err != nil {
			return err
		}
	}

	return nil
}

func (ms *MemoryStorage) Get(metric metrics.Metrics) (*metrics.Metrics, error) {
	if valid, err := metric.Validate(); !valid {
		return nil, err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	if value, ok := ms.metrics[metric.ID]; ok {
		return &value, nil
	}

	return nil, metrics.ErrNoValue
}

func (ms *MemoryStorage) Store(metric metrics.Metrics) error {
	if valid, err := metric.Validate(); !valid {
		return err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	storedMetric, ok := ms.metrics[metric.ID]
	if !ok {
		ms.metrics[metric.ID] = *metrics.New(metric.Type, metric.ID)
		storedMetric = ms.metrics[metric.ID]
	}

	metrics.Update(&storedMetric, &metric)

	ms.metrics[metric.ID] = storedMetric

	return nil
}

func (ms *MemoryStorage) CleanUp() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.metrics = make(map[string]metrics.Metrics)

	return nil
}

func (ms *MemoryStorage) Ping(_ context.Context) error {
	if ms.metrics == nil {
		return fmt.Errorf("metrics map collection is not initialized")
	}
	return nil
}

func (ms *MemoryStorage) Close() {
}
