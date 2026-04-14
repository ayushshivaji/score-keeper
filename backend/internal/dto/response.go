package dto

import "github.com/ayush-sr/score-keeper/backend/internal/model"

type APIResponse struct {
	Data  interface{} `json:"data"`
	Error *APIError   `json:"error"`
	Meta  *Meta       `json:"meta,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

func Success(data interface{}) APIResponse {
	return APIResponse{Data: data, Error: nil}
}

func SuccessWithMeta(data interface{}, meta *Meta) APIResponse {
	return APIResponse{Data: data, Error: nil, Meta: meta}
}

func ErrorResponse(code, message string) APIResponse {
	return APIResponse{Data: nil, Error: &APIError{Code: code, Message: message}}
}

// UserProfileResponse is the richer shape returned by GET /users/:id.
// It embeds the base User row and adds computed streak/form stats.
type UserProfileResponse struct {
	model.User
	Losses            int      `json:"losses"`
	WinRate           float64  `json:"win_rate"`
	CurrentStreak     int      `json:"current_streak"`
	LongestWinStreak  int      `json:"longest_win_streak"`
	LongestLossStreak int      `json:"longest_loss_streak"`
	RecentForm        []string `json:"recent_form"`
}

// HeadToHeadResponse summarises two players' shared history.
type HeadToHeadResponse struct {
	Player1       model.User                `json:"player1"`
	Player2       model.User                `json:"player2"`
	TotalMatches  int                       `json:"total_matches"`
	Player1Wins   int                       `json:"player1_wins"`
	Player2Wins   int                       `json:"player2_wins"`
	Player1Points int                       `json:"player1_points"`
	Player2Points int                       `json:"player2_points"`
	Matches       []model.MatchWithDetails  `json:"matches"`
}
