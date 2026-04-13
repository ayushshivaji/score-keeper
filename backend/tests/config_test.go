package tests

import (
	"os"
	"testing"

	"github.com/ayush-sr/score-keeper/backend/internal/config"
)

func setEnvVars(t *testing.T, vars map[string]string) {
	t.Helper()
	for k, v := range vars {
		t.Setenv(k, v)
	}
}

func validEnvVars() map[string]string {
	return map[string]string{
		"DATABASE_URL":         "postgres://localhost:5432/test",
		"GOOGLE_CLIENT_ID":     "test-client-id",
		"GOOGLE_CLIENT_SECRET": "test-client-secret",
		"JWT_SECRET":           "test-jwt-secret",
	}
}

func TestLoad_AllRequired(t *testing.T) {
	setEnvVars(t, validEnvVars())
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load() failed: %v", err)
	}
	if cfg.DatabaseURL != "postgres://localhost:5432/test" {
		t.Errorf("unexpected DatabaseURL: %s", cfg.DatabaseURL)
	}
	if cfg.GoogleClientID != "test-client-id" {
		t.Errorf("unexpected GoogleClientID: %s", cfg.GoogleClientID)
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	vars := validEnvVars()
	delete(vars, "DATABASE_URL")
	setEnvVars(t, vars)
	os.Unsetenv("DATABASE_URL")
	if _, err := config.Load(); err == nil {
		t.Fatal("expected error for missing DATABASE_URL")
	}
}

func TestLoad_MissingGoogleClientID(t *testing.T) {
	vars := validEnvVars()
	delete(vars, "GOOGLE_CLIENT_ID")
	setEnvVars(t, vars)
	os.Unsetenv("GOOGLE_CLIENT_ID")
	if _, err := config.Load(); err == nil {
		t.Fatal("expected error for missing GOOGLE_CLIENT_ID")
	}
}

func TestLoad_MissingGoogleClientSecret(t *testing.T) {
	vars := validEnvVars()
	delete(vars, "GOOGLE_CLIENT_SECRET")
	setEnvVars(t, vars)
	os.Unsetenv("GOOGLE_CLIENT_SECRET")
	if _, err := config.Load(); err == nil {
		t.Fatal("expected error for missing GOOGLE_CLIENT_SECRET")
	}
}

func TestLoad_MissingJWTSecret(t *testing.T) {
	vars := validEnvVars()
	delete(vars, "JWT_SECRET")
	setEnvVars(t, vars)
	os.Unsetenv("JWT_SECRET")
	if _, err := config.Load(); err == nil {
		t.Fatal("expected error for missing JWT_SECRET")
	}
}

func TestLoad_Defaults(t *testing.T) {
	setEnvVars(t, validEnvVars())
	cfg, _ := config.Load()
	if cfg.Port != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.Port)
	}
	if cfg.FrontendURL != "http://localhost:3000" {
		t.Errorf("expected default frontend URL, got %s", cfg.FrontendURL)
	}
	if cfg.GoogleRedirectURL != "http://localhost:8080/api/v1/auth/google/callback" {
		t.Errorf("expected default redirect URL, got %s", cfg.GoogleRedirectURL)
	}
}

func TestLoad_CustomPort(t *testing.T) {
	vars := validEnvVars()
	vars["PORT"] = "9090"
	setEnvVars(t, vars)
	cfg, _ := config.Load()
	if cfg.Port != "9090" {
		t.Errorf("expected 9090, got %s", cfg.Port)
	}
}

func TestLoad_CustomFrontendURL(t *testing.T) {
	vars := validEnvVars()
	vars["FRONTEND_URL"] = "https://myapp.com"
	setEnvVars(t, vars)
	cfg, _ := config.Load()
	if cfg.FrontendURL != "https://myapp.com" {
		t.Errorf("expected https://myapp.com, got %s", cfg.FrontendURL)
	}
}
