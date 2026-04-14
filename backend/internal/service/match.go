package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ayush-sr/score-keeper/backend/internal/dto"
	"github.com/ayush-sr/score-keeper/backend/internal/model"
	"github.com/ayush-sr/score-keeper/backend/internal/repository"
	"github.com/ayush-sr/score-keeper/backend/internal/validator"
	"github.com/google/uuid"
)

type MatchService struct {
	matchRepo *repository.MatchRepository
	userRepo  *repository.UserRepository
}

func NewMatchService(matchRepo *repository.MatchRepository, userRepo *repository.UserRepository) *MatchService {
	return &MatchService{matchRepo: matchRepo, userRepo: userRepo}
}

func (s *MatchService) CreateMatch(ctx context.Context, req *dto.CreateMatchRequest, createdBy uuid.UUID) (*model.MatchWithDetails, error) {
	if req.Player1ID == req.Player2ID {
		return nil, fmt.Errorf("player 1 and player 2 must be different")
	}

	if err := validator.ValidateMatchScore(req.Player1Score, req.Player2Score); err != nil {
		return nil, err
	}

	winnerID := req.Player1ID
	if req.Player2Score > req.Player1Score {
		winnerID = req.Player2ID
	}

	tx, err := s.matchRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	match := &model.Match{
		Player1ID:    req.Player1ID,
		Player2ID:    req.Player2ID,
		WinnerID:     winnerID,
		Player1Score: req.Player1Score,
		Player2Score: req.Player2Score,
		PlayedAt:     req.PlayedAt,
		CreatedBy:    createdBy,
	}

	if err := s.matchRepo.CreateMatch(ctx, tx, match); err != nil {
		return nil, err
	}

	p1Won, p2Won := 0, 0
	if winnerID == req.Player1ID {
		p1Won = 1
	} else {
		p2Won = 1
	}
	if err := s.userRepo.AdjustMatchStats(ctx, tx, req.Player1ID, 1, p1Won, req.Player1Score); err != nil {
		return nil, err
	}
	if err := s.userRepo.AdjustMatchStats(ctx, tx, req.Player2ID, 1, p2Won, req.Player2Score); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return s.matchRepo.GetByID(ctx, match.ID)
}

func (s *MatchService) GetMatch(ctx context.Context, id uuid.UUID) (*model.MatchWithDetails, error) {
	return s.matchRepo.GetByID(ctx, id)
}

func (s *MatchService) ListMatches(ctx context.Context, playerID *uuid.UUID, page, perPage int) ([]model.MatchWithDetails, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 50 {
		perPage = 20
	}
	return s.matchRepo.List(ctx, playerID, page, perPage)
}

func (s *MatchService) DeleteMatch(ctx context.Context, matchID, userID uuid.UUID) error {
	createdBy, err := s.matchRepo.GetCreatedBy(ctx, matchID)
	if err != nil {
		return fmt.Errorf("match not found")
	}
	if createdBy != userID {
		return fmt.Errorf("only the creator can delete a match")
	}

	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil || match == nil {
		return fmt.Errorf("match not found")
	}

	if time.Since(match.CreatedAt) > 24*time.Hour {
		return fmt.Errorf("matches can only be deleted within 24 hours of creation")
	}

	tx, err := s.matchRepo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	p1Won, p2Won := 0, 0
	if match.WinnerID == match.Player1ID {
		p1Won = 1
	} else {
		p2Won = 1
	}
	if err := s.userRepo.AdjustMatchStats(ctx, tx, match.Player1ID, -1, -p1Won, -match.Player1Score); err != nil {
		return err
	}
	if err := s.userRepo.AdjustMatchStats(ctx, tx, match.Player2ID, -1, -p2Won, -match.Player2Score); err != nil {
		return err
	}

	if err := s.matchRepo.Delete(ctx, tx, matchID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
