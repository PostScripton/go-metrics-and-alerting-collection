package config

type CommonConfig struct {
	Address   string `env:"ADDRESS"`
	Key       string `env:"KEY"`
	CryptoKey string `env:"CRYPTO_KEY"`
}

const defaultAddress = "localhost:8080"
const defaultKey = ""
const defaultCryptoKey = ""
