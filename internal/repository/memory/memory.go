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

func (storage *memoryStorage) Get(t string, name string) (metrics.MetricType, error) {
	switch t {
	case metrics.StringCounterType:
		if value, ok := storage.counterMetrics[name]; ok {
			return value, nil
		}
	case metrics.StringGaugeType:
		if value, ok := storage.gaugeMetrics[name]; ok {
			return value, nil
		}
	default:
		return nil, metrics.ErrUnsupportedType
	}

	return nil, metrics.ErrNoValue
}

func (storage *memoryStorage) Store(name string, value metrics.MetricType) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	switch v := value.(type) {
	case metrics.Counter:
		storage.counterMetrics[name] += v
	case metrics.Gauge:
		storage.gaugeMetrics[name] = v
	}
}
