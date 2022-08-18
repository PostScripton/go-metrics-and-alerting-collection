package main

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/config"
)

func main() {
	cfg := config.NewConfig()

	storage := memory.NewMemoryStorage()

	coreServer := server.NewServer(cfg.Address, storage)
	coreServer.Run()
}
