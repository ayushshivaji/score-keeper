"use client";

import Link from "next/link";
import { Match } from "@/types/match";
import { formatDate } from "@/lib/utils";

interface MatchCardProps {
  match: Match;
}

export function MatchCard({ match }: MatchCardProps) {
  return (
    <Link
      href={`/matches/${match.id}`}
      className="block rounded-lg border bg-card p-4 hover:border-ring transition-colors"
    >
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2">
            {match.player1.avatar_url && (
              <img src={match.player1.avatar_url} alt="" className="h-6 w-6 rounded-full" />
            )}
            <span
              className={`font-medium ${match.winner_id === match.player1_id ? "text-green-600 dark:text-green-400" : "text-foreground"}`}
            >
              {match.player1.name}
            </span>
          </div>
          <span className="text-muted-foreground">vs</span>
          <div className="flex items-center gap-2">
            {match.player2.avatar_url && (
              <img src={match.player2.avatar_url} alt="" className="h-6 w-6 rounded-full" />
            )}
            <span
              className={`font-medium ${match.winner_id === match.player2_id ? "text-green-600 dark:text-green-400" : "text-foreground"}`}
            >
              {match.player2.name}
            </span>
          </div>
        </div>
        <div className="text-right">
          <p className="font-mono text-sm font-bold text-foreground">
            {match.player1_score}-{match.player2_score}
          </p>
        </div>
      </div>
      <div className="mt-2 flex items-center justify-end">
        <p className="text-xs text-muted-foreground">{formatDate(match.played_at)}</p>
      </div>
    </Link>
  );
}
