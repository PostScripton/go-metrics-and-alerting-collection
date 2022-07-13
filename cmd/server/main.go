package main

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/handlers"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/service"
	"log"
	"net/http"
)

func main() {
	storage := memory.New()
	metricService := service.NewMetricService(storage)

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handlers.UpdateMetricHandler(metricService))

	fmt.Println("The server has just started on port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
