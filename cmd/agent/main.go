package main

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent/config"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/monitoring"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "03:04:05PM"})

	cfg := config.NewConfig()
	log.Info().Interface("config", cfg).Send()

	baseURI := fmt.Sprintf("http://%s", cfg.Address)

	storage := memory.NewMemoryStorage()
	sender := client.NewClient(baseURI, 5*time.Second, cfg.Key)
	monitor := monitoring.NewMonitor(storage, sender)

	var metricsAgent agent.MetricAgenter = agent.NewMetricAgent(monitor)
	go metricsAgent.RunPolling(cfg.PollInterval)
	go metricsAgent.RunReporting(cfg.ReportInterval)

	select {}
}
