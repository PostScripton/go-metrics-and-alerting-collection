package server

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"net/http"
)

type JSONObj map[string]any

var notFoundResponse = JSONObj{"message": "404 page not found"}
var methodNotAllowed = JSONObj{"message": "405 method not allowed"}

func (s *server) UpdateMetricJSONHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		JSON(rw, http.StatusBadRequest, JSONObj{"message": "Invalid Content-Type"})
		return
	}

	var metricsRequest metrics.Metrics
	if err := ParseJSON(r, &metricsRequest); err != nil {
		JSON(rw, http.StatusInternalServerError, JSONObj{"message": "Unable to parse JSON"})
		return
	}

	if metricsRequest.ID == "" {
		JSON(rw, http.StatusNotFound, JSONObj{"message": "No metric ID specified"})
		return
	}

	switch metricsRequest.Type {
	case metrics.StringCounterType:
		if metricsRequest.Delta == nil {
			JSON(rw, http.StatusNotFound, notFoundResponse)
			return
		}
	case metrics.StringGaugeType:
		if metricsRequest.Value == nil {
			JSON(rw, http.StatusNotFound, notFoundResponse)
			return
		}
	default:
		JSON(rw, http.StatusNotImplemented, JSONObj{"message": "Invalid metric type"})
		return
	}

	if err := s.storage.Store(metricsRequest); err != nil {
		JSON(rw, http.StatusInternalServerError, JSONObj{"message": err.Error()})
		return
	}

	JSON(rw, http.StatusOK, JSONObj{})

	var value interface{}
	if metricsRequest.Delta != nil {
		value = *metricsRequest.Delta
	} else if metricsRequest.Value != nil {
		value = *metricsRequest.Value
	}

	fmt.Printf("Metric updated! [%s] \"%s\": %v\n", metricsRequest.Type, metricsRequest.ID, value)
}

func (s *server) GetMetricJSONHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		JSON(rw, http.StatusBadRequest, JSONObj{"message": "Invalid Content-Type"})
		return
	}
	var metricsReq metrics.Metrics
	if err := ParseJSON(r, &metricsReq); err != nil {
		JSON(rw, http.StatusInternalServerError, JSONObj{"message": "Unable to parse JSON"})
		return
	}

	if metricsReq.ID == "" {
		JSON(rw, http.StatusNotFound, JSONObj{"message": "No metric ID specified"})
		return
	}
	switch metricsReq.Type {
	case metrics.StringCounterType:
	case metrics.StringGaugeType:
	default:
		JSON(rw, http.StatusNotImplemented, JSONObj{"message": "Invalid metric type"})
		return
	}

	value, err := s.storage.Get(metricsReq)
	if err != nil {
		if err == metrics.ErrNoValue {
			JSON(rw, http.StatusNotFound, JSONObj{"message": "No value"})
			return
		}
		JSON(rw, http.StatusInternalServerError, JSONObj{"message": err.Error()})
		return
	}

	JSON(rw, http.StatusOK, value)
}
