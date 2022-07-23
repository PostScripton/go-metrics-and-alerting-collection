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

func (s *memoryStorage) GetCounterMetrics() map[string]metrics.Counter {
	return s.counterMetrics
}

func (s *memoryStorage) GetGaugeMetrics() map[string]metrics.Gauge {
	return s.gaugeMetrics
}

func (s *memoryStorage) Get(t string, name string) (metrics.MetricType, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch t {
	case metrics.StringCounterType:
		if value, ok := s.counterMetrics[name]; ok {
			return value, nil
		}
	case metrics.StringGaugeType:
		if value, ok := s.gaugeMetrics[name]; ok {
			return value, nil
		}
	default:
		return nil, metrics.ErrUnsupportedType
	}

	return nil, metrics.ErrNoValue
}

func (s *memoryStorage) Store(name string, value metrics.MetricType) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch v := value.(type) {
	case metrics.Counter:
		s.counterMetrics[name] += v
	case metrics.Gauge:
		s.gaugeMetrics[name] = v
	}
}
