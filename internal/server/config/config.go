package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"time"
)

type config struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func NewConfig() *config {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Parsing .env error: %s\n", err)
		return nil
	}

	return &cfg
}
