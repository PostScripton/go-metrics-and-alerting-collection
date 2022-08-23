package agent

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/monitoring"
	"time"
)

type MetricAgenter interface {
	RunPolling(interval time.Duration)
	RunReporting(interval time.Duration)
}

type metricAgent struct {
	monitor monitoring.Monitorer
}

func NewMetricAgent(monitor monitoring.Monitorer) *metricAgent {
	return &metricAgent{
		monitor: monitor,
	}
}

func (a *metricAgent) RunPolling(interval time.Duration) {
	pollInterval := time.NewTicker(interval)
	for {
		<-pollInterval.C
		a.monitor.Gather()
	}
}

func (a *metricAgent) RunReporting(interval time.Duration) {
	reportInterval := time.NewTicker(interval)
	for {
		<-reportInterval.C
		a.monitor.Send()
	}
}
