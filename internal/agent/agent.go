package agent

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent/service"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"time"
)

type metricAgent struct {
	storage repository.Storage
	service *service.MetricService
}

func New(storage repository.Storage, client client.Client) *metricAgent {
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
