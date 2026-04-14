package service

// ComputeStreakStats walks a chronological list of match results (oldest first)
// and returns the player's current streak, longest win streak, and longest loss
// streak.
//
// The current streak is signed: positive values mean consecutive wins, negative
// values mean consecutive losses. Zero means the player has no matches.
func ComputeStreakStats(resultsOldestFirst []bool) (current, longestWin, longestLoss int) {
	for _, won := range resultsOldestFirst {
		if won {
			if current >= 0 {
				current++
			} else {
				current = 1
			}
			if current > longestWin {
				longestWin = current
			}
		} else {
			if current <= 0 {
				current--
			} else {
				current = -1
			}
			if -current > longestLoss {
				longestLoss = -current
			}
		}
	}
	return current, longestWin, longestLoss
}

// RecentForm returns the last n results as single-character labels, newest first.
// Expects resultsOldestFirst and returns ["W", "L", ...] newest first.
func RecentForm(resultsOldestFirst []bool, n int) []string {
	start := len(resultsOldestFirst) - n
	if start < 0 {
		start = 0
	}
	slice := resultsOldestFirst[start:]
	out := make([]string, 0, len(slice))
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] {
			out = append(out, "W")
		} else {
			out = append(out, "L")
		}
	}
	return out
}
