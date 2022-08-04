package main

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/handlers"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type storager interface {
	repository.Getter
	repository.Storer
}

type config struct {
	Address string `env:"ADDRESS" envDefault:"localhost:8080"`
}

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("Parsing .env error: %s\n", err)
		return
	}

	storage := memory.New()

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
