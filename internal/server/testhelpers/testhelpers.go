// Package testhelpers provides help function for testing.
package testhelpers

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apb "github.com/keycubasic/gophkeeper/internal/auth/proto"
	gpb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
)

// PostgresContainer stores postgres container and database connection string.
type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
}

// AuthClient stores client for auth server.
type AuthClient struct {
	Conn   *grpc.ClientConn
	Client apb.AuthServiceClient
}

// GophKeeperClient stores client for data server.
type GophKeeperClient struct {
	Conn   *grpc.ClientConn
	Client gpb.GophKeeperServiceClient
}

// CreatePostgresContainer creates postgres container.
func CreatePostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2-alpine"),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return nil, err
	}
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &PostgresContainer{
		PostgresContainer: pgContainer,
		ConnectionString:  connStr,
	}, nil
}

// NewAuthClient creates auth client.
func NewAuthClient(ctx context.Context, serverAddress string) (*AuthClient, error) {
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.DialContext(ctx, serverAddress, transportOption)
	if err != nil {
		return nil, err
	}

	return &AuthClient{
		Conn:   conn,
		Client: apb.NewAuthServiceClient(conn),
	}, nil
}

// NewGophKeeperClient creates data client.
func NewGophKeeperClient(ctx context.Context, serverAddress string) (*GophKeeperClient, error) {
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.DialContext(ctx, serverAddress, transportOption)
	if err != nil {
		return nil, err
	}

	return &GophKeeperClient{
		Conn:   conn,
		Client: gpb.NewGophKeeperServiceClient(conn),
	}, nil
}
