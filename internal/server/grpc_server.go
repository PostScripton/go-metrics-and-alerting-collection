package server

import (
	"context"
	"net"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/factory/storage"
	pb "github.com/PostScripton/go-metrics-and-alerting-collection/internal/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	listener net.Listener
	core     *grpc.Server
	storage  storage.Storager
}

func newGRPCServer(address string, storage storage.Storager, key string) *GRPCServer {
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("Splitting host and port")
		return nil
	}

	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal().Err(err).Str("address", address).Msg("Unable to listen address")
	}

	core := grpc.NewServer()
	pb.RegisterMetricsServer(core, &MetricsServer{
		storage: storage,
		key:     key,
	})

	s := &GRPCServer{
		listener: listen,
		core:     core,
		storage:  storage,
	}

	return s
}

func (s *GRPCServer) Run() error {
	return s.core.Serve(s.listener)
}

func (s *GRPCServer) Shutdown(_ context.Context) error {
	s.core.GracefulStop()

	return nil
}
