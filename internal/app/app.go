// Package app contains the basic functions for the application.
package app

import (
	"context"
	"fmt"

	"github.com/keycubasic/gophkeeper/internal/config"
	"github.com/keycubasic/gophkeeper/internal/model"
	"github.com/keycubasic/gophkeeper/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

// Application interface.
type Application interface {
	// Ping checks the availability of the data store.
	PingContext(ctx context.Context) error
	// RegisterUser registers new user and returns user Id.
	RegisterUser(ctx context.Context, username, password string) (*model.User, error)
	// IsCorrectPassword checks password.
	IsCorrectPassword(user *model.User, password string) bool
	// FindUser finds user.
	FindUser(ctx context.Context, username string) (*model.User, error)
}

// App stores an interface for application configuration and data storage.
type App struct {
	cfg *config.Config
	db  storage.Database
}

// NewApp creates a new App object.
func NewApp(ctx context.Context, cfg *config.Config, db storage.Database) *App {
	return &App{
		cfg: cfg,
		db:  db,
	}
}

// Ping checks the availability of the data store.
func (app *App) PingContext(ctx context.Context) error {
	return app.db.PingContext(ctx)
}

// RegisterUser registers new user.
func (app *App) RegisterUser(ctx context.Context, username string, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("app register user hash password: %w", err)
	}

	userID, err := app.db.RegisterUser(ctx, username, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("app register user: %w", err)
	}

	user := &model.User{
		ID:             userID,
		Username:       username,
		HashedPassword: string(hashedPassword),
	}

	return user, nil
}

// IsCorrectPassword checks password.
func (app *App) IsCorrectPassword(user *model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

// FindUser finds user.
func (app *App) FindUser(ctx context.Context, username string) (*model.User, error) {

	// TODO: implement FindUser

	return nil, nil
}
