package tests

import (
	"context"
	"testing"

	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/google/uuid"
)

// Only the "same player" validation path runs before the repo is touched,
// so that's the only branch we can unit-test without a DB. This mirrors
// the convention in user_service_test.go and auth_service_test.go.

func TestQuickPairCreate_SamePlayerRejected(t *testing.T) {
	svc := service.NewQuickPairService(nil, nil)
	playerID := uuid.New()
	userID := uuid.New()

	_, err := svc.Create(context.Background(), userID, playerID, playerID)
	if err != service.ErrSamePlayer {
		t.Fatalf("expected ErrSamePlayer, got %v", err)
	}
}
