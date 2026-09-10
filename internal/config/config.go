package config

import (
	"os"
)

// Config holds configuration loaded from environment variables.
type Config struct {
	Port        string
	Env         string
	DatabaseURL string
}

// Load reads configuration from environment variables.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://agentbox:acer@localhost:5432/agentbox_db"
	}
	return Config{
		Port:        port,
		Env:         env,
		DatabaseURL: dbURL,
	}
}
