package model

import (
	"time"

	"github.com/google/uuid"
)

type Match struct {
	ID           uuid.UUID `json:"id"`
	Player1ID    uuid.UUID `json:"player1_id"`
	Player2ID    uuid.UUID `json:"player2_id"`
	WinnerID     uuid.UUID `json:"winner_id"`
	Player1Score int       `json:"player1_score"`
	Player2Score int       `json:"player2_score"`
	PlayedAt     time.Time `json:"played_at"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    uuid.UUID `json:"created_by"`
}

type MatchWithDetails struct {
	Match
	Player1 User `json:"player1"`
	Player2 User `json:"player2"`
}
