package model

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID               uuid.UUID  `json:"id"`
	Player1ID        uuid.UUID  `json:"player1_id"`
	Player2ID        uuid.UUID  `json:"player2_id"`
	WinnerID         uuid.UUID  `json:"winner_id"`
	MatchFormat      int        `json:"match_format"`
	Player1SetsWon   int        `json:"player1_sets_won"`
	Player2SetsWon   int        `json:"player2_sets_won"`
	TournamentMatchID *uuid.UUID `json:"tournament_match_id,omitempty"`
	PlayedAt         time.Time  `json:"played_at"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        uuid.UUID  `json:"created_by"`
}

type MatchSet struct {
	ID           uuid.UUID `json:"id"`
	MatchID      uuid.UUID `json:"match_id"`
	SetNumber    int       `json:"set_number"`
	Player1Score int       `json:"player1_score"`
	Player2Score int       `json:"player2_score"`
}

type MatchWithDetails struct {
	Match
	Player1 User       `json:"player1"`
	Player2 User       `json:"player2"`
	Sets    []MatchSet `json:"sets"`
}
