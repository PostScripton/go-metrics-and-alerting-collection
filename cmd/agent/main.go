package main

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent/config"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/monitoring"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"time"
)

func main() {
	cfg := config.NewConfig()
	fmt.Printf("Config: %v\n", cfg)

	baseURI := fmt.Sprintf("http://%s", cfg.Address)

	storage := memory.NewMemoryStorage()
	sender := client.NewClient(baseURI, 5*time.Second)
	monitor := monitoring.NewMonitor(storage, sender)

	var metricsAgent agent.MetricAgenter = agent.NewMetricAgent(monitor)
	go metricsAgent.RunPolling(cfg.PollInterval)
	go metricsAgent.RunReporting(cfg.ReportInterval)

	select {}
}
