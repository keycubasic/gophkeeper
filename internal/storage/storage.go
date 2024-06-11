// Package storage contains interface of database.
package storage

import (
	"context"
	"errors"

	"github.com/keycubasic/gophkeeper/internal/model"
)

// ErrUniqueViolation an unique violetion error.
var ErrUniqueViolation = errors.New("storage unique violation")

// Interface of database.
//
//go:generate mockgen -source app.go -destination mock_storage.go -package mock Database
type Database interface {
	// PingContext checks the database connection.
	PingContext(ctx context.Context) error
	// RegisterUser registers new user and returns user Id.
	RegisterUser(ctx context.Context, username, hashedPassword string) (int, error)
	// FindUser finds user.
	FindUser(ctx context.Context, username string) (*model.User, error)
	// Close closes the connection to the database.
	Close() error
}
