package handlers

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

var notFoundResponse = gin.H{"message": "404 page not found"}

func UpdateMetricJSONHandler(storer repository.Storer) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Content-Type"})
		}

		var metricsRequest metrics.Metrics
		if err := c.BindJSON(&metricsRequest); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to parse JSON"})
		}

		if metricsRequest.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": "No metric ID specified"})
			return
		}

		switch metricsRequest.Type {
		case metrics.StringCounterType:
			if metricsRequest.Delta == nil {
				c.JSON(http.StatusNotFound, notFoundResponse)
				return
			}
			storer.Store(metricsRequest.ID, metrics.Counter(*metricsRequest.Delta))
		case metrics.StringGaugeType:
			if metricsRequest.Value == nil {
				c.JSON(http.StatusNotFound, notFoundResponse)
				return
			}
			storer.Store(metricsRequest.ID, metrics.Gauge(*metricsRequest.Value))
		default:
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Invalid metric type"})
			return
		}

		c.JSON(http.StatusOK, gin.H{})

		fmt.Printf("Metric updated! [%s] \"%s\" (", metricsRequest.Type, metricsRequest.ID)
		switch metricsRequest.Type {
		case metrics.StringCounterType:
			fmt.Print(*metricsRequest.Delta)
		case metrics.StringGaugeType:
			fmt.Print(*metricsRequest.Value)
		}
		fmt.Printf(")\n")
	}
}

func GetMetricJSONHandler(getter repository.Getter) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.GetHeader("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Content-Type"})
		}
		var metricsReq metrics.Metrics
		if err := c.BindJSON(&metricsReq); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to parse JSON"})
		}

		if metricsReq.ID == "" {
			c.JSON(http.StatusNotFound, gin.H{"message": "No metric ID specified"})
			return
		}
		switch metricsReq.Type {
		case metrics.StringCounterType:
		case metrics.StringGaugeType:
		default:
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Invalid metric type"})
			return
		}

		value, err := getter.Get(metricsReq.Type, metricsReq.ID)
		if err != nil {
			if err == metrics.ErrNoValue {
				c.JSON(http.StatusNotFound, notFoundResponse)
				return
			}
		}

		metricsRes := metrics.Metrics{
			ID:   metricsReq.ID,
			Type: metricsReq.Type,
		}
		switch metricsRes.Type {
		case metrics.StringCounterType:
			counter := value.(metrics.MetricIntCaster).ToInt64()
			metricsRes.Delta = &counter
		case metrics.StringGaugeType:
			gauge := value.(metrics.MetricFloatCaster).ToFloat64()
			metricsRes.Value = &gauge
		}
		c.JSON(200, metricsRes)
	}
}
