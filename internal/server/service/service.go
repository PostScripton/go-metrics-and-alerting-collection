package service

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"strconv"
)

type MetricService struct {
	storer repository.Storer
}

func NewMetricService(storer repository.Storer) *MetricService {
	return &MetricService{
		storer: storer,
	}
}

func (s *MetricService) UpdateMetric(t string, name string, value string) {
	switch t {
	case metrics.StringCounterType:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(err)
		}
		s.storer.Store(name, metrics.Counter(v))
	case metrics.StringGaugeType:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(err)
		}
		s.storer.Store(name, metrics.Gauge(v))
	}
}
