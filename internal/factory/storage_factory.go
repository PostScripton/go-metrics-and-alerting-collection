package factory

import (
	"github.com/rs/zerolog/log"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/database/postgres"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/file"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/memory"
)

// StorageFactory реализация абстрактной фабрики для хранилища
type StorageFactory struct {
	DSN      string // Если передана строка для подключения к БД, то будет БД-хранилище
	FilePath string // Если передан путь до файла, то будет файловое хранилище
	Testing  bool   // (Временное решение) Нужно выставить в true для тестов
}

// CreateStorage возвращает новый экземпляр хранилища
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
