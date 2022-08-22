package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

type config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

const defaultAddress = "localhost:8080"
const defaultReportInterval = 10 * time.Second
const defaultPollInterval = 2 * time.Second

func NewConfig() *config {
	var cfg config

	flag.StringVar(&cfg.Address, "a", defaultAddress, "An address of the server")
	flag.DurationVar(&cfg.ReportInterval, "r", defaultReportInterval, "An interval for reporting to the server")
	flag.DurationVar(&cfg.PollInterval, "p", defaultPollInterval, "An interval for polling metrics data")

	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Parsing .env error: %s\n", err)
		return nil
	}

	return &cfg
}
