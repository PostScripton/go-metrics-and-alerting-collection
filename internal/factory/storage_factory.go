package factory

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/database/postgres"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/file"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/jackc/pgx/v4/pgxpool"
)

type StorageFactory struct {
	Pool     *pgxpool.Pool
	DSN      string
	FilePath string
}

func (sf *StorageFactory) CreateStorage() repository.Storager {
	if sf.DSN != "" && sf.Pool != nil {
		return postgres.NewPostgres(sf.Pool)
	}
	if sf.FilePath != "" {
		return file.NewFileStorage(sf.FilePath)
	}
	return memory.NewMemoryStorage()
}
