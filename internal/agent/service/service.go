package service

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"math/rand"
	"runtime"
)

type Storage interface {
	repository.CollectionGetter
	repository.Storer
}

type MetricService struct {
	storage  Storage
	client   client.Client
	memStats runtime.MemStats
}

func NewMetricService(storage Storage, client client.Client) *MetricService {
	return &MetricService{
		storage:  storage,
		client:   client,
		memStats: runtime.MemStats{},
	}
}

func (s *MetricService) GatherMetrics() {
	fmt.Println("Gathering...")
	s.storage.Store("PollCount", metrics.Counter(1))

	s.storage.Store("Alloc", metrics.Gauge(s.memStats.Alloc))
	s.storage.Store("BuckHashSys", metrics.Gauge(s.memStats.BuckHashSys))
	s.storage.Store("Frees", metrics.Gauge(s.memStats.Frees))
	s.storage.Store("GCCPUFraction", metrics.Gauge(s.memStats.GCCPUFraction))
	s.storage.Store("GCSys", metrics.Gauge(s.memStats.GCSys))
	s.storage.Store("HeapAlloc", metrics.Gauge(s.memStats.HeapAlloc))
	s.storage.Store("HeapIdle", metrics.Gauge(s.memStats.HeapIdle))
	s.storage.Store("HeapInuse", metrics.Gauge(s.memStats.HeapInuse))
	s.storage.Store("HeapObjects", metrics.Gauge(s.memStats.HeapObjects))
	s.storage.Store("HeapReleased", metrics.Gauge(s.memStats.HeapReleased))
	s.storage.Store("HeapSys", metrics.Gauge(s.memStats.HeapSys))
	s.storage.Store("LastGC", metrics.Gauge(s.memStats.LastGC))
	s.storage.Store("MCacheInuse", metrics.Gauge(s.memStats.MCacheInuse))
	s.storage.Store("MCacheSys", metrics.Gauge(s.memStats.MCacheSys))
	s.storage.Store("MSpanInuse", metrics.Gauge(s.memStats.MSpanInuse))
	s.storage.Store("MSpanSys", metrics.Gauge(s.memStats.MSpanSys))
	s.storage.Store("Mallocs", metrics.Gauge(s.memStats.Mallocs))
	s.storage.Store("NextGC", metrics.Gauge(s.memStats.NextGC))
	s.storage.Store("NumForcedGC", metrics.Gauge(s.memStats.NumForcedGC))
	s.storage.Store("NumGC", metrics.Gauge(s.memStats.NumGC))
	s.storage.Store("OtherSys", metrics.Gauge(s.memStats.OtherSys))
	s.storage.Store("PauseTotalNs", metrics.Gauge(s.memStats.PauseTotalNs))
	s.storage.Store("StackInuse", metrics.Gauge(s.memStats.StackInuse))
	s.storage.Store("StackSys", metrics.Gauge(s.memStats.StackSys))
	s.storage.Store("Sys", metrics.Gauge(s.memStats.Sys))
	s.storage.Store("TotalAlloc", metrics.Gauge(s.memStats.TotalAlloc))
	s.storage.Store("RandomValue", metrics.Gauge(randomMetric(0, 10000)))
}

func (s *MetricService) SendMetrics() {
	fmt.Println("Reporting...")
	for counterName, counter := range s.storage.GetCounterMetrics() {
		s.client.UpdateMetric(counter.Type(), counterName, fmt.Sprintf("%v", counter))
	}
	for gaugeName, gauge := range s.storage.GetGaugeMetrics() {
		s.client.UpdateMetric(gauge.Type(), gaugeName, fmt.Sprintf("%v", gauge))
	}
}

func randomMetric(min float64, max float64) float64 {
	// [min, max)
	return min + rand.Float64()*(max-min)
}
