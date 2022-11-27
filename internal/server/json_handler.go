package server

import (
	"errors"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/PostScripton/go-metrics-and-alerting-collection/pkg/hashing/hmac"
	"github.com/rs/zerolog/log"
	"net/http"
)

type JSONObj map[string]any

var notFoundResponse = JSONObj{"message": "404 page not found"}
var methodNotAllowed = JSONObj{"message": "405 method not allowed"}

func (s *Server) UpdateMetricJSONHandler(rw http.ResponseWriter, r *http.Request) {
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
		if !metricsRequest.ValidHash(hmac.NewHmacSigner(), metricsRequest.Hash, s.key) {
			JSON(rw, http.StatusBadRequest, JSONObj{"message": "Signature does not match"})
			return
		}
	}

	if err := s.storage.Store(metricsRequest); err != nil {
		JSON(rw, http.StatusInternalServerError, JSONObj{"message": fmt.Sprintf("Error on storing data: %s", err)})
		return
	}

	JSON(rw, http.StatusOK, JSONObj{})

	log.Debug().Interface("metric", metricsRequest).Msg("Metric updated!")
}

func (s *Server) GetMetricJSONHandler(rw http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, metrics.ErrNoValue) {
			JSON(rw, http.StatusNotFound, JSONObj{"message": "No value"})
			return
		}
		JSON(rw, http.StatusInternalServerError, JSONObj{"message": err.Error()})
		return
	}

	if s.key != "" {
		value.Hash = value.ToHexHash(hmac.NewHmacSigner(), s.key)
	}

	JSON(rw, http.StatusOK, value)
}

func (s *Server) UpdateMetricsBatchJSONHandler(rw http.ResponseWriter, r *http.Request) {
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
			if !m.ValidHash(hmac.NewHmacSigner(), m.Hash, s.key) {
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
		log.Debug().Interface("metric", m).Msg("Metric of collection updated!")
	}

	if err := s.storage.StoreCollection(metricsMap); err != nil {
		JSON(rw, http.StatusInternalServerError, JSONObj{"message": fmt.Sprintf("Error on storing data: %s", err)})
		return
	}

	log.Info().Msg("Metrics collection updated")

	JSON(rw, http.StatusOK, JSONObj{})
}
