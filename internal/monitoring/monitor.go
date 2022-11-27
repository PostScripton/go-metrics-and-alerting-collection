package monitoring

import (
	"math/rand"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
)

type Monitorer interface {
	GatherMain()
	GatherAdditional()
	Send()
}

type Monitor struct {
	storage  storage.Storager
	client   *client.Client
	memStats *runtime.MemStats
}

func NewMonitor(storage storage.Storager, client *client.Client) Monitorer {
	return &Monitor{
		storage:  storage,
		client:   client,
		memStats: &runtime.MemStats{},
	}
}

func (m *Monitor) GatherMain() {
	log.Info().Msg("Gathering main...")
	runtime.ReadMemStats(m.memStats)

	_ = m.storage.Store(*metrics.NewCounter("PollCount", 1))

	_ = m.storage.Store(*metrics.NewGauge("Alloc", float64(m.memStats.Alloc)))
	_ = m.storage.Store(*metrics.NewGauge("BuckHashSys", float64(m.memStats.BuckHashSys)))
	_ = m.storage.Store(*metrics.NewGauge("Frees", float64(m.memStats.Frees)))
	_ = m.storage.Store(*metrics.NewGauge("GCCPUFraction", float64(m.memStats.GCCPUFraction)))
	_ = m.storage.Store(*metrics.NewGauge("GCSys", float64(m.memStats.GCSys)))
	_ = m.storage.Store(*metrics.NewGauge("HeapAlloc", float64(m.memStats.HeapAlloc)))
	_ = m.storage.Store(*metrics.NewGauge("HeapIdle", float64(m.memStats.HeapIdle)))
	_ = m.storage.Store(*metrics.NewGauge("HeapInuse", float64(m.memStats.HeapInuse)))
	_ = m.storage.Store(*metrics.NewGauge("HeapObjects", float64(m.memStats.HeapObjects)))
	_ = m.storage.Store(*metrics.NewGauge("HeapReleased", float64(m.memStats.HeapReleased)))
	_ = m.storage.Store(*metrics.NewGauge("HeapSys", float64(m.memStats.HeapSys)))
	_ = m.storage.Store(*metrics.NewGauge("LastGC", float64(m.memStats.LastGC)))
	_ = m.storage.Store(*metrics.NewGauge("Lookups", float64(m.memStats.Lookups)))
	_ = m.storage.Store(*metrics.NewGauge("MCacheInuse", float64(m.memStats.MCacheInuse)))
	_ = m.storage.Store(*metrics.NewGauge("MCacheSys", float64(m.memStats.MCacheSys)))
	_ = m.storage.Store(*metrics.NewGauge("MSpanInuse", float64(m.memStats.MSpanInuse)))
	_ = m.storage.Store(*metrics.NewGauge("MSpanSys", float64(m.memStats.MSpanSys)))
	_ = m.storage.Store(*metrics.NewGauge("Mallocs", float64(m.memStats.Mallocs)))
	_ = m.storage.Store(*metrics.NewGauge("NextGC", float64(m.memStats.NextGC)))
	_ = m.storage.Store(*metrics.NewGauge("NumForcedGC", float64(m.memStats.NumForcedGC)))
	_ = m.storage.Store(*metrics.NewGauge("NumGC", float64(m.memStats.NumGC)))
	_ = m.storage.Store(*metrics.NewGauge("OtherSys", float64(m.memStats.OtherSys)))
	_ = m.storage.Store(*metrics.NewGauge("PauseTotalNs", float64(m.memStats.PauseTotalNs)))
	_ = m.storage.Store(*metrics.NewGauge("StackInuse", float64(m.memStats.StackInuse)))
	_ = m.storage.Store(*metrics.NewGauge("StackSys", float64(m.memStats.StackSys)))
	_ = m.storage.Store(*metrics.NewGauge("Sys", float64(m.memStats.Sys)))
	_ = m.storage.Store(*metrics.NewGauge("TotalAlloc", float64(m.memStats.TotalAlloc)))

	random := rand.Float64() * (10000)
	_ = m.storage.Store(*metrics.NewGauge("RandomValue", random))
}

func (m *Monitor) GatherAdditional() {
	log.Info().Msg("Gathering additional...")
	v, _ := mem.VirtualMemory()
	CPUUtilization, _ := cpu.Percent(0, false)

	_ = m.storage.Store(*metrics.NewGauge("TotalMemory", float64(v.Total)))
	_ = m.storage.Store(*metrics.NewGauge("FreeMemory", float64(v.Free)))
	_ = m.storage.Store(*metrics.NewGauge("CPUutilization1", CPUUtilization[0]))
}

func (m *Monitor) Send() {
	log.Debug().Msg("Reporting...")
	collection, err := m.storage.GetCollection()
	if err != nil {
		log.Error().Err(err).Send()
		return
	}

	if len(collection) == 0 {
		log.Error().Msg("Empty collection, nothing to send to the server")
		return
	}

	if err = m.client.UpdateMetricsBatchJSON(collection); err != nil {
		log.Error().Err(err).Send()
		return
	}

	log.Info().Msg("A collection of metrics was sent for update")

	if err = m.storage.CleanUp(); err != nil {
		log.Error().Err(err).Send()
		return
	}
}
