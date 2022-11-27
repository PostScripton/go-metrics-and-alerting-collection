package factory

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/database/postgres"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/file"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/memory"
	"github.com/rs/zerolog/log"
)

type StorageFactory struct {
	DSN      string
	FilePath string
	Testing  bool
}

func (sf *StorageFactory) CreateStorage() storage.Storager {
	if sf.DSN != "" {
		if sf.Testing {
			return &postgres.Postgres{}
		}

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
