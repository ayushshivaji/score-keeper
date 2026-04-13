package tests

import (
	"testing"

	"github.com/ayush-sr/score-keeper/backend/internal/dto"
	"github.com/ayush-sr/score-keeper/backend/internal/validator"
)

// ---------------------------------------------------------------------------
// ValidateMatchFormat
// ---------------------------------------------------------------------------

func TestValidateMatchFormat_ValidFormats(t *testing.T) {
	for _, format := range []int{3, 5, 7} {
		if err := validator.ValidateMatchFormat(format); err != nil {
			t.Errorf("validator.ValidateMatchFormat(%d) returned error: %v", format, err)
		}
	}
}

func TestValidateMatchFormat_InvalidFormats(t *testing.T) {
	for _, format := range []int{0, 1, 2, 4, 6, 8, 9, -1, 100} {
		if err := validator.ValidateMatchFormat(format); err == nil {
			t.Errorf("validator.ValidateMatchFormat(%d) expected error, got nil", format)
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateSetScore — Standard wins (21-X where X < 20)
// ---------------------------------------------------------------------------

func TestValidateSetScore_StandardWins(t *testing.T) {
	for loser := 0; loser <= 19; loser++ {
		if err := validator.ValidateSetScore(21, loser); err != nil {
			t.Errorf("validator.ValidateSetScore(21, %d) returned error: %v", loser, err)
		}
		if err := validator.ValidateSetScore(loser, 21); err != nil {
			t.Errorf("validator.ValidateSetScore(%d, 21) returned error: %v", loser, err)
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateSetScore — Deuce wins (both >= 20, win by exactly 2)
// ---------------------------------------------------------------------------

func TestValidateSetScore_DeuceWins(t *testing.T) {
	deuces := [][2]int{
		{22, 20}, {23, 21}, {24, 22}, {25, 23}, {30, 28}, {40, 38},
	}
	for _, d := range deuces {
		if err := validator.ValidateSetScore(d[0], d[1]); err != nil {
			t.Errorf("validator.ValidateSetScore(%d, %d) returned error: %v", d[0], d[1], err)
		}
		if err := validator.ValidateSetScore(d[1], d[0]); err != nil {
			t.Errorf("validator.ValidateSetScore(%d, %d) returned error: %v", d[1], d[0], err)
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateSetScore — Negative scores
// ---------------------------------------------------------------------------

func TestValidateSetScore_NegativeScores(t *testing.T) {
	if err := validator.ValidateSetScore(-1, 21); err == nil {
		t.Error("expected error for negative player1 score")
	}
	if err := validator.ValidateSetScore(21, -5); err == nil {
		t.Error("expected error for negative player2 score")
	}
	if err := validator.ValidateSetScore(-3, -2); err == nil {
		t.Error("expected error for both negative scores")
	}
}

// ---------------------------------------------------------------------------
// ValidateSetScore — Ties
// ---------------------------------------------------------------------------

func TestValidateSetScore_Tie(t *testing.T) {
	ties := []int{0, 5, 10, 20, 21, 25}
	for _, score := range ties {
		if err := validator.ValidateSetScore(score, score); err == nil {
			t.Errorf("validator.ValidateSetScore(%d, %d) expected error for tie", score, score)
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateSetScore — Winner below 21
// ---------------------------------------------------------------------------

func TestValidateSetScore_WinnerBelow21(t *testing.T) {
	cases := [][2]int{{20, 5}, {18, 3}, {5, 0}, {20, 18}, {19, 17}}
	for _, c := range cases {
		if err := validator.ValidateSetScore(c[0], c[1]); err == nil {
			t.Errorf("validator.ValidateSetScore(%d, %d) expected error: winner < 21", c[0], c[1])
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateSetScore — Win by less than 2
// ---------------------------------------------------------------------------

func TestValidateSetScore_WinByLessThan2(t *testing.T) {
	cases := [][2]int{{21, 20}, {22, 21}, {25, 24}, {30, 29}}
	for _, c := range cases {
		if err := validator.ValidateSetScore(c[0], c[1]); err == nil {
			t.Errorf("validator.ValidateSetScore(%d, %d) expected error: margin < 2", c[0], c[1])
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateSetScore — Deuce win by more than 2
// ---------------------------------------------------------------------------

func TestValidateSetScore_DeuceWinByMoreThan2(t *testing.T) {
	cases := [][2]int{{24, 20}, {25, 21}, {30, 25}}
	for _, c := range cases {
		if err := validator.ValidateSetScore(c[0], c[1]); err == nil {
			t.Errorf("validator.ValidateSetScore(%d, %d) expected error: deuce win by > 2", c[0], c[1])
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateSetScore — Winner above 21 when loser below 20
// ---------------------------------------------------------------------------

func TestValidateSetScore_WinnerAbove21WhenLoserBelow20(t *testing.T) {
	cases := [][2]int{{22, 5}, {23, 18}, {25, 19}, {30, 0}}
	for _, c := range cases {
		if err := validator.ValidateSetScore(c[0], c[1]); err == nil {
			t.Errorf("validator.ValidateSetScore(%d, %d) expected error: winner != 21 with loser < 20", c[0], c[1])
		}
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchSets — Best of 3
// ---------------------------------------------------------------------------

func TestValidateMatchSets_BestOf3_Player1Wins2_0(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 3,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 21, Player2Score: 19},
		},
	}
	if err := validator.ValidateMatchSets(req); err != nil {
		t.Errorf("expected valid 2-0, got error: %v", err)
	}
}

func TestValidateMatchSets_BestOf3_Player2Wins2_1(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 3,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 17},
			{Player1Score: 19, Player2Score: 21},
			{Player1Score: 18, Player2Score: 21},
		},
	}
	if err := validator.ValidateMatchSets(req); err != nil {
		t.Errorf("expected valid 1-2, got error: %v", err)
	}
}

func TestValidateMatchSets_BestOf3_TooFewSets(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 3,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 5},
		},
	}
	if err := validator.ValidateMatchSets(req); err == nil {
		t.Error("expected error: only 1 set in best-of-3")
	}
}

func TestValidateMatchSets_BestOf3_TooManySets(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 3,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 21, Player2Score: 5},
		},
	}
	if err := validator.ValidateMatchSets(req); err == nil {
		t.Error("expected error: 4 sets in best-of-3")
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchSets — Best of 5
// ---------------------------------------------------------------------------

func TestValidateMatchSets_BestOf5_Player1Wins3_0(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 5,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 3},
			{Player1Score: 21, Player2Score: 17},
			{Player1Score: 21, Player2Score: 19},
		},
	}
	if err := validator.ValidateMatchSets(req); err != nil {
		t.Errorf("expected valid 3-0, got error: %v", err)
	}
}

func TestValidateMatchSets_BestOf5_Player1Wins3_2(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 5,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 17},
			{Player1Score: 19, Player2Score: 21},
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 18, Player2Score: 21},
			{Player1Score: 21, Player2Score: 19},
		},
	}
	if err := validator.ValidateMatchSets(req); err != nil {
		t.Errorf("expected valid 3-2, got error: %v", err)
	}
}

func TestValidateMatchSets_BestOf5_Player2Wins3_1(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 5,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 17},
			{Player1Score: 5, Player2Score: 21},
			{Player1Score: 18, Player2Score: 21},
			{Player1Score: 19, Player2Score: 21},
		},
	}
	if err := validator.ValidateMatchSets(req); err != nil {
		t.Errorf("expected valid 1-3, got error: %v", err)
	}
}

func TestValidateMatchSets_BestOf5_DeuceInSets(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 5,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 22, Player2Score: 20},
			{Player1Score: 24, Player2Score: 22},
			{Player1Score: 23, Player2Score: 21},
		},
	}
	if err := validator.ValidateMatchSets(req); err != nil {
		t.Errorf("expected valid with deuce sets, got error: %v", err)
	}
}

func TestValidateMatchSets_BestOf5_TooFewSets(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 5,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 21, Player2Score: 5},
		},
	}
	if err := validator.ValidateMatchSets(req); err == nil {
		t.Error("expected error: only 2 sets in best-of-5")
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchSets — Best of 7
// ---------------------------------------------------------------------------

func TestValidateMatchSets_BestOf7_Player1Wins4_0(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 7,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 3},
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 21, Player2Score: 17},
			{Player1Score: 21, Player2Score: 19},
		},
	}
	if err := validator.ValidateMatchSets(req); err != nil {
		t.Errorf("expected valid 4-0, got error: %v", err)
	}
}

func TestValidateMatchSets_BestOf7_Player2Wins4_3(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 7,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 5, Player2Score: 21},
			{Player1Score: 21, Player2Score: 17},
			{Player1Score: 19, Player2Score: 21},
			{Player1Score: 21, Player2Score: 18},
			{Player1Score: 6, Player2Score: 21},
			{Player1Score: 19, Player2Score: 21},
		},
	}
	if err := validator.ValidateMatchSets(req); err != nil {
		t.Errorf("expected valid 3-4, got error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchSets — Extra sets after match decided
// ---------------------------------------------------------------------------

func TestValidateMatchSets_ExtraSetsAfterWin_BestOf3(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 3,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 21, Player2Score: 17},
			{Player1Score: 21, Player2Score: 3},
		},
	}
	if err := validator.ValidateMatchSets(req); err == nil {
		t.Error("expected error: extra set after match decided")
	}
}

func TestValidateMatchSets_ExtraSetsAfterWin_BestOf5(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 5,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 21, Player2Score: 17},
			{Player1Score: 21, Player2Score: 3},
			{Player1Score: 21, Player2Score: 19},
		},
	}
	if err := validator.ValidateMatchSets(req); err == nil {
		t.Error("expected error: extra set after match decided")
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchSets — No winner
// ---------------------------------------------------------------------------

func TestValidateMatchSets_NoWinner_BestOf5(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 5,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 5, Player2Score: 21},
			{Player1Score: 21, Player2Score: 17},
			{Player1Score: 19, Player2Score: 21},
		},
	}
	if err := validator.ValidateMatchSets(req); err == nil {
		t.Error("expected error: 2-2 is not a winner in best-of-5")
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchSets — Invalid set score within match
// ---------------------------------------------------------------------------

func TestValidateMatchSets_InvalidSetScoreInMiddle(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 5,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 18, Player2Score: 5},
			{Player1Score: 21, Player2Score: 17},
		},
	}
	if err := validator.ValidateMatchSets(req); err == nil {
		t.Error("expected error: set 2 has invalid score (18-5)")
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchSets — Invalid match format
// ---------------------------------------------------------------------------

func TestValidateMatchSets_InvalidFormat(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 4,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 21, Player2Score: 5},
			{Player1Score: 21, Player2Score: 17},
		},
	}
	if err := validator.ValidateMatchSets(req); err == nil {
		t.Error("expected error: invalid format 4")
	}
}

// ---------------------------------------------------------------------------
// ValidateMatchSets — Edge cases
// ---------------------------------------------------------------------------

func TestValidateMatchSets_EmptySets(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 3,
		Sets:        []dto.SetScoreRequest{},
	}
	if err := validator.ValidateMatchSets(req); err == nil {
		t.Error("expected error: no sets provided")
	}
}

func TestValidateMatchSets_ZeroZeroSet(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 3,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 0, Player2Score: 0},
			{Player1Score: 21, Player2Score: 5},
		},
	}
	if err := validator.ValidateMatchSets(req); err == nil {
		t.Error("expected error: 0-0 set is a tie")
	}
}

func TestValidateMatchSets_AllDeuceMatch(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 3,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 22, Player2Score: 20},
			{Player1Score: 20, Player2Score: 22},
			{Player1Score: 25, Player2Score: 23},
		},
	}
	if err := validator.ValidateMatchSets(req); err != nil {
		t.Errorf("expected valid all-deuce match, got error: %v", err)
	}
}

func TestValidateMatchSets_HighDeuceScore(t *testing.T) {
	req := &dto.CreateMatchRequest{
		MatchFormat: 3,
		Sets: []dto.SetScoreRequest{
			{Player1Score: 40, Player2Score: 38},
			{Player1Score: 21, Player2Score: 19},
		},
	}
	if err := validator.ValidateMatchSets(req); err != nil {
		t.Errorf("expected valid high-deuce match, got error: %v", err)
	}
}
