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

func (a *metricAgent) RunPolling() {
	pollInterval := time.NewTicker(2 * time.Second)
	for {
		<-pollInterval.C
		a.monitor.Gather()
	}
}

func (a *metricAgent) RunReporting() {
	reportInterval := time.NewTicker(10 * time.Second)
	for {
		<-reportInterval.C
		a.monitor.Send()
	}
}
