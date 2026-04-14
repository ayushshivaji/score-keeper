package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	GoogleID      string    `json:"google_id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	AvatarURL     *string   `json:"avatar_url"`
	MatchesPlayed int       `json:"matches_played"`
	MatchesWon    int       `json:"matches_won"`
	TotalPoints   int       `json:"total_points"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
