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

type ServerConfig struct {
	CommonConfig
	StoreInterval types.Duration `env:"STORE_INTERVAL" json:"store_interval"`
	StoreFile     string         `env:"STORE_FILE" json:"store_file"`
	Restore       bool           `env:"RESTORE" json:"restore"`
	DatabaseDSN   string         `env:"DATABASE_DSN" json:"database_dsn"`
}

const defaultRestore = true
const defaultStoreFile = "/tmp/devops-metrics-db.json"
const defaultStoreInterval = 5 * time.Minute
const defaultDatabaseDSN = ""

func NewServerConfig() *ServerConfig {
	var jsonCfg ServerConfig
	var flagCfg ServerConfig
	var envCfg ServerConfig

	flag.StringVar(&flagCfg.Address, "a", defaultAddress, "An address of the server")
	flag.BoolVar(&flagCfg.Restore, "r", defaultRestore, "Whether restore state from a file")
	flag.StringVar(&flagCfg.StoreFile, "f", defaultStoreFile, "A file to store to or restore from")
	flag.DurationVar(&flagCfg.StoreInterval.Duration, "i", defaultStoreInterval, "An interval for storing into a file")
	flag.StringVar(&flagCfg.Key, "k", defaultKey, "A key for encrypting data")
	flag.StringVar(&flagCfg.DatabaseDSN, "d", defaultDatabaseDSN, "A DSN for connecting to database")
	flag.StringVar(&flagCfg.CryptoKey, "crypto-key", defaultCryptoKey, "A private key file")

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

	var cfg ServerConfig
	cfg.merge(&envCfg).merge(&flagCfg).merge(&jsonCfg)

	log.Info().Interface("config", cfg).Send()
	return &cfg
}

func (c *ServerConfig) merge(other *ServerConfig) *ServerConfig {
	if c.Address == "" {
		c.Address = other.Address
	}
	if c.StoreInterval.Duration == 0 {
		c.StoreInterval = other.StoreInterval
	}
	if c.StoreFile == "" {
		c.StoreFile = other.StoreFile
	}
	if !c.Restore {
		c.Restore = other.Restore
	}
	if c.DatabaseDSN == "" {
		c.DatabaseDSN = other.DatabaseDSN
	}
	if c.Key == "" {
		c.Key = other.Key
	}
	if c.CryptoKey == "" {
		c.CryptoKey = other.CryptoKey
	}

	return c
}
