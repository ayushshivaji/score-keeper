import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(dateString: string): string {
  return new Date(dateString).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

export function formatDateTime(dateString: string): string {
  return new Date(dateString).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function computeStandingsRow(matchesPlayed: number, matchesWon: number): {
  losses: number;
  winRate: number;
} {
  const losses = Math.max(0, matchesPlayed - matchesWon);
  const winRate = matchesPlayed > 0 ? Math.round((matchesWon / matchesPlayed) * 100) : 0;
  return { losses, winRate };
}

// Formats a signed streak integer (positive = wins, negative = losses, 0 = none)
// into a short human label like "W3", "L2", or "—".
export function formatStreak(streak: number): string {
  if (streak > 0) return `W${streak}`;
  if (streak < 0) return `L${-streak}`;
  return "—";
}

// Validates the final score of a standalone table tennis match.
// Rules: not tied, winner >= 21, win by >= 2, with the standard deuce rule
// (once both reach 20, the winner must win by exactly 2).
// Mirrors backend/internal/validator/match.go ValidateMatchScore.
export function validateMatchScore(player1: number, player2: number): string | null {
  if (player1 < 0 || player2 < 0) return "Scores must be non-negative";
  if (player1 === player2) return "Scores cannot be tied";
  const winner = Math.max(player1, player2);
  const loser = Math.min(player1, player2);
  if (winner < 21) return "Winner needs at least 21";
  if (winner - loser < 2) return "Must win by 2";
  if (loser >= 20 && winner !== loser + 2) return "In deuce, must win by exactly 2";
  if (loser < 20 && winner !== 21) return "Score must be 21 when opponent < 20";
  return null;
}
