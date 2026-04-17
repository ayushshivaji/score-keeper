package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/ayush-sr/score-keeper/backend/internal/service"
)

// These tests cover the input-validation paths of UserService.CreatePlayer,
// which run before the repository is ever touched. The happy path (actual
// insert) is a DB-backed concern and lives outside this unit suite.

func TestCreatePlayer_EmptyName(t *testing.T) {
	svc := service.NewUserService(nil, nil)
	_, err := svc.CreatePlayer(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCreatePlayer_WhitespaceOnlyName(t *testing.T) {
	svc := service.NewUserService(nil, nil)
	_, err := svc.CreatePlayer(context.Background(), "   \t\n  ")
	if err == nil {
		t.Fatal("expected error for whitespace-only name")
	}
}

func TestCreatePlayer_NameTooLong(t *testing.T) {
	svc := service.NewUserService(nil, nil)
	longName := strings.Repeat("a", 256)
	_, err := svc.CreatePlayer(context.Background(), longName)
	if err == nil {
		t.Fatal("expected error for name > 255 chars")
	}
}
