package server

import (
	"context"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage"
	"github.com/PostScripton/go-metrics-and-alerting-collection/pkg/key_management/rsakeys"
	"github.com/rs/zerolog/log"
)

const (
	HTTPType = "http"
	GRPCType = "grpc"
)

type IServer interface {
	Run() error                         // Запускает сервер
	Shutdown(ctx context.Context) error // Останавливает сервер
}

func NewServer(serverType string, address string, storage storage.Storager, hashKey, cryptoKey, trustedSubnet string) IServer {
	privateKey, err := rsakeys.ImportPrivateKeyFromFile(cryptoKey)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get private key from file")
	}

	switch serverType {
	case HTTPType:
		return newHTTPServer(address, storage, hashKey, privateKey, trustedSubnet)
	case GRPCType:
		return newGRPCServer(address, storage, hashKey)
	}

	return nil
}
