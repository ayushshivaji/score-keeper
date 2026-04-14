package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMatchRequest struct {
	Player1ID    uuid.UUID `json:"player1_id" binding:"required"`
	Player2ID    uuid.UUID `json:"player2_id" binding:"required"`
	Player1Score int       `json:"player1_score" binding:"min=0"`
	Player2Score int       `json:"player2_score" binding:"min=0"`
	PlayedAt     time.Time `json:"played_at" binding:"required"`
}
