package service

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"os"
	"strconv"
)

type MetricService struct {
	storage repository.Storer
}

func NewMetricService(storage repository.Storage) *MetricService {
	return &MetricService{
		storage: storage,
	}
}

func (s *MetricService) UpdateMetric(metricType string, metricName string, metricValue string) {
	fmt.Printf("[%s] \"%s\": %s\n", metricType, metricName, metricValue)
	switch metricType {
	case metrics.StringCounterType:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		s.storage.StoreCounter(metricName, metrics.Counter(value))
	case metrics.StringGaugeType:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		s.storage.StoreGauge(metricName, metrics.Gauge(value))
	}
}
