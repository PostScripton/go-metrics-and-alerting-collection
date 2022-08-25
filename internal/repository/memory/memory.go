package memory

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"sync"
)

type memoryStorage struct {
	mu      sync.Mutex
	metrics map[string]metrics.Metrics
}

var _ repository.Storager = &memoryStorage{}

func NewMemoryStorage() repository.Storager {
	return &memoryStorage{
		mu:      sync.Mutex{},
		metrics: make(map[string]metrics.Metrics),
	}
}

func (s *memoryStorage) GetCollection() (map[string]metrics.Metrics, error) {
	return s.metrics, nil
}

func (s *memoryStorage) StoreCollection(metrics map[string]metrics.Metrics) error {
	s.mu.Lock()
	s.metrics = metrics
	s.mu.Unlock()

	return nil
}

func (s *memoryStorage) Get(metric metrics.Metrics) (*metrics.Metrics, error) {
	if valid, err := metric.Validate(); !valid {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if value, ok := s.metrics[metric.ID]; ok {
		return &value, nil
	}

	return nil, metrics.ErrNoValue
}

func (s *memoryStorage) Store(metric metrics.Metrics) error {
	if valid, err := metric.Validate(); !valid {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	storedMetric, ok := s.metrics[metric.ID]
	if !ok {
		s.metrics[metric.ID] = *metrics.New(metric.Type, metric.ID)
		storedMetric = s.metrics[metric.ID]
	}

	switch metric.Type {
	case metrics.StringCounterType:
		var delta int64
		if storedMetric.Delta != nil {
			delta = *storedMetric.Delta
		}
		delta += *metric.Delta
		storedMetric.Delta = &delta
	case metrics.StringGaugeType:
		storedMetric.Value = metric.Value
	}

	s.metrics[metric.ID] = storedMetric

	return nil
}
