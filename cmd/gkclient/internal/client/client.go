// Package client contains auth and data clients description.
package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"fmt"

	apb "github.com/keycubasic/gophkeeper/internal/auth/proto"
	gpb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	//go:embed cert/ca-cert.pem
	pemServerCA []byte

	//go:embed cert/client-cert.pem
	certFile []byte

	//go:embed cert/client-key.pem
	keyFile []byte
)

// AuthClient stores auth client.
type AuthClient struct {
	Conn   *grpc.ClientConn
	Client apb.AuthServiceClient
}

// GophKeeperClient stores data client.
type GophKeeperClient struct {
	Conn   *grpc.ClientConn
	Client gpb.GophKeeperServiceClient
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Load client's certificate and private key
	clientCert, err := tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

// NewAuthClient creates new auth client.
func NewAuthClient(ctx context.Context, serverAddress string, enableTLS bool) *AuthClient {
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	if enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatal().Msgf("cannot load TLS credentials: %v", err)
		}
		transportOption = grpc.WithTransportCredentials(tlsCredentials)
	}

	conn, err := grpc.DialContext(ctx, serverAddress, transportOption)
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("Auth server %s connected", serverAddress)

	return &AuthClient{
		Conn:   conn,
		Client: apb.NewAuthServiceClient(conn),
	}
}

// NewGophKeeperClient creates new data client.
func NewGophKeeperClient(ctx context.Context, serverAddress string, enableTLS bool) *GophKeeperClient {
	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	if enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatal().Msgf("cannot load TLS credentials: %v", err)
		}
		transportOption = grpc.WithTransportCredentials(tlsCredentials)
	}
	conn, err := grpc.DialContext(ctx, serverAddress, transportOption)
	if err != nil {
		log.Fatal().Err(err)
	}

	log.Info().Msgf("GophKepeer server %s connected", serverAddress)

	return &GophKeeperClient{
		Conn:   conn,
		Client: gpb.NewGophKeeperServiceClient(conn),
	}
}
