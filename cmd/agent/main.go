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

const notAssigned = "N/A"

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

// go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildTime=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git rev-parse HEAD)'" cmd/agent/main.go

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "03:04:05PM"})

	if buildVersion == "" {
		buildVersion = notAssigned
	}
	if buildTime == "" {
		buildTime = notAssigned
	}
	if buildCommit == "" {
		buildCommit = notAssigned
	}

	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildTime)
	log.Printf("Build commit: %s", buildCommit)

	cfg := config.NewAgentConfig()

	baseURI := fmt.Sprintf("http://%s", cfg.Address)

	storage := memory.NewMemoryStorage()
	sender := client.NewClient(baseURI, 5*time.Second, cfg.Key, cfg.CryptoKey)
	monitor := monitoring.NewMonitor(storage, sender)

	metricsAgent := agent.NewMetricAgent(monitor)
	go metricsAgent.RunPolling(cfg.PollInterval)
	go metricsAgent.RunReporting(cfg.ReportInterval)

	select {}
}
