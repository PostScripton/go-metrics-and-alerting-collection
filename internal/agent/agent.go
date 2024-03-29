// Package agent позволяет собирать метрики с компьютера на котором он запущен и отправлять их на сервер на сохранение.
package agent

import (
	"context"
	"sync"
	"time"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/monitoring"
)

// Runner интерфейс агента управляет сбором и отправкой метрик
type Runner interface {
	RunPolling(interval time.Duration)   // Запускает сбор метрик раз в какой-то интервал
	RunReporting(interval time.Duration) // Запускает отправку метрик раз в какой-то интервал
}

type MetricAgent struct {
	wg      *sync.WaitGroup
	monitor monitoring.IMonitor
}

func NewMetricAgent(monitor monitoring.IMonitor) *MetricAgent {
	return &MetricAgent{
		wg:      &sync.WaitGroup{},
		monitor: monitor,
	}
}

func (a *MetricAgent) RunPolling(ctx context.Context, interval time.Duration) {
	pollInterval := time.NewTicker(interval)
	for {
		select {
		case <-pollInterval.C:
			a.wg.Add(2)
			go func() {
				a.monitor.GatherMain()
				a.wg.Done()
			}()
			go func() {
				a.monitor.GatherAdditional()
				a.wg.Done()
			}()
		case <-ctx.Done():
			return
		}
	}
}

func (a *MetricAgent) RunReporting(ctx context.Context, interval time.Duration) {
	reportInterval := time.NewTicker(interval)
	for {
		select {
		case <-reportInterval.C:
			a.wg.Wait()
			a.monitor.Send()
		case <-ctx.Done():
			return
		}
	}
}
