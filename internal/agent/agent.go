package agent

import (
	"time"
)

type Monitorer interface {
	Gather()
	Send()
}

type metricAgent struct {
	monitor Monitorer
}

func New(monitor Monitorer) *metricAgent {
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
