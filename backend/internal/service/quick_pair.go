package service

import (
	"context"
	"errors"

	"github.com/ayush-sr/score-keeper/backend/internal/model"
	"github.com/ayush-sr/score-keeper/backend/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrSamePlayer        = errors.New("player1 and player2 must be different")
	ErrPlayerNotFound    = errors.New("player not found")
	ErrQuickPairExists   = errors.New("this pair is already saved")
	ErrQuickPairNotFound = errors.New("quick pair not found")
)

type QuickPairService struct {
	repo     *repository.QuickPairRepository
	userRepo *repository.UserRepository
}

func NewQuickPairService(repo *repository.QuickPairRepository, userRepo *repository.UserRepository) *QuickPairService {
	return &QuickPairService{repo: repo, userRepo: userRepo}
}

func (s *QuickPairService) List(ctx context.Context, userID uuid.UUID) ([]model.QuickPairWithPlayers, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *QuickPairService) Create(ctx context.Context, userID, player1ID, player2ID uuid.UUID) (*model.QuickPairWithPlayers, error) {
	if player1ID == player2ID {
		return nil, ErrSamePlayer
	}

	p1, err := s.userRepo.GetByID(ctx, player1ID)
	if err != nil {
		return nil, err
	}
	if p1 == nil {
		return nil, ErrPlayerNotFound
	}
	p2, err := s.userRepo.GetByID(ctx, player2ID)
	if err != nil {
		return nil, err
	}
	if p2 == nil {
		return nil, ErrPlayerNotFound
	}

	qp, err := s.repo.Create(ctx, userID, player1ID, player2ID)
	if err != nil {
		if errors.Is(err, repository.ErrQuickPairDuplicate) {
			return nil, ErrQuickPairExists
		}
		return nil, err
	}
	return qp, nil
}

func (s *QuickPairService) Delete(ctx context.Context, userID, id uuid.UUID) error {
	deleted, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrQuickPairNotFound
	}
	return nil
}
