package postgres

import (
	"context"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/keycubasic/gophkeeper/internal/config"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	glog "go.finelli.dev/gooseloggers/zerolog"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Storage stores the application configuration and database engine.
type Storage struct {
	cfg *config.Config
	db  *sqlx.DB
}

// NewLinks creates a new Links object
func NewStorage(ctx context.Context, cfg *config.Config) (*Storage, error) {
	log.Info().Msgf("Database DSN: %s", cfg.DatabaseDSN)
	db, err := sqlx.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("storage open connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("stogate ping connection: %w", err)
	}

	log.Info().Msg("Database connected...")
	return &Storage{cfg: cfg, db: db}, nil
}

// Migrate migrates database objects.
func (s Storage) Migrate(ctx context.Context) error {
	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(glog.GooseZerologLogger(&log.Logger))

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("storage migrate set dialect: %w", err)
	}

	if err := goose.Up(s.db.DB, "migrations"); err != nil {
		return fmt.Errorf("storage migrate up: %w", err)
	}

	return nil
}

// PingContext checks the database connection.
func (s Storage) PingContext(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// RegisterUser registers new user and returns user Id.
func (s Storage) RegisterUser(ctx context.Context, username, hashedPassword string) (int, error) {
	// TODO: implement RegisterUser

	return 0, nil
}

// FindUser finds user.
func (s Storage) FindUser(ctx context.Context, username string) (int, error) {
	// TODO: implement FindUser

	return 0, nil
}

// Close closes the connection to the database.
func (s Storage) Close() error {
	return s.db.Close()
}
