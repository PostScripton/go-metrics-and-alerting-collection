package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/PostScripton/go-metrics-and-alerting-collection/config"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/monitoring"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "03:04:05PM"})

	cfg := config.NewAgentConfig()

	baseURI := fmt.Sprintf("http://%s", cfg.Address)

	storage := memory.NewMemoryStorage()
	sender := client.NewClient(baseURI, 5*time.Second, cfg.Key)
	monitor := monitoring.NewMonitor(storage, sender)

	metricsAgent := agent.NewMetricAgent(monitor)
	go metricsAgent.RunPolling(cfg.PollInterval)
	go metricsAgent.RunReporting(cfg.ReportInterval)

	select {}
}
