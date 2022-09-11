package server

import (
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/hashing"
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

	if s.key != "" {
		sign := hashing.HashMetric(&metricsRequest, s.key)
		if !hashing.ValidHash(sign, metricsRequest.Hash) {
			JSON(rw, http.StatusBadRequest, JSONObj{"message": "Signature does not match"})
			return
		}
	}

	if err := s.storage.Store(metricsRequest); err != nil {
		JSON(rw, http.StatusInternalServerError, JSONObj{"message": fmt.Sprintf("Error on storing data: %s", err)})
		return
	}

	JSON(rw, http.StatusOK, JSONObj{})

	if metricsRequest.Delta != nil {
		fmt.Printf("Metric updated! [%s] \"%s\": %d\n", metricsRequest.Type, metricsRequest.ID, *metricsRequest.Delta)
	} else if metricsRequest.Value != nil {
		fmt.Printf("Metric updated! [%s] \"%s\": %f\n", metricsRequest.Type, metricsRequest.ID, *metricsRequest.Value)
	}
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

	if s.key != "" {
		value.Hash = hashing.HashToHexMetric(value, s.key)
	}

	JSON(rw, http.StatusOK, value)
}

func (s *server) UpdateMetricsBatchJSONHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		JSON(rw, http.StatusBadRequest, JSONObj{"message": "Invalid Content-Type"})
		return
	}
	var metricsCollection []metrics.Metrics
	if err := ParseJSON(r, &metricsCollection); err != nil {
		JSON(rw, http.StatusInternalServerError, JSONObj{"message": "Unable to parse JSON"})
		return
	}

	if s.key != "" {
		for _, m := range metricsCollection {
			sign := hashing.HashMetric(&m, s.key)
			if !hashing.ValidHash(sign, m.Hash) {
				JSON(rw, http.StatusBadRequest, JSONObj{"message": fmt.Sprintf("Signature for [%s] does not match", m.ID)})
				return
			}
		}
	}

	var metricsMap = make(map[string]metrics.Metrics)
	for _, m := range metricsCollection {
		if old, ok := metricsMap[m.ID]; ok {
			metrics.Update(&old, &m)
			metricsMap[m.ID] = old
		} else {
			metricsMap[m.ID] = m
		}
	}

	if err := s.storage.StoreCollection(metricsMap); err != nil {
		JSON(rw, http.StatusInternalServerError, JSONObj{"message": fmt.Sprintf("Error on storing data: %s", err)})
		return
	}

	JSON(rw, http.StatusOK, JSONObj{})
}
