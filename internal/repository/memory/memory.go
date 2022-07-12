package memory

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"sync"
)

type memoryStorage struct {
	mu             sync.Mutex
	counterMetrics map[string]metrics.Counter
	gaugeMetrics   map[string]metrics.Gauge
}

func New() *memoryStorage {
	return &memoryStorage{
		mu:             sync.Mutex{},
		counterMetrics: make(map[string]metrics.Counter),
		gaugeMetrics:   make(map[string]metrics.Gauge),
	}
}

func (storage *memoryStorage) GetCounterMetrics() map[string]metrics.Counter {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	return storage.counterMetrics
}

func (storage *memoryStorage) GetGaugeMetrics() map[string]metrics.Gauge {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	return storage.gaugeMetrics
}

func (storage *memoryStorage) StoreCounter(name string, value metrics.Counter) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.counterMetrics[name] += value
}

func (storage *memoryStorage) StoreGauge(name string, value metrics.Gauge) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.gaugeMetrics[name] = value
}
