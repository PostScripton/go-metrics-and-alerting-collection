package service

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"os"
	"strconv"
)

type storage interface {
	repository.Storer
}

type MetricService struct {
	storage storage
}

func NewMetricService(storage storage) *MetricService {
	return &MetricService{
		storage: storage,
	}
}

func (s *MetricService) UpdateMetric(t string, name string, value string) {
	switch t {
	case metrics.StringCounterType:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		s.storage.Store(name, metrics.Counter(v))
	case metrics.StringGaugeType:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		s.storage.Store(name, metrics.Gauge(v))
	}
}
