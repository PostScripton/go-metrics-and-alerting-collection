package service

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"math/rand"
	"runtime"
)

type MetricService struct {
	storage  repository.Storage
	client   client.Client
	memStats runtime.MemStats
}

func NewMetricService(storage repository.Storage, client client.Client) *MetricService {
	return &MetricService{
		storage:  storage,
		client:   client,
		memStats: runtime.MemStats{},
	}
}

func (s *MetricService) GatherMetrics() {
	fmt.Println("Gathering...")
	s.storage.StoreCounter("PollCount", 1)

	s.storage.StoreGauge("Alloc", metrics.Gauge(s.memStats.Alloc))
	s.storage.StoreGauge("BuckHashSys", metrics.Gauge(s.memStats.BuckHashSys))
	s.storage.StoreGauge("Frees", metrics.Gauge(s.memStats.Frees))
	s.storage.StoreGauge("GCCPUFraction", metrics.Gauge(s.memStats.GCCPUFraction))
	s.storage.StoreGauge("GCSys", metrics.Gauge(s.memStats.GCSys))
	s.storage.StoreGauge("HeapAlloc", metrics.Gauge(s.memStats.HeapAlloc))
	s.storage.StoreGauge("HeapIdle", metrics.Gauge(s.memStats.HeapIdle))
	s.storage.StoreGauge("HeapInuse", metrics.Gauge(s.memStats.HeapInuse))
	s.storage.StoreGauge("HeapObjects", metrics.Gauge(s.memStats.HeapObjects))
	s.storage.StoreGauge("HeapReleased", metrics.Gauge(s.memStats.HeapReleased))
	s.storage.StoreGauge("HeapSys", metrics.Gauge(s.memStats.HeapSys))
	s.storage.StoreGauge("LastGC", metrics.Gauge(s.memStats.LastGC))
	s.storage.StoreGauge("MCacheInuse", metrics.Gauge(s.memStats.MCacheInuse))
	s.storage.StoreGauge("MCacheSys", metrics.Gauge(s.memStats.MCacheSys))
	s.storage.StoreGauge("MSpanInuse", metrics.Gauge(s.memStats.MSpanInuse))
	s.storage.StoreGauge("MSpanSys", metrics.Gauge(s.memStats.MSpanSys))
	s.storage.StoreGauge("Mallocs", metrics.Gauge(s.memStats.Mallocs))
	s.storage.StoreGauge("NextGC", metrics.Gauge(s.memStats.NextGC))
	s.storage.StoreGauge("NumForcedGC", metrics.Gauge(s.memStats.NumForcedGC))
	s.storage.StoreGauge("NumGC", metrics.Gauge(s.memStats.NumGC))
	s.storage.StoreGauge("OtherSys", metrics.Gauge(s.memStats.OtherSys))
	s.storage.StoreGauge("PauseTotalNs", metrics.Gauge(s.memStats.PauseTotalNs))
	s.storage.StoreGauge("StackInuse", metrics.Gauge(s.memStats.StackInuse))
	s.storage.StoreGauge("StackSys", metrics.Gauge(s.memStats.StackSys))
	s.storage.StoreGauge("Sys", metrics.Gauge(s.memStats.Sys))
	s.storage.StoreGauge("TotalAlloc", metrics.Gauge(s.memStats.TotalAlloc))
	s.storage.StoreGauge("RandomValue", metrics.Gauge(randomMetric(0, 10000)))
}

func (s *MetricService) SendMetrics() {
	fmt.Println("Reporting...")
	for counterName, counter := range s.storage.GetCounterMetrics() {
		s.client.UpdateMetric(counter.Type(), counterName, fmt.Sprintf("%d", counter))
	}
	for gaugeName, gauge := range s.storage.GetGaugeMetrics() {
		s.client.UpdateMetric(gauge.Type(), gaugeName, fmt.Sprintf("%v", gauge))
	}
}

func randomMetric(min float64, max float64) float64 {
	// [min, max)
	return min + rand.Float64()*(max-min)
}
