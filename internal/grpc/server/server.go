// Package server contains gRPC server.
package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/keycubasic/gophkeeper/internal/app"
	pb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
	"github.com/keycubasic/gophkeeper/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		return nil, status.Errorf(codes.Internal, "server ping: %v", err)
	}
	return &pb.PingResponse{}, nil
}

// StoreData stores data.
func (s *GrpcServer) StoreData(ctx context.Context, r *pb.StoreDataRequest) (*pb.StoreDataResponse, error) {
	userID, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "server store data: %v", err)
	}

	creds := model.Creds{
		UserID:   userID,
		Site:     r.Site,
		Login:    r.Credentials.Login,
		Password: r.Credentials.Password,
		Metadata: r.Metadata,
	}

	if err := s.app.AppendData(ctx, creds); err != nil {
		return nil, status.Errorf(codes.Internal, "server store data: %v", err)
	}

	return &pb.StoreDataResponse{Info: "user data will be stored"}, nil
}

// GetData gets data.
func (s *GrpcServer) GetData(ctx context.Context, r *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	userID, err := s.getUserFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "server get data: %v", err)
	}

	creds, err := s.app.GetData(ctx, userID, r.Site)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "server get data: %v", err)
	}

	pbCreds := &pb.Credentials{
		Login:    creds.Login,
		Password: creds.Password,
	}

	response := &pb.GetDataResponse{
		Credentials: pbCreds,
		Metadata:    creds.Metadata,
	}

	return response, nil
}

func (s *GrpcServer) getUserFromContext(ctx context.Context) (int, error) {
	username := ctx.Value("username")
	if username == nil {
		return 0, fmt.Errorf("server get user from context: %w", errors.New("empty user in context"))
	}

	user, err := s.app.FindUser(ctx, username.(string))
	if err != nil {
		return 0, fmt.Errorf("server get user from context: %w", err)
	}

	return user.ID, nil
}
