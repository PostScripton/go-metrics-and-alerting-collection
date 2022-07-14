package main

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/handlers"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	storage := memory.New()
	metricService := service.NewMetricService(storage)

	router := gin.New()
	router.RedirectTrailingSlash = false
	//router.LoadHTMLGlob("../../internal/templates/**/*")
	//router.GET("/", handlers.GetAllMetricsHandler(storage))
	router.GET("/value/:type/:name", handlers.GetMetricHandler(storage))
	router.POST("/update/:type/:name/:value", handlers.UpdateMetricHandler(metricService))

	fmt.Println("The server has just started on port :8080")
	log.Fatal(router.Run(":8080"))
}
