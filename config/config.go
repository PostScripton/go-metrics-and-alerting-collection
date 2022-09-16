package config

type CommonConfig struct {
	Address string `env:"ADDRESS"`
	Key     string `env:"KEY"`
}

const defaultAddress = "localhost:8080"
const defaultKey = ""
