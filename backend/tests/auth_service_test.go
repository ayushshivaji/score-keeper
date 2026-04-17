package tests

import (
	"context"
	"testing"
	"time"

	"github.com/ayush-sr/score-keeper/backend/internal/config"
	"github.com/ayush-sr/score-keeper/backend/internal/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func newTestConfig() *config.Config {
	return &config.Config{
		JWTSecret: "test-secret-key-at-least-32-chars-long!!",
	}
}

// ---------------------------------------------------------------------------
// JWT Access Token — Generation
// ---------------------------------------------------------------------------

func TestGenerateAccessToken_ReturnsNonEmptyString(t *testing.T) {
	svc := service.NewAuthService(nil, newTestConfig())
	token, err := svc.GenerateAccessToken(uuid.New())
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestGenerateAccessToken_DifferentUsersGetDifferentTokens(t *testing.T) {
	svc := service.NewAuthService(nil, newTestConfig())
	token1, _ := svc.GenerateAccessToken(uuid.New())
	token2, _ := svc.GenerateAccessToken(uuid.New())
	if token1 == token2 {
		t.Fatal("expected different tokens for different users")
	}
}

// ---------------------------------------------------------------------------
// JWT Access Token — Validation
// ---------------------------------------------------------------------------

func TestValidateAccessToken_ValidToken(t *testing.T) {
	svc := service.NewAuthService(nil, newTestConfig())
	userID := uuid.New()
	token, _ := svc.GenerateAccessToken(userID)

	claims, err := svc.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected %s, got %s", userID, claims.UserID)
	}
}

func TestValidateAccessToken_InvalidToken(t *testing.T) {
	svc := service.NewAuthService(nil, newTestConfig())
	_, err := svc.ValidateAccessToken("not-a-real-token")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateAccessToken_EmptyString(t *testing.T) {
	svc := service.NewAuthService(nil, newTestConfig())
	_, err := svc.ValidateAccessToken("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	svc1 := service.NewAuthService(nil, &config.Config{JWTSecret: "secret-one-at-least-32-characters!!"})
	svc2 := service.NewAuthService(nil, &config.Config{JWTSecret: "secret-two-at-least-32-characters!!"})

	token, _ := svc1.GenerateAccessToken(uuid.New())
	_, err := svc2.ValidateAccessToken(token)
	if err == nil {
		t.Fatal("expected error with wrong secret")
	}
}

func TestValidateAccessToken_ExpiredToken(t *testing.T) {
	cfg := newTestConfig()
	claims := service.JWTClaims{
		UserID: uuid.New(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(cfg.JWTSecret))

	svc := service.NewAuthService(nil, cfg)
	_, err := svc.ValidateAccessToken(tokenString)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestValidateAccessToken_WrongSigningMethod(t *testing.T) {
	svc := service.NewAuthService(nil, newTestConfig())

	token := jwt.NewWithClaims(jwt.SigningMethodNone, &service.JWTClaims{
		UserID: uuid.New(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	})
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	_, err := svc.ValidateAccessToken(tokenString)
	if err == nil {
		t.Fatal("expected error for 'none' signing method")
	}
}

// ---------------------------------------------------------------------------
// JWT Access Token — Claims content
// ---------------------------------------------------------------------------

func TestGenerateAccessToken_HasCorrectExpiry(t *testing.T) {
	svc := service.NewAuthService(nil, newTestConfig())
	tokenString, _ := svc.GenerateAccessToken(uuid.New())
	claims, _ := svc.ValidateAccessToken(tokenString)

	expectedExpiry := time.Now().Add(15 * time.Minute)
	diff := claims.ExpiresAt.Time.Sub(expectedExpiry)
	if diff > 5*time.Second || diff < -5*time.Second {
		t.Errorf("expiry too far from expected 15min: diff=%v", diff)
	}
}

func TestGenerateAccessToken_HasIssuedAt(t *testing.T) {
	svc := service.NewAuthService(nil, newTestConfig())
	tokenString, _ := svc.GenerateAccessToken(uuid.New())
	claims, _ := svc.ValidateAccessToken(tokenString)

	if claims.IssuedAt == nil {
		t.Fatal("expected IssuedAt to be set")
	}
	if time.Since(claims.IssuedAt.Time) > 5*time.Second {
		t.Error("IssuedAt too far in the past")
	}
}

// ---------------------------------------------------------------------------
// Static Login
// ---------------------------------------------------------------------------

func TestLoginStatic_DisabledWhenUsernameEmpty(t *testing.T) {
	cfg := newTestConfig()
	cfg.StaticLoginUsername = ""
	cfg.StaticLoginPassword = "pw"
	svc := service.NewAuthService(nil, cfg)

	_, err := svc.LoginStatic(context.Background(), "admin", "pw")
	if err != service.ErrStaticLoginDisabled {
		t.Fatalf("expected ErrStaticLoginDisabled, got %v", err)
	}
}

func TestLoginStatic_DisabledWhenPasswordEmpty(t *testing.T) {
	cfg := newTestConfig()
	cfg.StaticLoginUsername = "admin"
	cfg.StaticLoginPassword = ""
	svc := service.NewAuthService(nil, cfg)

	_, err := svc.LoginStatic(context.Background(), "admin", "")
	if err != service.ErrStaticLoginDisabled {
		t.Fatalf("expected ErrStaticLoginDisabled, got %v", err)
	}
}

func TestLoginStatic_WrongUsername(t *testing.T) {
	cfg := newTestConfig()
	cfg.StaticLoginUsername = "admin"
	cfg.StaticLoginPassword = "pw"
	svc := service.NewAuthService(nil, cfg)

	_, err := svc.LoginStatic(context.Background(), "root", "pw")
	if err != service.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginStatic_WrongPassword(t *testing.T) {
	cfg := newTestConfig()
	cfg.StaticLoginUsername = "admin"
	cfg.StaticLoginPassword = "pw"
	svc := service.NewAuthService(nil, cfg)

	_, err := svc.LoginStatic(context.Background(), "admin", "nope")
	if err != service.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLoginStatic_LengthMismatchStillRejected(t *testing.T) {
	// ConstantTimeCompare returns 0 for differing lengths — make sure that
	// path still surfaces as an invalid-credentials error rather than panicking.
	cfg := newTestConfig()
	cfg.StaticLoginUsername = "admin"
	cfg.StaticLoginPassword = "longpassword"
	svc := service.NewAuthService(nil, cfg)

	_, err := svc.LoginStatic(context.Background(), "admin", "short")
	if err != service.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

