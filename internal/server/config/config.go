package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
	"time"
)

type config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
}

const defaultAddress = "localhost:8080"
const defaultRestore = true
const defaultStoreFile = "/tmp/devops-metrics-db.json"
const defaultStoreInterval = 5 * time.Minute
const defaultKey = ""
const defaultDatabaseDSN = ""

func NewConfig() *config {
	var cfg config

	flag.StringVar(&cfg.Address, "a", defaultAddress, "An address of the server")
	flag.BoolVar(&cfg.Restore, "r", defaultRestore, "Whether restore state from a file")
	flag.StringVar(&cfg.StoreFile, "f", defaultStoreFile, "A file to store to or restore from")
	flag.DurationVar(&cfg.StoreInterval, "i", defaultStoreInterval, "An interval for storing into a file")
	flag.StringVar(&cfg.Key, "k", defaultKey, "A key for encrypting data")
	flag.StringVar(&cfg.DatabaseDSN, "d", defaultDatabaseDSN, "A DSN for connecting to database")

	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Err(err).Msg("Parsing env")
		return nil
	}

	return &cfg
}
