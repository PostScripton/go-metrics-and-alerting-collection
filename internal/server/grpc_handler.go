package server

import (
	"context"
	"fmt"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	pb "github.com/PostScripton/go-metrics-and-alerting-collection/internal/proto"
	"github.com/PostScripton/go-metrics-and-alerting-collection/pkg/hashing/hmac"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MetricsServer struct {
	pb.UnimplementedMetricsServer

	storage storage.Storager
	key     string
}

var _ pb.MetricsServer = (*MetricsServer)(nil)

func (s *MetricsServer) UpdateMetric(_ context.Context, in *pb.MetricRequest) (*emptypb.Empty, error) {
	if in.Metric.GetID() == "" {
		return nil, status.Error(codes.InvalidArgument, "No metric ID specified")
	}

	switch in.Metric.GetType() {
	case metrics.StringCounterType:
		if in.Metric.Delta == nil {
			return nil, status.Error(codes.InvalidArgument, "No delta passed")
		}
	case metrics.StringGaugeType:
		if in.Metric.Value == nil {
			return nil, status.Error(codes.InvalidArgument, "No value passed")
		}
	default:
		return nil, status.Error(codes.Unimplemented, "Invalid metric type")
	}

	internalMetric := toInternalMetrics(in.Metric)

	if s.key != "" {
		if !internalMetric.ValidHash(hmac.NewHmacSigner(), internalMetric.Hash, s.key) {
			return nil, status.Error(codes.InvalidArgument, "Signature does not match")
		}
	}

	if err := s.storage.Store(internalMetric); err != nil {
		return nil, status.Errorf(codes.Unknown, fmt.Sprintf("Error on storing data: %s", err))
	}

	log.Debug().Interface("metric", internalMetric).Msg("Metric updated!")

	return &emptypb.Empty{}, nil
}

func (s *MetricsServer) BatchUpdateMetrics(_ context.Context, in *pb.BatchMetricsRequest) (*emptypb.Empty, error) {
	internalMetrics := toInternalMetricsCollection(in.Metrics)

	if s.key != "" {
		for _, m := range internalMetrics {
			if !m.ValidHash(hmac.NewHmacSigner(), m.Hash, s.key) {
				return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Signature for [%s] does not match", m.ID))
			}
		}
	}

	var metricsMap = make(map[string]metrics.Metrics)
	for _, m := range internalMetrics {
		if old, ok := metricsMap[m.ID]; ok {
			metrics.Update(&old, &m)
			metricsMap[m.ID] = old
		} else {
			metricsMap[m.ID] = m
		}
		log.Debug().Interface("metric", m).Msg("Metric of collection updated!")
	}

	if err := s.storage.StoreCollection(metricsMap); err != nil {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Error on storing data: %s", err))
	}

	log.Info().Msg("Metrics collection updated")

	return &emptypb.Empty{}, nil
}

func toInternalMetrics(metric *pb.Metric) metrics.Metrics {
	return metrics.Metrics{
		ID:    metric.ID,
		Type:  metric.Type,
		Delta: metric.Delta,
		Value: metric.Value,
		Hash:  metric.Hash,
	}
}

func toInternalMetricsCollection(grpcMetrics []*pb.Metric) []metrics.Metrics {
	collection := make([]metrics.Metrics, 0, len(grpcMetrics))

	for _, metric := range grpcMetrics {
		collection = append(collection, toInternalMetrics(metric))
	}

	return collection
}
