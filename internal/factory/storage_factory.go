package factory

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/database/postgres"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/file"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/rs/zerolog/log"
)

type StorageFactory struct {
	DSN      string
	FilePath string
}

func (sf *StorageFactory) CreateStorage() repository.Storager {
	if sf.DSN != "" {
		db, err := postgres.NewPostgres(sf.DSN)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		return db
	}
	if sf.FilePath != "" {
		return file.NewFileStorage(sf.FilePath)
	}
	return memory.NewMemoryStorage()
}
