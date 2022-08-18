package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func NewConfig() *config {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Parsing .env error: %s\n", err)
		return nil
	}

	return &cfg
}
