package main

import (
	"context"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/database/postgres"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "03:04:05PM"})

	cfg := config.NewConfig()
	log.Info().Interface("config", cfg).Send()

	dbConn, err := postgres.ConnectToDB(context.Background(), cfg.DatabaseDSN)
	if dbConn != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err = dbConn.Ping(ctx); err != nil {
			log.Warn().Err(err).Msg("Ping DB err")
		}
		defer dbConn.Close()

		postgres.Migrate(dbConn)
	} else {
		log.Warn().Err(err).Msg("Postgres err")
	}

	mainStorageFactory := &factory.StorageFactory{
		Pool: dbConn,
		DSN:  cfg.DatabaseDSN,
	}
	mainStorage := mainStorageFactory.CreateStorage()

	backupStorageFactory := &factory.StorageFactory{
		Pool:     dbConn,
		DSN:      cfg.DatabaseDSN,
		FilePath: cfg.StoreFile,
	}
	backupStorage := backupStorageFactory.CreateStorage()
	restorer := repository.NewRestorer(backupStorage, mainStorage)
	restorer.Run(cfg.Restore, cfg.StoreInterval)

	coreServer := server.NewServer(cfg.Address, mainStorage, cfg.Key, dbConn)
	coreServer.Run()
}
