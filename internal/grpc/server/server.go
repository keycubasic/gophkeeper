// Package server contains gRPC server.
package server

import (
	"context"

	"github.com/keycubasic/gophkeeper/internal/app"
	pb "github.com/keycubasic/gophkeeper/internal/grpc/proto"

	"github.com/keycubasic/gophkeeper/internal/config"
)

// GrpcServer stores unimplemented server, configuration and functions of the application.
type GrpcServer struct {
	pb.UnimplementedGophKeeperServiceServer
	cfg *config.Config
	app app.Application
}

// NewGrpcServer creates a new GrpcServer object.
func NewGrpcServer(cfg *config.Config, app app.Application) *GrpcServer {
	return &GrpcServer{cfg: cfg, app: app}
}

// Ping checks the availability of the data store.
func (s *GrpcServer) Ping(ctx context.Context, _ *pb.PingRequest) (*pb.PingResponse, error) {
	if err := s.app.PingContext(ctx); err != nil {
		return nil, err
	}
	return &pb.PingResponse{}, nil
}
