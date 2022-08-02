package handlers

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func UpdateMetricHandler(storer repository.Storer) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		//if r.Header.Get("Content-Type") != "text/plain" {
		//	rw.WriteHeader(http.StatusBadRequest)
		//	rw.Write([]byte("Invalid Content-Type"))
		//	return
		//}

		metricType := chi.URLParam(r, "type")
		if !(metricType == metrics.StringCounterType || metricType == metrics.StringGaugeType) {
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte("Wrong metric type"))
			return
		}
		metricName := chi.URLParam(r, "name")
		if metricName == "" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("No metric ID specified"))
			return
		}
		metricValue := chi.URLParam(r, "value")
		if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Invalid metric value"))
			return
		}

		fmt.Printf("Got updating metric request: ")
		fmt.Printf("[%s] \"%s\": %v\n", metricType, metricName, metricValue)
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
	}
}

func GetMetricHandler(getter repository.Getter) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		if !(metricType == metrics.StringCounterType || metricType == metrics.StringGaugeType) {
			rw.WriteHeader(http.StatusNotImplemented)
			rw.Write([]byte("Wrong metric type"))
			return
		}
		metricName := chi.URLParam(r, "name")
		if metricName == "" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("No metric ID specified"))
			return
		}

		value, err := getter.Get(metricType, metricName)
		if err != nil {
			if err == metrics.ErrNoValue {
				rw.WriteHeader(http.StatusNotFound)
				return
			}
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("%v", value)))
	}
}

func GetAllMetricsHandler(storage repository.CollectionGetter) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	}
}
