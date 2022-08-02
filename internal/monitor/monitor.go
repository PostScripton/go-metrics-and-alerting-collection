package monitor

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

type Monitor struct {
	storage  Storager
	client   *client.Client
	memStats *runtime.MemStats
}

func New(storage Storager, client *client.Client) *Monitor {
	return &Monitor{
		storage:  storage,
		client:   client,
		memStats: &runtime.MemStats{},
	}
}

func (m *Monitor) Gather() {
	fmt.Println("Gathering...")
	runtime.ReadMemStats(m.memStats)

	m.storage.Store("PollCount", metrics.Counter(1))

	m.storage.Store("Alloc", metrics.Gauge(m.memStats.Alloc))
	m.storage.Store("BuckHashSys", metrics.Gauge(m.memStats.BuckHashSys))
	m.storage.Store("Frees", metrics.Gauge(m.memStats.Frees))
	m.storage.Store("GCCPUFraction", metrics.Gauge(m.memStats.GCCPUFraction))
	m.storage.Store("GCSys", metrics.Gauge(m.memStats.GCSys))
	m.storage.Store("HeapAlloc", metrics.Gauge(m.memStats.HeapAlloc))
	m.storage.Store("HeapIdle", metrics.Gauge(m.memStats.HeapIdle))
	m.storage.Store("HeapInuse", metrics.Gauge(m.memStats.HeapInuse))
	m.storage.Store("HeapObjects", metrics.Gauge(m.memStats.HeapObjects))
	m.storage.Store("HeapReleased", metrics.Gauge(m.memStats.HeapReleased))
	m.storage.Store("HeapSys", metrics.Gauge(m.memStats.HeapSys))
	m.storage.Store("LastGC", metrics.Gauge(m.memStats.LastGC))
	m.storage.Store("Lookups", metrics.Gauge(m.memStats.Lookups))
	m.storage.Store("MCacheInuse", metrics.Gauge(m.memStats.MCacheInuse))
	m.storage.Store("MCacheSys", metrics.Gauge(m.memStats.MCacheSys))
	m.storage.Store("MSpanInuse", metrics.Gauge(m.memStats.MSpanInuse))
	m.storage.Store("MSpanSys", metrics.Gauge(m.memStats.MSpanSys))
	m.storage.Store("Mallocs", metrics.Gauge(m.memStats.Mallocs))
	m.storage.Store("NextGC", metrics.Gauge(m.memStats.NextGC))
	m.storage.Store("NumForcedGC", metrics.Gauge(m.memStats.NumForcedGC))
	m.storage.Store("NumGC", metrics.Gauge(m.memStats.NumGC))
	m.storage.Store("OtherSys", metrics.Gauge(m.memStats.OtherSys))
	m.storage.Store("PauseTotalNs", metrics.Gauge(m.memStats.PauseTotalNs))
	m.storage.Store("StackInuse", metrics.Gauge(m.memStats.StackInuse))
	m.storage.Store("StackSys", metrics.Gauge(m.memStats.StackSys))
	m.storage.Store("Sys", metrics.Gauge(m.memStats.Sys))
	m.storage.Store("TotalAlloc", metrics.Gauge(m.memStats.TotalAlloc))

	random := rand.Float64() * (10000)
	m.storage.Store("RandomValue", metrics.Gauge(random))
}

func (m *Monitor) Send() {
	fmt.Println("Reporting...")
	for counterName, counter := range m.storage.GetCounterMetrics() {
		m.client.UpdateMetricJSON(counter.Type(), counterName, counter)
	}
	for gaugeName, gauge := range m.storage.GetGaugeMetrics() {
		m.client.UpdateMetricJSON(gauge.Type(), gaugeName, gauge)
	}
}
