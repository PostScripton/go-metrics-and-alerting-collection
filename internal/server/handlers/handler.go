package handlers

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NoRoute(c *gin.Context) {
	if c.ContentType() == "application/json" {
		c.JSON(404, notFoundResponse)
	} else {
		c.String(404, "404 page not found")
	}
}

func UpdateMetricHandler(storer repository.Storer) func(c *gin.Context) {
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

		switch metricType {
		case metrics.StringCounterType:
			v, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				panic(err)
			}
			storer.Store(metricName, metrics.Counter(v))
		case metrics.StringGaugeType:
			v, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				panic(err)
			}
			storer.Store(metricName, metrics.Gauge(v))
		}

		c.Status(http.StatusOK)
	}
}

func GetMetricHandler(getter repository.Getter) func(c *gin.Context) {
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

		value, err := getter.Get(metricType, metricName)
		if err != nil {
			if err == metrics.ErrNoValue {
				c.Status(http.StatusNotFound)
				return
			}
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("%v", value))
	}
}

func GetAllMetricsHandler(storage repository.CollectionGetter) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	}
}
