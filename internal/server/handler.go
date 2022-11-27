package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (s *Server) PingDBHandler(rw http.ResponseWriter, _ *http.Request) {
	if err := s.storage.Ping(context.Background()); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (s *Server) UpdateMetricHandler(rw http.ResponseWriter, r *http.Request) {
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

	switch metricType {
	case metrics.StringCounterType:
		v, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			panic(err)
		}
		if err := s.storage.Store(*metrics.NewCounter(metricName, v)); err != nil {
			String(rw, http.StatusInternalServerError, err.Error())
			return
		}
	case metrics.StringGaugeType:
		v, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			panic(err)
		}
		if err := s.storage.Store(*metrics.NewGauge(metricName, v)); err != nil {
			String(rw, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

func (s *Server) GetMetricHandler(rw http.ResponseWriter, r *http.Request) {
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

	value, err := s.storage.Get(*metrics.New(metricType, metricName))
	if err != nil {
		if errors.Is(err, metrics.ErrNoValue) {
			String(rw, http.StatusNotFound, "")
			return
		}
	}

	switch value.Type {
	case metrics.StringCounterType:
		String(rw, http.StatusOK, fmt.Sprintf("%v", *value.Delta))
	case metrics.StringGaugeType:
		String(rw, http.StatusOK, fmt.Sprintf("%v", *value.Value))
	}
}
