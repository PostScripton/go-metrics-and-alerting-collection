package agent

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/monitoring"
	"sync"
	"time"
)

type AgentRunner interface {
	RunPolling(interval time.Duration)
	RunReporting(interval time.Duration)
}

type metricAgent struct {
	wg      *sync.WaitGroup
	monitor monitoring.Monitorer
}

func NewMetricAgent(monitor monitoring.Monitorer) AgentRunner {
	return &metricAgent{
		wg:      &sync.WaitGroup{},
		monitor: monitor,
	}
}

func (a *metricAgent) RunPolling(interval time.Duration) {
	pollInterval := time.NewTicker(interval)
	for {
		<-pollInterval.C
		a.wg.Add(2)
		go func() {
			a.monitor.GatherMain()
			a.wg.Done()
		}()
		go func() {
			a.monitor.GatherAdditional()
			a.wg.Done()
		}()
	}
}

func (a *metricAgent) RunReporting(interval time.Duration) {
	reportInterval := time.NewTicker(interval)
	for {
		<-reportInterval.C
		a.wg.Wait()
		a.monitor.Send()
	}
}
