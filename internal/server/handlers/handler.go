package handlers

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/server/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UpdateMetricHandler(service *service.MetricService) func(c *gin.Context) {
	return func(c *gin.Context) {
		//if c.GetHeader("Content-Type") != "text/plain" {
		//	c.String(http.StatusBadRequest, "Invalid Content-Type")
		//}

		metricType := c.Param("type")
		if !(metricType == metrics.StringCounterType || metricType == metrics.StringGaugeType) {
			c.String(http.StatusNotImplemented, "Wrong metric type")
			return
		}
		metricName := c.Param("name")
		if metricName == "" {
			c.String(http.StatusNotFound, "No metric ID specified")
			return
		}
		metricValue := c.Param("value")
		if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
			c.String(http.StatusBadRequest, "Invalid metric value")
			return
		}

		service.UpdateMetric(metricType, metricName, metricValue)
	}
}

func GetMetricHandler(storage repository.Getter) func(c *gin.Context) {
	return func(c *gin.Context) {
		metricType := c.Param("type")
		if !(metricType == metrics.StringCounterType || metricType == metrics.StringGaugeType) {
			c.String(http.StatusNotImplemented, "Wrong metric type")
			return
		}
		metricName := c.Param("name")
		if metricName == "" {
			c.String(http.StatusNotFound, "No metric ID specified")
			return
		}

		value, err := storage.Get(metricType, metricName)
		if err != nil {
			if err == metrics.ErrNoValue {
				c.Status(http.StatusNotFound)
				return
			}
		}

		c.String(http.StatusOK, fmt.Sprintf("%v", value))
	}
}
