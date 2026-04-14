package validator

import "fmt"

// ValidateMatchScore checks the final score of a single standalone table tennis match.
// Rules: both scores non-negative, not tied, winner >= 21, win by >= 2, with the
// standard deuce rule (when both players reach 20, the winner must win by exactly 2).
func ValidateMatchScore(player1Score, player2Score int) error {
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
		return fmt.Errorf("a match cannot end in a tie")
	}
	if winner < 21 {
		return fmt.Errorf("winning score must be at least 21")
	}
	if winner-loser < 2 {
		return fmt.Errorf("winner must win by at least 2 points")
	}
	if loser >= 20 && winner != loser+2 {
		return fmt.Errorf("in deuce, winner must win by exactly 2 points")
	}
	if loser < 20 && winner != 21 {
		return fmt.Errorf("winning score must be 21 when opponent has less than 20")
	}
	return nil
}
