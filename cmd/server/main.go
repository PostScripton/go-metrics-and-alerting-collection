package main

import (
	"context"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/database"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/file"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/config"
)

func main() {
	cfg := config.NewConfig()
	fmt.Printf("Config: %v\n", cfg)

	memoryStorage := memory.NewMemoryStorage()
	fileStorage := file.NewFileStorage(cfg.StoreFile)
	db, dbErr := database.NewPostgres(context.Background(), cfg.DatabaseDSN)
	if dbErr != nil {
		fmt.Printf("Postgres connection err: %s\n", dbErr)
	} else {
		cancelPing, pingErr := db.Ping(context.Background())
		if pingErr != nil {
			fmt.Printf("Ping DB err: %s\n", pingErr)
		}
		defer cancelPing()
		defer func() {
			if err := db.Close(context.Background()); err != nil {
				fmt.Printf("Postgres close err: %s\n", err)
			}
		}()
	}

	var backupStorage repository.Storager
	if db != nil {
		backupStorage = db
		fmt.Println("backup storage is DB")
	} else {
		backupStorage = fileStorage
		fmt.Println("backup storage is file")
	}
	restorer := repository.NewRestorer(backupStorage, memoryStorage)
	restorer.Run(cfg.Restore, cfg.StoreInterval)

	coreServer := server.NewServer(cfg.Address, memoryStorage, db, cfg.Key)
	coreServer.Run()
}
