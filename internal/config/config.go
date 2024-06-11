// Package config contains functions for working with application configuration.
package config

import (
	_ "embed"
	"flag"
	"os"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	//go:embed config.json
	embedConfig []byte
)

// Config stores the application configuration.
type Config struct {
	ServerAddress string `json:"server_address"`
	LogLevel      string `json:"log_level"`
	DatabaseDSN   string `json:"database_dsn"`
	JWTSecretKey  string `json:"jwt_secret_key"`
	TokenDuration int    `json:"token_diration"`
	EnableHTTPS   bool   `json:"enable_https"`
}

// ParseConfig parses the application configuration.
func (cfg *Config) ParseConfig() {
	if embedConfig != nil && json.Valid(embedConfig) {
		json.Unmarshal(embedConfig, cfg)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server start address")
	flag.StringVar(&cfg.LogLevel, "l", cfg.LogLevel, "Log level")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "Database connection string")
	flag.StringVar(&cfg.JWTSecretKey, "k", cfg.JWTSecretKey, "JWT secret key")
	flag.IntVar(&cfg.TokenDuration, "t", cfg.TokenDuration, "JWT token duration (hours)")
	flag.BoolFunc("s", "Sign of https protocol support", func(s string) error {
		cfg.EnableHTTPS = true
		return nil
	})

	flag.Parse()

	if a := os.Getenv("GK_SERVER_ADDRESS"); a != "" {
		cfg.ServerAddress = a
	}

	if l := os.Getenv("GK_LOG_LEVEL"); l != "" {
		cfg.LogLevel = l
	}

	if d := os.Getenv("GK_DATABASE_DSN"); d != "" {
		cfg.DatabaseDSN = d
	}

	if k := os.Getenv("GK_JWT_SECRET_KEY"); k != "" {
		cfg.JWTSecretKey = k
	}

	if _, ok := os.LookupEnv("GK_ENABLE_HTTPS"); ok {
		cfg.EnableHTTPS = ok
	}
}
