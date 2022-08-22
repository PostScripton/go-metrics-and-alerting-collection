package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"os"
	"time"
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

func (fs *fileStorage) GetCollection() map[string]metrics.Metrics {
	collection := map[string]metrics.Metrics{}
	if err := fs.decoder.Decode(&collection); err != nil {
		fmt.Printf("Error when fetching from file: %s\n", err)
		return make(map[string]metrics.Metrics)
	}
	return collection
}

func (fs *fileStorage) Store(metric metrics.Metrics) error {
	return fs.encoder.Encode(metric)
}

func (fs *fileStorage) StoreCollection(metrics map[string]metrics.Metrics) error {
	return fs.encoder.Encode(metrics)
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

func RunStoring(interval time.Duration, from repository.CollectionGetter, to *fileStorage) {
	storeInterval := time.NewTicker(interval)
	for {
		<-storeInterval.C

		if err := to.Open(); err != nil {
			fmt.Println(err)
			return
		}
		if err := to.StoreCollection(from.GetCollection()); err != nil {
			fmt.Println(err)
			return
		}
		if err := to.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (fs *fileStorage) Restore(to repository.CollectionStorer) error {
	if err := fs.Open(); err != nil {
		return err
	}
	if err := to.StoreCollection(fs.GetCollection()); err != nil {
		return err
	}

	return nil
}
