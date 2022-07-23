package main

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/monitor"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"time"
)

type Agent interface {
	RunPolling()
	RunReporting()
}

func main() {
	const baseURI = "http://localhost:8080"

	storage := memory.New()
	sender := client.New(baseURI, 5*time.Second)
	metrics := monitor.New(storage, sender)

	var metricsAgent Agent = agent.New(metrics)
	go metricsAgent.RunPolling()
	go metricsAgent.RunReporting()

	select {}
}
