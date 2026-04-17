package model

import (
	"time"

	"github.com/google/uuid"
)

type QuickPair struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Player1ID uuid.UUID `json:"player1_id"`
	Player2ID uuid.UUID `json:"player2_id"`
	CreatedAt time.Time `json:"created_at"`
}

type QuickPairWithPlayers struct {
	QuickPair
	Player1 User `json:"player1"`
	Player2 User `json:"player2"`
}
