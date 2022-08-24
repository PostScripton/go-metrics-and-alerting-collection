package main

import (
	"fmt"
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

	if cfg.Restore {
		if err := fileStorage.Restore(memoryStorage); err != nil {
			fmt.Printf("Restore error: %s\n", err)
		}
	}

	if cfg.StoreInterval == 0 {
		fmt.Println("Synchronously save to disk")
		// todo не знаю как сделать, чтобы сохраняло синхронно
	} else {
		fmt.Printf("Asynchronous save to disk with [%s] interval\n", cfg.StoreInterval)
		go file.RunStoring(cfg.StoreInterval, memoryStorage, fileStorage)
	}

	coreServer := server.NewServer(cfg.Address, memoryStorage, cfg.Key)
	coreServer.Run()
}
