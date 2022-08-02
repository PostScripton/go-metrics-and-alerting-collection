package main

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type storager interface {
	repository.Getter
	repository.Storer
}

func main() {
	const port = "8080"
	address := fmt.Sprintf(":%s", port)

	storage := memory.New()

	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	registerRoutes(router, storage)

	fmt.Printf("The server has just started on port [%s]\n", port)
	log.Fatal(http.ListenAndServe(address, router))
}

func registerRoutes(router *chi.Mux, storage storager) {
	router.NotFound(handlers.NotFound)
	router.MethodNotAllowed(handlers.MethodNotAllowed)

	router.Get("/value/{type}/{name}", handlers.GetMetricHandler(storage))
	router.Post("/update/{type}/{name}/{value}", handlers.UpdateMetricHandler(storage))
	router.Post("/value", handlers.GetMetricJSONHandler(storage))
	router.Post("/update", handlers.UpdateMetricJSONHandler(storage))
}
