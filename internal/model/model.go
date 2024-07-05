// Package model contains models.
package model

// User stores user info.
type User struct {
	ID             int
	Username       string
	HashedPassword string
}

// Creds stores user credentials.
type Creds struct {
	ID       int    `db:"cred_id"`
	UserID   int    `db:"user_id"`
	Site     string `db:"user_site"`
	Login    string `db:"user_login"`
	Password string `db:"user_password"`
	Metadata string `db:"metadata"`
}
