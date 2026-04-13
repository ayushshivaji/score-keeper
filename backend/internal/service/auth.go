package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ayush-sr/score-keeper/backend/internal/config"
	"github.com/ayush-sr/score-keeper/backend/internal/model"
	"github.com/ayush-sr/score-keeper/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) UpsertUser(ctx context.Context, googleID, email, name string, avatarURL *string) (*model.User, error) {
	return s.userRepo.UpsertByGoogleID(ctx, googleID, email, name, avatarURL)
}

func (s *AuthService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)
	hash := hashToken(token)

	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.userRepo.StoreRefreshToken(ctx, userID, hash, expiresAt); err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, string, error) {
	hash := hashToken(refreshToken)
	userID, expiresAt, err := s.userRepo.GetRefreshToken(ctx, hash)
	if err != nil {
		return "", "", err
	}
	if userID == uuid.Nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}
	if time.Now().After(expiresAt) {
		_ = s.userRepo.DeleteRefreshToken(ctx, hash)
		return "", "", fmt.Errorf("refresh token expired")
	}

	// Rotate: delete old, create new
	_ = s.userRepo.DeleteRefreshToken(ctx, hash)

	accessToken, err := s.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}
	newRefresh, err := s.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return "", "", err
	}
	return accessToken, newRefresh, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.DeleteUserRefreshTokens(ctx, userID)
}

func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
