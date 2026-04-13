package service

import (
	"context"

	"github.com/ayush-sr/score-keeper/backend/internal/model"
	"github.com/ayush-sr/score-keeper/backend/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, search string, page, perPage int) ([]model.User, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 50 {
		perPage = 20
	}
	return s.userRepo.List(ctx, search, page, perPage)
}
