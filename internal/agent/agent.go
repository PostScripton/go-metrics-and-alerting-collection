package agent

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent/service"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"time"
)

type metricAgent struct {
	storage service.Storage
	service *service.MetricService
}

func New(storage service.Storage, client client.Client) *metricAgent {
	return &metricAgent{
		storage: storage,
		service: service.NewMetricService(storage, client),
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
