package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL        string
	Port               string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	JWTSecret          string
	FrontendURL        string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		Port:               getEnv("PORT", "8080"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.GoogleClientID == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID is required")
	}
	if cfg.GoogleClientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_SECRET is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
