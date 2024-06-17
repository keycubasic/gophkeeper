package client

import (
	"context"

	apb "github.com/keycubasic/gophkeeper/internal/auth/proto"
	gpb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Client apb.AuthServiceClient
}

type GophKeeperClient struct {
	Client gpb.GophKeeperServiceClient
}

func NewAuthClient(ctx context.Context, serverAddress string) *AuthClient {
	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err)
	}
	defer conn.Close()

	log.Info().Msgf("Auth server %s connected", serverAddress)

	return &AuthClient{
		Client: apb.NewAuthServiceClient(conn),
	}
}

func NewGophKeeperClient(ctx context.Context, serverAddress string) *GophKeeperClient {
	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err)
	}
	defer conn.Close()

	log.Info().Msgf("GophKepeer server %s connected", serverAddress)

	return &GophKeeperClient{
		Client: gpb.NewGophKeeperServiceClient(conn),
	}
}
