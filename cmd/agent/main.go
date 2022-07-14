package main

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent/service"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
)

func main() {
	storage := memory.New()
	sender := client.New("http://localhost:8080")
	metricService := service.NewMetricService(storage, sender)

	metricsAgent := agent.New(metricService)

	go metricsAgent.RunPolling()
	go metricsAgent.RunReporting()

	select {}
}
