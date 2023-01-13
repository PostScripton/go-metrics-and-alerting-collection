package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/types"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type AgentConfig struct {
	CommonConfig
	ReportInterval types.Duration `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   types.Duration `env:"POLL_INTERVAL" json:"poll_interval"`
}

const defaultReportInterval = 10 * time.Second
const defaultPollInterval = 2 * time.Second

func NewAgentConfig() *AgentConfig {
	var jsonCfg AgentConfig
	var flagCfg AgentConfig
	var envCfg AgentConfig

	flag.StringVar(&flagCfg.ServerType, "type", defaultServerType, "A server type: http or grpc")
	flag.StringVar(&flagCfg.Address, "a", defaultAddress, "An address of the server")
	flag.DurationVar(&flagCfg.ReportInterval.Duration, "r", defaultReportInterval, "An interval for reporting to the server")
	flag.DurationVar(&flagCfg.PollInterval.Duration, "p", defaultPollInterval, "An interval for polling metrics data")
	flag.StringVar(&flagCfg.Key, "k", defaultKey, "A key for encrypting data")
	flag.StringVar(&flagCfg.CryptoKey, "crypto-key", defaultCryptoKey, "A public key file")

	var configFile struct {
		Path string `env:"CONFIG"`
	}
	flag.StringVar(&configFile.Path, "c", "", "A path to the JSON config file")
	flag.StringVar(&configFile.Path, "config", "", "A path to the JSON config file")

	flag.Parse()

	if err := env.Parse(&configFile); err != nil {
		log.Fatal().Err(err).Msg("Parsing env to get config file")
		return nil
	}

	if configFile.Path != "" {
		jsonBytes, err := os.ReadFile(configFile.Path)
		if err != nil {
			log.Fatal().Err(err).Msgf("Read JSON config file at path: %s", configFile.Path)
			return nil
		}
		if err = json.Unmarshal(jsonBytes, &jsonCfg); err != nil {
			log.Fatal().Err(err).Msg("Parse JSON from config file")
			return nil
		}
	}

	if err := env.Parse(&envCfg); err != nil {
		log.Fatal().Err(err).Msg("Parsing env")
		return nil
	}

	var cfg AgentConfig
	cfg.merge(&envCfg).merge(&flagCfg).merge(&jsonCfg)

	log.Info().Interface("config", cfg).Send()
	return &cfg
}

func (c *AgentConfig) merge(other *AgentConfig) *AgentConfig {
	if c.ServerType == "" {
		c.ServerType = other.ServerType
	}
	if c.Address == "" {
		c.Address = other.Address
	}
	if c.ReportInterval.Duration == 0 {
		c.ReportInterval = other.ReportInterval
	}
	if c.PollInterval.Duration == 0 {
		c.PollInterval = other.PollInterval
	}
	if c.Key == "" {
		c.Key = other.Key
	}
	if c.CryptoKey == "" {
		c.CryptoKey = other.CryptoKey
	}

	return c
}
