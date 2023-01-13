package client

import (
	"context"
	"net/url"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	pb "github.com/PostScripton/go-metrics-and-alerting-collection/internal/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client pb.MetricsClient
}

func newGRPCClient(address string) *GRPCClient {
	uri, err := url.Parse(address)
	if err != nil {
		log.Fatal().Err(err).Str("address", address).Msg("Parsing address")
		return nil
	}

	conn, err := grpc.Dial(":"+uri.Port(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("Dealing with grpc")
		return nil
	}

	return &GRPCClient{
		conn:   conn,
		client: pb.NewMetricsClient(conn),
	}
}

func (c *GRPCClient) UpdateMetric(metric metrics.Metrics) error {
	_, err := c.client.UpdateMetric(context.Background(), &pb.MetricRequest{
		Metric: toGRPCMetric(metric),
	})

	return err
}

func (c *GRPCClient) BatchUpdateMetrics(collection map[string]metrics.Metrics) error {
	_, err := c.client.BatchUpdateMetrics(context.Background(), &pb.BatchMetricsRequest{
		Metrics: toGRPCMetricsCollection(collection),
	})

	return err
}

// Close закрывает gRPC соединение
func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

// toGRPCMetric переводит metrics.Metrics в proto.Metric для передачи по gRPC
func toGRPCMetric(metric metrics.Metrics) *pb.Metric {
	return &pb.Metric{
		ID:    metric.ID,
		Type:  metric.Type,
		Delta: metric.Delta,
		Value: metric.Value,
		Hash:  metric.Hash,
	}
}

// toGRPCMetricsCollection переводит мапу metrics.Metrics в слайc proto.Metric для передачи по gRPC
func toGRPCMetricsCollection(collection map[string]metrics.Metrics) []*pb.Metric {
	result := make([]*pb.Metric, 0, len(collection))

	for _, metric := range collection {
		result = append(result, toGRPCMetric(metric))
	}

	return result
}
