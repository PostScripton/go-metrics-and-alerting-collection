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
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
}

var cfg config

func init() {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "An address of the server")
	flag.DurationVar(&cfg.ReportInterval, "r", 10*time.Second, "An interval for reporting to the server")
	flag.DurationVar(&cfg.PollInterval, "p", 2*time.Second, "An interval for polling metrics data")
	flag.StringVar(&cfg.Key, "k", "", "A key for encrypting data")
}

func main() {
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Parsing environment variables error: %s\n", err)
		return
	}
	fmt.Printf("Config: %v\n", cfg)

	baseURI := fmt.Sprintf("http://%s", cfg.Address)

	storage := memory.New()
	sender := client.New(baseURI, 5*time.Second, cfg.Key)
	metrics := monitor.New(storage, sender)

	var metricsAgent Agent = agent.New(metrics)
	go metricsAgent.RunPolling(cfg.PollInterval)
	go metricsAgent.RunReporting(cfg.ReportInterval)

	select {}
}
