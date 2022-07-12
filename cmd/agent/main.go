package main

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
)

func main() {
	storage := memory.New()
	sender := client.New("http://localhost:8080")

	metricsAgent := agent.New(storage, sender)

	go metricsAgent.RunPolling()
	go metricsAgent.RunReporting()

	select {}
}
