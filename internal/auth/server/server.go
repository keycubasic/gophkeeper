// Package server contains auth gRPC server.
package server

import (
	"context"
	"time"

	"github.com/keycubasic/gophkeeper/internal/app"
	"github.com/keycubasic/gophkeeper/internal/auth"
	pb "github.com/keycubasic/gophkeeper/internal/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/keycubasic/gophkeeper/internal/config"
)

// AuthServer stores unimplemented server, configuration and functions of the application.
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	cfg        *config.Config
	app        app.Application
	jwtManager *auth.JWTManager
}

// NewAuthServer creates a new AuthServer object.
func NewAuthServer(cfg *config.Config, app app.Application) *AuthServer {
	return &AuthServer{
		cfg: cfg,
		app: app,
		jwtManager: auth.NewJWTManager(
			cfg.JWTSecretKey,
			time.Duration(cfg.TokenDuration)*time.Hour,
		),
	}
}

// RegisterUser registers new user.
func (server *AuthServer) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	_, err := server.app.FindUser(ctx, req.User.UserName)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "auth server username %s is already in use by another user: %v", req.User.UserName, err)
	}

	user, err := server.app.RegisterUser(ctx, req.User.UserName, req.User.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "auth server register user: %v", err)
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "auth server cannot generate access token")
	}

	res := &pb.RegisterUserResponse{AccessToken: token}
	return res, nil
}

// LoginUser find and login user.
func (server *AuthServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.app.FindUser(ctx, req.User.UserName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "auth server cannot find user: %v", err)
	}

	if user == nil || !server.app.IsCorrectPassword(user, req.User.Password) {
		return nil, status.Errorf(codes.NotFound, "auth server incorrect username/password")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "auth server cannot generate access token")
	}

	res := &pb.LoginUserResponse{AccessToken: token}
	return res, nil
}
