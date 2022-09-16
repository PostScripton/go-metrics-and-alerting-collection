package main

import (
	"context"
	"github.com/PostScripton/go-metrics-and-alerting-collection/config"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "03:04:05PM"})

	cfg := config.NewServerConfig()

	mainStorageFactory := &factory.StorageFactory{DSN: cfg.DatabaseDSN}
	mainStorage := mainStorageFactory.CreateStorage()
	pingCtx, cancelPing := context.WithTimeout(context.Background(), 1*time.Second)
	if err := mainStorage.Ping(pingCtx); err != nil {
		log.Warn().Err(err).Msg("Ping storage")
	} else {
		defer cancelPing()
	}
	defer mainStorage.Close()

	backupStorageFactory := &factory.StorageFactory{
		DSN:      cfg.DatabaseDSN,
		FilePath: cfg.StoreFile,
	}
	backupStorage := backupStorageFactory.CreateStorage()
	restorer := repository.NewRestorer(backupStorage, mainStorage)
	restorer.Run(cfg.Restore, cfg.StoreInterval)

	coreServer := server.NewServer(cfg.Address, mainStorage, cfg.Key)
	coreServer.Run()
}
