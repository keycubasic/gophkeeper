// Package app contains the basic functions for the application.
package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/keycubasic/gophkeeper/internal/config"
	"github.com/keycubasic/gophkeeper/internal/crypto"
	"github.com/keycubasic/gophkeeper/internal/model"
	"github.com/keycubasic/gophkeeper/internal/storage"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// Application interface.
type Application interface {
	// Ping checks the availability of the data store.
	PingContext(ctx context.Context) error
	// RegisterUser registers new user and returns user Id.
	RegisterUser(ctx context.Context, username, password string) (*model.User, error)
	// IsCorrectPassword checks password.
	IsCorrectPassword(hashedPassword, password string) bool
	// FindUser finds user.
	FindUser(ctx context.Context, username string) (*model.User, error)
	// AppendData stores data.
	AppendData(ctx context.Context, creds model.Creds) error
	// GetData gets data.
	GetData(ctx context.Context, userID int, site string) (*model.Creds, error)
}

// App stores an interface for application configuration and data storage.
type App struct {
	cfg       *config.Config
	db        storage.Database
	credsChan chan model.Creds
}

// NewApp creates a new App object.
func NewApp(ctx context.Context, cfg *config.Config, db storage.Database) *App {
	instance := &App{
		cfg:       cfg,
		db:        db,
		credsChan: make(chan model.Creds, 256),
	}

	go instance.processData(ctx)

	return instance
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
func (app *App) IsCorrectPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// FindUser finds user.
func (app *App) FindUser(ctx context.Context, username string) (*model.User, error) {
	return app.db.FindUser(ctx, username)
}

// AppendData stores data.
func (app *App) AppendData(ctx context.Context, creds model.Creds) (err error) {
	creds.Password, err = crypto.Encode(app.cfg.SecretKey, creds.Password)
	if err != nil {
		return
	}
	app.credsChan <- creds
	return nil
}

func (app *App) processData(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	creds := make([]model.Creds, 0, 256)

	for {
		select {
		case <-ctx.Done():
			ctxClosing, closingCancel := context.
				WithTimeout(context.Background(), 10*time.Second)
			defer closingCancel()

			if err := app.PingContext(ctxClosing); err == nil {
				app.storeData(ctxClosing, creds)
			}

			log.Info().Msg("Goroutine processData stops")
			return
		case cred := <-app.credsChan:
			creds = append(creds, cred)
		case <-ticker.C:
			app.storeData(ctx, creds)
			creds = nil
		}
	}
}

func (app *App) storeData(ctx context.Context, creds []model.Creds) {
	if len(creds) > 0 {
		if err := app.db.StoreData(ctx, creds); err != nil {
			log.Warn().Msg(err.Error())
		}
	}
}

// GetData gets data.
func (app *App) GetData(ctx context.Context, userID int, site string) (*model.Creds, error) {
	creds, err := app.db.GetData(ctx, userID, site)
	if err != nil {
		return nil, err
	}

	if creds.ID == 0 {
		return nil, fmt.Errorf("app get data: %w", errors.New("creds not found"))
	}

	creds.Password, err = crypto.Decode(app.cfg.SecretKey, creds.Password)
	if err != nil {
		return nil, err
	}

	return creds, nil
}
