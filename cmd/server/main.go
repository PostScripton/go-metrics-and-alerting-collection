package main

import (
	"context"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/database/postgres"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/config"
	"time"
)

func main() {
	cfg := config.NewConfig()
	fmt.Printf("Config: %v\n", cfg)

	dbConn, err := postgres.ConnectToDB(context.Background(), cfg.DatabaseDSN)
	if dbConn != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err = dbConn.Ping(ctx); err != nil {
			fmt.Printf("Ping DB err: %s\n", err)
		}
		defer dbConn.Close()

		postgres.Migrate(dbConn)
	} else {
		fmt.Printf("Postgres err: %s\n", err)
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
