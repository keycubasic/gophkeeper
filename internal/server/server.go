// Package server contains logging interceptor, registration of servers and starts the main server.
package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"embed"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/keycubasic/gophkeeper/internal/app"
	"github.com/keycubasic/gophkeeper/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/swaggest/swgui/v5emb"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/keycubasic/gophkeeper/internal/auth"
	apb "github.com/keycubasic/gophkeeper/internal/auth/proto"
	asrv "github.com/keycubasic/gophkeeper/internal/auth/server"
	gpb "github.com/keycubasic/gophkeeper/internal/grpc/proto"
	gsrv "github.com/keycubasic/gophkeeper/internal/grpc/server"
)

var (
	//go:embed swagger/*
	swaggerJSONPath embed.FS

	//go:embed hosts/whitelist
	embedWhitelist string

	//go:embed cert/ca-cert.pem
	clientCACert []byte

	//go:embed cert/server-cert.pem
	certFile []byte

	//go:embed cert/server-key.pem
	keyFile []byte
)

// httpGrpcRouter separates HTTP and gRPC calls.
func httpGrpcRouter(gsrv *grpc.Server, hsvr http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor != 2 {
			hsvr.ServeHTTP(w, r)
		}

		if strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			gsrv.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

// interceptorLogger adapts zerolog logger to interceptor logger.
func interceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		l := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		}
	})
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(clientCACert) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}

	// Load server's certificate and private key
	serverCert, err := tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}

// RegisterGrpcServer registers the gRPC server.
func RegisterGrpcServer(ctx context.Context, cfg *config.Config, app app.Application) *grpc.Server {
	authServer := asrv.NewAuthServer(cfg, app)
	grpcServer := gsrv.NewGrpcServer(cfg, app)

	logger := zerolog.New(os.Stderr)
	opts := []logging.Option{logging.WithTimestampFormat(time.RFC3339)}
	authInterceptor := auth.NewAuthInterceptor(authServer.JwtManager)

	serverOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(interceptorLogger(logger), opts...),
			grpc.UnaryServerInterceptor(authInterceptor.Unary()),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(interceptorLogger(logger), opts...),
			grpc.StreamServerInterceptor(authInterceptor.Stream()),
		),
	}

	if cfg.EnableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			log.Fatal().Msgf("cannot load TLS credentials: %v", err)
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}

	server := grpc.NewServer(serverOptions...)
	apb.RegisterAuthServiceServer(server, authServer)
	gpb.RegisterGophKeeperServiceServer(server, grpcServer)
	reflection.Register(server)
	log.Info().Msg("gRPC server registered")
	return server
}

// RunServer starts the server.
func RunServer(ctx context.Context, cfg *config.Config, gsrv *grpc.Server) {
	addr := cfg.ServerAddress
	var s http.Server
	g, gCtx := errgroup.WithContext(ctx)
	gsrvmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := apb.RegisterAuthServiceHandlerFromEndpoint(ctx, gsrvmux, cfg.ServerAddress, opts); err != nil {
		log.Fatal().Err(err)
	}
	if err := gpb.RegisterGophKeeperServiceHandlerFromEndpoint(ctx, gsrvmux, cfg.ServerAddress, opts); err != nil {
		log.Fatal().Err(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", gsrvmux.ServeHTTP)

	// Swagger
	mux.Handle("/api/docs/", http.StripPrefix("/api/docs/", http.FileServer(http.FS(swaggerJSONPath))))
	mux.Handle("/swagger/", v5emb.New("GophKeeper", "/api/docs/swagger/swagger.json", "/swagger/"))

	hgr := httpGrpcRouter(gsrv, mux)
	if cfg.EnableHTTPS {
		log.Info().Msg("HTTPS enabled")
		manager := &autocert.Manager{
			Cache:      autocert.DirCache("cache-dir"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(strings.Split(embedWhitelist, "\n")...),
		}
		s = http.Server{Addr: addr, Handler: hgr, TLSConfig: manager.TLSConfig()}
		g.Go(func() error {
			return s.ListenAndServeTLS("", "")
		})
	} else {
		log.Info().Msg("HTTPS disabled")
		s = http.Server{Addr: addr, Handler: hgr}
		g.Go(func() error {
			return s.ListenAndServe()
		})
	}

	log.Info().Msgf("Server listening address %s", addr)

	g.Go(func() error {
		<-gCtx.Done()
		log.Info().Msg("Server stops")
		ctxWithTimeout, cancel := context.
			WithTimeout(context.Background(), time.Duration(10)*time.Second)
		defer cancel()

		gsrv.GracefulStop()
		log.Info().Msg("gRPC server stops")

		return s.Shutdown(ctxWithTimeout)
	})

	if err := g.Wait(); err != nil {
		log.Info().Err(err).Msg("Exit reason")
	}
}
