package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/keycubasic/gophkeeper/internal/app"
	"github.com/keycubasic/gophkeeper/internal/config"
	"github.com/keycubasic/gophkeeper/internal/server"
	"github.com/keycubasic/gophkeeper/internal/storage/postgres"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	cfg := new(config.Config)
	cfg.ParseConfig()

	db, err := postgres.NewStorage(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal().Err(err)
		}
		log.Info().Msg("Database connection closed")
	}()

	err = db.Migrate(ctx)
	if err != nil {
		log.Fatal().Err(err)
	}

	app := app.NewApp(ctx, cfg, db)
	gsrv := server.RegisterGrpcServer(ctx, cfg, app)
	server.RunServer(ctx, cfg, gsrv)
}
