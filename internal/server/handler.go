package server

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (s *server) UpdateMetricHandler(rw http.ResponseWriter, r *http.Request) {
	//if r.Header.Get("Content-Type") != "text/plain" {
	//	String(rw, http.StatusBadRequest, "Invalid Content-Type")
	//	return
	//}

	metricType := chi.URLParam(r, "type")
	if !(metricType == metrics.StringCounterType || metricType == metrics.StringGaugeType) {
		String(rw, http.StatusNotImplemented, "Wrong metric type")
		return
	}
	metricName := chi.URLParam(r, "name")
	if metricName == "" {
		String(rw, http.StatusNotFound, "No metric ID specified")
		return
	}
	metricValue := chi.URLParam(r, "value")
	if _, err := strconv.ParseFloat(metricValue, 64); err != nil {
		String(rw, http.StatusBadRequest, "Invalid metric value")
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
		s.storage.Store(metricName, metrics.Counter(v))
	case metrics.StringGaugeType:
		v, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			panic(err)
		}
		s.storage.Store(metricName, metrics.Gauge(v))
	}
}

func (s *server) GetMetricHandler(rw http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	if !(metricType == metrics.StringCounterType || metricType == metrics.StringGaugeType) {
		String(rw, http.StatusNotImplemented, "Wrong metric type")
		return
	}
	metricName := chi.URLParam(r, "name")
	if metricName == "" {
		String(rw, http.StatusNotFound, "No metric ID specified")
		return
	}

	value, err := s.storage.Get(metricType, metricName)
	if err != nil {
		if err == metrics.ErrNoValue {
			String(rw, http.StatusNotFound, "")
			return
		}
	}

	String(rw, http.StatusOK, fmt.Sprintf("%v", value))
}