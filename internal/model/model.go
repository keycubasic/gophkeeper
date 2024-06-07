// Package model contains models.
package model

// User stores user info.
type User struct {
	ID             int
	Username       string
	HashedPassword string
}
