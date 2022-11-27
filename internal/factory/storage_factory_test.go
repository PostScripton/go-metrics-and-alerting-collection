package factory

import (
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/database/postgres"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/file"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorageFactory_CreateStorage(t *testing.T) {
	type fields struct {
		DSN      string
		FilePath string
	}
	tests := []struct {
		name   string
		fields fields
		want   storage.Storager
	}{
		{
			name:   "memory storage",
			fields: fields{},
			want:   &memory.MemoryStorage{},
		},
		{
			name: "file storage",
			fields: fields{
				FilePath: "/some/path",
			},
			want: &file.FileStorage{},
		},
		{
			name: "postgres storage",
			fields: fields{
				DSN: "postgres://postgres:postgres@postgres:5432/db_name",
			},
			want: &postgres.Postgres{},
		},
		{
			name: "postgres storage",
			fields: fields{
				DSN:      "postgres://postgres:postgres@postgres:5432/db_name",
				FilePath: "/some/path",
			},
			want: &postgres.Postgres{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &StorageFactory{
				DSN:      tt.fields.DSN,
				FilePath: tt.fields.FilePath,
				Testing:  true,
			}
			if got := sf.CreateStorage(); !assert.IsTypef(t, tt.want, got, "") {
				t.Errorf("CreateStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
