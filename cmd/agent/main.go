package main

import (
	"flag"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/agent"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/client"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/monitor"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/caarlos0/env/v6"
	"time"
)

type Agent interface {
	RunPolling(interval time.Duration)
	RunReporting(interval time.Duration)
}

type config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

var cfg config

func init() {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "An address of the server")
	flag.DurationVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "An interval for reporting to the server")
	flag.DurationVar(&cfg.PollInterval, "p", cfg.PollInterval, "An interval for polling metrics data")
}

func main() {
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Parsing environment variables error: %s\n", err)
		return
	}
	flag.Parse()

	baseURI := fmt.Sprintf("http://%s", cfg.Address)

	storage := memory.New()
	sender := client.New(baseURI, 5*time.Second)
	metrics := monitor.New(storage, sender)

	var metricsAgent Agent = agent.New(metrics)
	go metricsAgent.RunPolling(cfg.PollInterval)
	go metricsAgent.RunReporting(cfg.ReportInterval)

	select {}
}
