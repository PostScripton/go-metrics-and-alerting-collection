package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
	"time"
)

type config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
}

const defaultAddress = "localhost:8080"
const defaultReportInterval = 10 * time.Second
const defaultPollInterval = 2 * time.Second
const defaultKey = ""

func NewConfig() *config {
	var cfg config

	flag.StringVar(&cfg.Address, "a", defaultAddress, "An address of the server")
	flag.DurationVar(&cfg.ReportInterval, "r", defaultReportInterval, "An interval for reporting to the server")
	flag.DurationVar(&cfg.PollInterval, "p", defaultPollInterval, "An interval for polling metrics data")
	flag.StringVar(&cfg.Key, "k", defaultKey, "A key for encrypting data")

	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Parsing env")
		return nil
	}

	return &cfg
}
