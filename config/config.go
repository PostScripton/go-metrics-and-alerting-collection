package config

import "github.com/PostScripton/go-metrics-and-alerting-collection/internal/server"

type CommonConfig struct {
	ServerType string `env:"SERVER_TYPE" json:"server_type"`
	Address    string `env:"ADDRESS" json:"address"`
	Key        string `env:"KEY" json:"-"`
	CryptoKey  string `env:"CRYPTO_KEY" json:"crypto_key"`
}

const defaultServerType = server.HTTPType
const defaultAddress = "localhost:8080"
const defaultKey = ""
const defaultCryptoKey = ""
