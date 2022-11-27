package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"os"
)

type FileStorage struct {
	path    string
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{
		path: path,
	}
}

func (fs *FileStorage) Get(_ metrics.Metrics) (*metrics.Metrics, error) {
	return nil, fmt.Errorf("no implementation")
}

func (fs *FileStorage) GetCollection() (map[string]metrics.Metrics, error) {
	if err := fs.OpenFile(); err != nil {
		return nil, err
	}

	collection := map[string]metrics.Metrics{}
	if err := fs.decoder.Decode(&collection); err != nil {
		return nil, fmt.Errorf("fetching collection from file: %w", err)
	}

	if err := fs.CloseFile(); err != nil {
		return nil, err
	}

	return collection, nil
}

func (fs *FileStorage) Store(metric metrics.Metrics) error {
	if err := fs.OpenFile(); err != nil {
		return err
	}

	encErr := fs.encoder.Encode(metric)

	if err := fs.CloseFile(); err != nil {
		return err
	}

	return encErr
}

func (fs *FileStorage) StoreCollection(metrics map[string]metrics.Metrics) error {
	if err := fs.OpenFile(); err != nil {
		return err
	}

	encErr := fs.encoder.Encode(metrics)

	if err := fs.CloseFile(); err != nil {
		return err
	}

	return encErr
}

func (fs *FileStorage) OpenFile() error {
	if fs.path == "" {
		return errors.New("empty path")
	}

	file, err := os.OpenFile(fs.path, os.O_RDWR|os.O_CREATE, 0744)
	if err != nil {
		return err
	}

	fs.file = file
	fs.encoder = json.NewEncoder(fs.file)
	fs.decoder = json.NewDecoder(fs.file)

	return nil
}

func (fs *FileStorage) CloseFile() error {
	return fs.file.Close()
}

func (fs *FileStorage) CleanUp() error {
	return fs.file.Truncate(0)
}

func (fs *FileStorage) Ping(_ context.Context) error {
	if _, err := os.Stat(fs.path); errors.Is(err, os.ErrNotExist) {
		return os.ErrNotExist
	}
	return nil
}

func (fs *FileStorage) Close() {
}