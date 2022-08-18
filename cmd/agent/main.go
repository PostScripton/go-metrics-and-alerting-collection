package main

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent/config"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/monitor"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"time"
)

func main() {
	cfg := config.NewConfig()

	baseURI := fmt.Sprintf("http://%s", cfg.Address)

	storage := memory.New()
	sender := client.New(baseURI, 5*time.Second)
	metrics := monitor.New(storage, sender)

	var metricsAgent agent.Agent = agent.New(metrics)
	go metricsAgent.RunPolling(cfg.PollInterval)
	go metricsAgent.RunReporting(cfg.ReportInterval)

	select {}
}
