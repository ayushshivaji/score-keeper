package validator

import (
	"fmt"

	"github.com/ayush-sr/score-keeper/backend/internal/dto"
)

func ValidateMatchFormat(format int) error {
	if format != 3 && format != 5 && format != 7 {
		return fmt.Errorf("match_format must be 3, 5, or 7")
	}
	return nil
}

func ValidateSetScore(player1Score, player2Score int) error {
	if player1Score < 0 || player2Score < 0 {
		return fmt.Errorf("scores must be non-negative")
	}

	winner := player1Score
	loser := player2Score
	if player2Score > player1Score {
		winner = player2Score
		loser = player1Score
	}

	if winner == loser {
		return fmt.Errorf("a set cannot end in a tie")
	}

	if winner < 21 {
		return fmt.Errorf("winning score must be at least 21")
	}

	if winner-loser < 2 {
		return fmt.Errorf("winner must win by at least 2 points")
	}

	// If both players reached 20+, must win by exactly 2 (deuce scenario)
	if loser >= 20 && winner != loser+2 {
		return fmt.Errorf("in deuce, winner must win by exactly 2 points")
	}

	// If loser has less than 20, winner must have exactly 21
	if loser < 20 && winner != 21 {
		return fmt.Errorf("winning score must be 21 when opponent has less than 20")
	}

	return nil
}

func ValidateMatchSets(req *dto.CreateMatchRequest) error {
	if err := ValidateMatchFormat(req.MatchFormat); err != nil {
		return err
	}

	setsToWin := (req.MatchFormat / 2) + 1
	minSets := setsToWin
	maxSets := req.MatchFormat

	if len(req.Sets) < minSets || len(req.Sets) > maxSets {
		return fmt.Errorf("expected between %d and %d sets for best-of-%d", minSets, maxSets, req.MatchFormat)
	}

	player1Wins := 0
	player2Wins := 0

	for i, set := range req.Sets {
		if err := ValidateSetScore(set.Player1Score, set.Player2Score); err != nil {
			return fmt.Errorf("set %d: %w", i+1, err)
		}

		if set.Player1Score > set.Player2Score {
			player1Wins++
		} else {
			player2Wins++
		}

		// Check if match was already decided before this set
		if i < len(req.Sets)-1 {
			if player1Wins == setsToWin || player2Wins == setsToWin {
				return fmt.Errorf("match was decided after set %d, but %d sets were provided", i+1, len(req.Sets))
			}
		}
	}

	if player1Wins != setsToWin && player2Wins != setsToWin {
		return fmt.Errorf("no player has won %d sets", setsToWin)
	}

	return nil
}
