package main

import (
	"context"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/PostScripton/go-metrics-and-alerting-collection/config"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server"
)

const notAssigned = "N/A"

var (
	buildVersion string
	buildTime    string
	buildCommit  string
)

// go run -ldflags "-X main.buildVersion=v1.0.0 -X 'main.buildTime=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git rev-parse HEAD)'" cmd/server/main.go

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "03:04:05PM"})

	if buildVersion == "" {
		buildVersion = notAssigned
	}
	if buildTime == "" {
		buildTime = notAssigned
	}
	if buildCommit == "" {
		buildCommit = notAssigned
	}

	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildTime)
	log.Printf("Build commit: %s", buildCommit)

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
	restorer := storage.NewRestorer(backupStorage, mainStorage)
	restorer.Run(cfg.Restore, cfg.StoreInterval)

	coreServer := server.NewServer(cfg.Address, mainStorage, cfg.Key)
	coreServer.Run()
}
