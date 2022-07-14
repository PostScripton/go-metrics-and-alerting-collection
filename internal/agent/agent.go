package agent

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent/service"
	"time"
)

type metricAgent struct {
	service *service.MetricService
}

func New(metricService *service.MetricService) *metricAgent {
	return &metricAgent{
		service: metricService,
	}
}

func (a *metricAgent) RunPolling() {
	pollInterval := time.NewTicker(2 * time.Second)
	for {
		<-pollInterval.C
		a.service.GatherMetrics()
	}
}

func (a *metricAgent) RunReporting() {
	reportInterval := time.NewTicker(10 * time.Second)
	for {
		<-reportInterval.C
		a.service.SendMetrics()
	}
}
