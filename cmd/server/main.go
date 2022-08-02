package main

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	address := fmt.Sprintf(":%s", port)

	storage := memory.New()

	router := chi.NewRouter()
	router.Get("/value/{type}/{name}", handlers.GetMetricHandler(storage))
	router.Post("/update/{type}/{name}/{value}", handlers.UpdateMetricHandler(storage))

	fmt.Printf("The server has just started on port [%s]\n", port)
	log.Fatal(http.ListenAndServe(address, router))
}
