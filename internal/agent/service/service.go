package service

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"math/rand"
	"runtime"
)

type Storager interface {
	repository.CollectionGetter
	repository.Storer
}

type MetricService struct {
	storager Storager
	client   *client.Client
	memStats *runtime.MemStats
}

func NewMetricService(storager Storager, client *client.Client) *MetricService {
	return &MetricService{
		storager: storager,
		client:   client,
		memStats: &runtime.MemStats{},
	}
}

func (s *MetricService) GatherMetrics() {
	fmt.Println("Gathering...")
	s.storager.Store("PollCount", metrics.Counter(1))

	s.storager.Store("Alloc", metrics.Gauge(s.memStats.Alloc))
	s.storager.Store("BuckHashSys", metrics.Gauge(s.memStats.BuckHashSys))
	s.storager.Store("Frees", metrics.Gauge(s.memStats.Frees))
	s.storager.Store("GCCPUFraction", metrics.Gauge(s.memStats.GCCPUFraction))
	s.storager.Store("GCSys", metrics.Gauge(s.memStats.GCSys))
	s.storager.Store("HeapAlloc", metrics.Gauge(s.memStats.HeapAlloc))
	s.storager.Store("HeapIdle", metrics.Gauge(s.memStats.HeapIdle))
	s.storager.Store("HeapInuse", metrics.Gauge(s.memStats.HeapInuse))
	s.storager.Store("HeapObjects", metrics.Gauge(s.memStats.HeapObjects))
	s.storager.Store("HeapReleased", metrics.Gauge(s.memStats.HeapReleased))
	s.storager.Store("HeapSys", metrics.Gauge(s.memStats.HeapSys))
	s.storager.Store("LastGC", metrics.Gauge(s.memStats.LastGC))
	s.storager.Store("MCacheInuse", metrics.Gauge(s.memStats.MCacheInuse))
	s.storager.Store("MCacheSys", metrics.Gauge(s.memStats.MCacheSys))
	s.storager.Store("MSpanInuse", metrics.Gauge(s.memStats.MSpanInuse))
	s.storager.Store("MSpanSys", metrics.Gauge(s.memStats.MSpanSys))
	s.storager.Store("Mallocs", metrics.Gauge(s.memStats.Mallocs))
	s.storager.Store("NextGC", metrics.Gauge(s.memStats.NextGC))
	s.storager.Store("NumForcedGC", metrics.Gauge(s.memStats.NumForcedGC))
	s.storager.Store("NumGC", metrics.Gauge(s.memStats.NumGC))
	s.storager.Store("OtherSys", metrics.Gauge(s.memStats.OtherSys))
	s.storager.Store("PauseTotalNs", metrics.Gauge(s.memStats.PauseTotalNs))
	s.storager.Store("StackInuse", metrics.Gauge(s.memStats.StackInuse))
	s.storager.Store("StackSys", metrics.Gauge(s.memStats.StackSys))
	s.storager.Store("Sys", metrics.Gauge(s.memStats.Sys))
	s.storager.Store("TotalAlloc", metrics.Gauge(s.memStats.TotalAlloc))
	s.storager.Store("RandomValue", metrics.Gauge(randomMetric(0, 10000)))
}

func (s *MetricService) SendMetrics() {
	fmt.Println("Reporting...")
	for counterName, counter := range s.storager.GetCounterMetrics() {
		s.client.UpdateMetric(counter.Type(), counterName, fmt.Sprintf("%v", counter))
	}
	for gaugeName, gauge := range s.storager.GetGaugeMetrics() {
		s.client.UpdateMetric(gauge.Type(), gaugeName, fmt.Sprintf("%v", gauge))
	}
}

func randomMetric(min float64, max float64) float64 {
	// [min, max)
	return min + rand.Float64()*(max-min)
}
