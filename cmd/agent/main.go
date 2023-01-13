package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PostScripton/go-metrics-and-alerting-collection/config"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/monitoring"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

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
	sender := client.NewClient(cfg.ServerType, baseURI, 5*time.Second, cfg.Key, cfg.CryptoKey, cfg.Address)
	defer sender.Close()
	monitor := monitoring.NewMonitor(storage, sender)

	metricsAgent := agent.NewMetricAgent(monitor)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		metricsAgent.RunPolling(gCtx, cfg.PollInterval.Duration)
		return nil
	})
	g.Go(func() error {
		metricsAgent.RunReporting(gCtx, cfg.ReportInterval.Duration)
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()

		monitor.Send()
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Info().Err(err).Msg("Reason for graceful shutdown")
	}

	log.Info().Msg("The application is shutdown")
}
