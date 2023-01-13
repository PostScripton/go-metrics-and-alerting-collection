package client

import (
	"io"
	"time"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server"
	"github.com/PostScripton/go-metrics-and-alerting-collection/pkg/key_management/rsakeys"
	"github.com/rs/zerolog/log"
)

type IClient interface {
	io.Closer                                                       // Закрывает соединение
	UpdateMetric(metrics metrics.Metrics) error                     // Обновляет метрику
	BatchUpdateMetrics(collection map[string]metrics.Metrics) error // Обновляет коллекцию метрик
}

func NewClient(clientType string, baseURI string, timeout time.Duration, hashKey string, cryptoKey string, address string) IClient {
	publicKey, err := rsakeys.ImportPublicKeyFromFile(cryptoKey)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get public hashKey from file")
	}

	switch clientType {
	case server.HTTPType:
		return newHTTPClient(baseURI, timeout, hashKey, publicKey, address)
	case server.GRPCType:
		return newGRPCClient(baseURI)
	}

	return nil
}
