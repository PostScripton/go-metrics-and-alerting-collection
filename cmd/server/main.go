package main

import (
	"flag"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/file"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/handlers"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

type storager interface {
	repository.Getter
	repository.Storer
	repository.CollectionGetter
	repository.CollectionStorer
}

type config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

var cfg config

func init() {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "An address of the server")
	flag.BoolVar(&cfg.Restore, "r", true, "Whether restore state from a file")
	flag.StringVar(&cfg.StoreFile, "f", "/tmp/devops-metrics-db.json", "A file to store to or restore from")
	flag.DurationVar(&cfg.StoreInterval, "i", 5*time.Minute, "An interval for storing into a file")
}

func main() {
	flag.Parse()
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Parsing environment variables error: %s\n", err)
		return
	}
	fmt.Printf("Config: %v\n", cfg)

	var storage storager = memory.New()
	fileStorage := file.New(cfg.StoreFile)

	if cfg.Restore {
		if err := fileStorage.Restore(storage); err != nil {
			fmt.Printf("Restore error: %s\n", err)
		}
	}

	if cfg.StoreInterval == 0 {
		fmt.Println("Synchronously save to disk")
		// todo не знаю как сделать, чтобы сохраняло синхронно
	} else {
		fmt.Printf("Asynchronous save to disk with [%s] interval\n", cfg.StoreInterval)
		go file.RunStoring(cfg.StoreInterval, storage, fileStorage)
	}

	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	registerRoutes(router, storage)

	fmt.Printf("The server has just started on [%s]\n", cfg.Address)
	log.Fatal(http.ListenAndServe(cfg.Address, router))
}

func registerRoutes(router *chi.Mux, storage storager) {
	router.NotFound(handlers.NotFound)
	router.MethodNotAllowed(handlers.MethodNotAllowed)

	router.Get("/value/{type}/{name}", handlers.GetMetricHandler(storage))
	router.Post("/update/{type}/{name}/{value}", handlers.UpdateMetricHandler(storage))
	router.Post("/value", handlers.GetMetricJSONHandler(storage))
	router.Post("/update", handlers.UpdateMetricJSONHandler(storage))
}
