package main

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository/memory"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	const port = "8080"
	address := fmt.Sprintf(":%s", port)

	storage := memory.New()

	router := gin.New()
	router.RedirectTrailingSlash = false
	//router.LoadHTMLGlob("../../internal/templates/**/*")
	//router.GET("/", handlers.GetAllMetricsHandler(storage))
	router.GET("/value/:type/:name", handlers.GetMetricHandler(storage))
	router.POST("/update/:type/:name/:value", handlers.UpdateMetricHandler(storage))

	fmt.Printf("The server has just started on port [%s]\n", port)
	log.Fatal(router.Run(address))
}
