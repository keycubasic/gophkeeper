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
	Conn   *grpc.ClientConn
	Client apb.AuthServiceClient
}

type GophKeeperClient struct {
	Conn   *grpc.ClientConn
	Client gpb.GophKeeperServiceClient
}

func NewAuthClient(ctx context.Context, serverAddress string) *AuthClient {
	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("Auth server %s connected", serverAddress)

	return &AuthClient{
		Conn:   conn,
		Client: apb.NewAuthServiceClient(conn),
	}
}

func NewGophKeeperClient(ctx context.Context, serverAddress string) *GophKeeperClient {
	conn, err := grpc.DialContext(ctx, serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("GophKepeer server %s connected", serverAddress)

	return &GophKeeperClient{
		Conn:   conn,
		Client: gpb.NewGophKeeperServiceClient(conn),
	}
}
