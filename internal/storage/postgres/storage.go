// Package postgre contains functions for working with the PostgreSQL database.
package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/keycubasic/gophkeeper/internal/config"
	"github.com/keycubasic/gophkeeper/internal/model"
	"github.com/keycubasic/gophkeeper/internal/storage"
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

// RegisterUser registers new user and returns user ID.
func (s Storage) RegisterUser(ctx context.Context, username, hashedPassword string) (int, error) {
	sql := sq.Insert("users").Columns("username", "hash").Values(username, hashedPassword).
		Suffix("RETURNING user_id").RunWith(s.db).PlaceholderFormat(sq.Dollar)

	var userID int
	err := sql.QueryRowContext(ctx).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return 0, fmt.Errorf("storage register user: %w", storage.ErrUniqueViolation)
		}

		return 0, fmt.Errorf("storage register user: %w", err)
	}

	return userID, nil
}

// FindUser finds user.
func (s Storage) FindUser(ctx context.Context, username string) (*model.User, error) {
	sql := sq.Select("*").From("users").Where(sq.Eq{"username": username})
	rows, err := sql.RunWith(s.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage find user query: %w", err)
	}

	user := new(model.User)
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.HashedPassword)
		if err != nil {
			return nil, fmt.Errorf("storage scan user row: %w", err)
		}
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("storage find user rows: %w", err)
	}

	if user.ID == 0 {
		return nil, fmt.Errorf("storage find user: %w", errors.New("user not found"))
	}

	return user, nil
}

// StoreData stores data to database.
func (s Storage) StoreData(ctx context.Context, creds []model.Creds) error {
	sql := `INSERT INTO creds (user_id, user_site, user_login, user_password, metadata)
	 		VALUES (:user_id, :user_site, :user_login, :user_password, :metadata)`
	if _, err := s.db.NamedExecContext(ctx, sql, creds); err != nil {
		return fmt.Errorf("stogare store data: %w", err)
	}

	return nil
}

// GetData gets data from database.
func (s Storage) GetData(ctx context.Context, userID int, site string) (*model.Creds, error) {
	sql := sq.Select("user_login, user_password, metadata").
		From("creds").Where(sq.Eq{"user_id": userID}, sq.Eq{"site": site})
	rows, err := sql.RunWith(s.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage get data query: %w", err)
	}

	creds := new(model.Creds)
	for rows.Next() {
		err := rows.Scan(&creds.Login, &creds.Password, &creds.Metadata)
		if err != nil {
			return nil, fmt.Errorf("storage get data scan row: %w", err)
		}
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("storage get data rows: %w", err)
	}

	return creds, nil
}

// Close closes the connection to the database.
func (s Storage) Close() error {
	return s.db.Close()
}
