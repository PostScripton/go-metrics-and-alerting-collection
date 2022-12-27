package config

type CommonConfig struct {
	Address   string `env:"ADDRESS" json:"address"`
	Key       string `env:"KEY" json:"-"`
	CryptoKey string `env:"CRYPTO_KEY" json:"crypto_key"`
}

const defaultAddress = "localhost:8080"
const defaultKey = ""
const defaultCryptoKey = ""
