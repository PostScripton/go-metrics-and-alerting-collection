package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type AgentConfig struct {
	CommonConfig
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

const defaultReportInterval = 10 * time.Second
const defaultPollInterval = 2 * time.Second

func NewAgentConfig() *AgentConfig {
	var cfg AgentConfig

	flag.StringVar(&cfg.Address, "a", defaultAddress, "An address of the server")
	flag.DurationVar(&cfg.ReportInterval, "r", defaultReportInterval, "An interval for reporting to the server")
	flag.DurationVar(&cfg.PollInterval, "p", defaultPollInterval, "An interval for polling metrics data")
	flag.StringVar(&cfg.Key, "k", defaultKey, "A key for encrypting data")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", defaultCryptoKey, "A public key file")

	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Parsing env")
		return nil
	}

	log.Info().Interface("config", cfg).Send()
	return &cfg
}
