package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ayush-sr/score-keeper/backend/internal/dto"
	"github.com/ayush-sr/score-keeper/backend/internal/model"
	"github.com/ayush-sr/score-keeper/backend/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo  *repository.UserRepository
	matchRepo *repository.MatchRepository
}

func NewUserService(userRepo *repository.UserRepository, matchRepo *repository.MatchRepository) *UserService {
	return &UserService{userRepo: userRepo, matchRepo: matchRepo}
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// CreatePlayer adds a player who does not sign in via Google. Used when the
// app is operated via the static login and players are managed by the admin.
func (s *UserService) CreatePlayer(ctx context.Context, name string) (*model.User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if len(name) > 255 {
		return nil, fmt.Errorf("name is too long")
	}
	playerID := uuid.NewString()
	syntheticGoogleID := "player:" + playerID
	syntheticEmail := "player-" + playerID + "@local"
	return s.userRepo.CreatePlayer(ctx, syntheticGoogleID, syntheticEmail, name)
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

func (s *UserService) GetLeaderboard(ctx context.Context, page, perPage int) ([]model.User, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 50 {
		perPage = 20
	}
	return s.userRepo.ListLeaderboard(ctx, page, perPage)
}

// GetUserProfile returns the player's user row plus computed streaks and recent form.
func (s *UserService) GetUserProfile(ctx context.Context, id uuid.UUID) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, nil
	}

	results, err := s.matchRepo.GetPlayerResultsChronological(ctx, id)
	if err != nil {
		return nil, err
	}

	current, longestWin, longestLoss := ComputeStreakStats(results)

	losses := user.MatchesPlayed - user.MatchesWon
	if losses < 0 {
		losses = 0
	}
	var winRate float64
	if user.MatchesPlayed > 0 {
		winRate = float64(user.MatchesWon) / float64(user.MatchesPlayed)
	}

	return &dto.UserProfileResponse{
		User:              *user,
		Losses:            losses,
		WinRate:           winRate,
		CurrentStreak:     current,
		LongestWinStreak:  longestWin,
		LongestLossStreak: longestLoss,
		RecentForm:        RecentForm(results, 5),
	}, nil
}

// GetHeadToHead returns aggregated and per-match head-to-head stats between two players.
func (s *UserService) GetHeadToHead(ctx context.Context, p1, p2 uuid.UUID) (*dto.HeadToHeadResponse, error) {
	if p1 == p2 {
		return nil, fmt.Errorf("players must be different")
	}

	player1, err := s.userRepo.GetByID(ctx, p1)
	if err != nil || player1 == nil {
		return nil, fmt.Errorf("player not found")
	}
	player2, err := s.userRepo.GetByID(ctx, p2)
	if err != nil || player2 == nil {
		return nil, fmt.Errorf("player not found")
	}

	total, p1Wins, p2Wins, p1Points, p2Points, err := s.matchRepo.HeadToHeadAggregates(ctx, p1, p2)
	if err != nil {
		return nil, err
	}

	matches, err := s.matchRepo.HeadToHeadMatches(ctx, p1, p2)
	if err != nil {
		return nil, err
	}

	return &dto.HeadToHeadResponse{
		Player1:       *player1,
		Player2:       *player2,
		TotalMatches:  total,
		Player1Wins:   p1Wins,
		Player2Wins:   p2Wins,
		Player1Points: p1Points,
		Player2Points: p2Points,
		Matches:       matches,
	}, nil
}
