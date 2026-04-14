package tests

import (
	"testing"

	"github.com/ayush-sr/score-keeper/backend/internal/validator"
)

// ---------------------------------------------------------------------------
// ValidateMatchScore — Standard wins (21-X where X < 20)
// ---------------------------------------------------------------------------

func TestValidateMatchScore_StandardWins(t *testing.T) {
	for loser := 0; loser <= 19; loser++ {
		if err := validator.ValidateMatchScore(21, loser); err != nil {
			t.Errorf("ValidateMatchScore(21, %d) returned error: %v", loser, err)
		}
		if err := validator.ValidateMatchScore(loser, 21); err != nil {
			t.Errorf("ValidateMatchScore(%d, 21) returned error: %v", loser, err)
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchScore — Deuce wins (both >= 20, win by exactly 2)
// ---------------------------------------------------------------------------

func TestValidateMatchScore_DeuceWins(t *testing.T) {
	deuces := [][2]int{
		{22, 20}, {23, 21}, {24, 22}, {25, 23}, {30, 28}, {40, 38},
	}
	for _, d := range deuces {
		if err := validator.ValidateMatchScore(d[0], d[1]); err != nil {
			t.Errorf("ValidateMatchScore(%d, %d) returned error: %v", d[0], d[1], err)
		}
		if err := validator.ValidateMatchScore(d[1], d[0]); err != nil {
			t.Errorf("ValidateMatchScore(%d, %d) returned error: %v", d[1], d[0], err)
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchScore — Negative scores
// ---------------------------------------------------------------------------

func TestValidateMatchScore_NegativeScores(t *testing.T) {
	if err := validator.ValidateMatchScore(-1, 21); err == nil {
		t.Error("expected error for negative player1 score")
	}
	if err := validator.ValidateMatchScore(21, -5); err == nil {
		t.Error("expected error for negative player2 score")
	}
	if err := validator.ValidateMatchScore(-3, -2); err == nil {
		t.Error("expected error for both negative scores")
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchScore — Ties
// ---------------------------------------------------------------------------

func TestValidateMatchScore_Tie(t *testing.T) {
	ties := []int{0, 5, 10, 20, 21, 25}
	for _, score := range ties {
		if err := validator.ValidateMatchScore(score, score); err == nil {
			t.Errorf("ValidateMatchScore(%d, %d) expected error for tie", score, score)
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchScore — Winner below 21
// ---------------------------------------------------------------------------

func TestValidateMatchScore_WinnerBelow21(t *testing.T) {
	cases := [][2]int{{20, 5}, {18, 3}, {5, 0}, {20, 18}, {19, 17}}
	for _, c := range cases {
		if err := validator.ValidateMatchScore(c[0], c[1]); err == nil {
			t.Errorf("ValidateMatchScore(%d, %d) expected error: winner < 21", c[0], c[1])
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchScore — Win by less than 2
// ---------------------------------------------------------------------------

func TestValidateMatchScore_WinByLessThan2(t *testing.T) {
	cases := [][2]int{{21, 20}, {22, 21}, {25, 24}, {30, 29}}
	for _, c := range cases {
		if err := validator.ValidateMatchScore(c[0], c[1]); err == nil {
			t.Errorf("ValidateMatchScore(%d, %d) expected error: margin < 2", c[0], c[1])
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchScore — Deuce win by more than 2
// ---------------------------------------------------------------------------

func TestValidateMatchScore_DeuceWinByMoreThan2(t *testing.T) {
	cases := [][2]int{{24, 20}, {25, 21}, {30, 25}}
	for _, c := range cases {
		if err := validator.ValidateMatchScore(c[0], c[1]); err == nil {
			t.Errorf("ValidateMatchScore(%d, %d) expected error: deuce win by > 2", c[0], c[1])
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchScore — Winner above 21 when loser below 20
// ---------------------------------------------------------------------------

func TestValidateMatchScore_WinnerAbove21WhenLoserBelow20(t *testing.T) {
	cases := [][2]int{{22, 5}, {23, 18}, {25, 19}, {30, 0}}
	for _, c := range cases {
		if err := validator.ValidateMatchScore(c[0], c[1]); err == nil {
			t.Errorf("ValidateMatchScore(%d, %d) expected error: winner != 21 with loser < 20", c[0], c[1])
		}
	}
}
