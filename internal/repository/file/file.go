package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"os"
)

type fileStorage struct {
	path    string
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewFileStorage(path string) *fileStorage {
	return &fileStorage{
		path: path,
	}
}

func (fs *fileStorage) Get(_ metrics.Metrics) (*metrics.Metrics, error) {
	return nil, fmt.Errorf("no implementation")
}

func (fs *fileStorage) GetCollection() (map[string]metrics.Metrics, error) {
	if err := fs.Open(); err != nil {
		return nil, err
	}

	collection := map[string]metrics.Metrics{}
	if err := fs.decoder.Decode(&collection); err != nil {
		return nil, fmt.Errorf("fetching collection from file: %w", err)
	}

	if err := fs.Close(); err != nil {
		return nil, err
	}

	return collection, nil
}

func (fs *fileStorage) Store(metric metrics.Metrics) error {
	if err := fs.Open(); err != nil {
		return err
	}

	encErr := fs.encoder.Encode(metric)

	if err := fs.Close(); err != nil {
		return err
	}

	return encErr
}

func (fs *fileStorage) StoreCollection(metrics map[string]metrics.Metrics) error {
	if err := fs.Open(); err != nil {
		return err
	}

	encErr := fs.encoder.Encode(metrics)

	if err := fs.Close(); err != nil {
		return err
	}

	return encErr
}

func (fs *fileStorage) Open() error {
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

func (fs *fileStorage) Close() error {
	return fs.file.Close()
}

func (fs *fileStorage) CleanUp() error {
	return fs.file.Truncate(0)
}
