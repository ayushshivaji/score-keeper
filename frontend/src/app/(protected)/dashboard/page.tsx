"use client";

import { useAuth } from "@/context/auth-context";
import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { Match } from "@/types/match";
import { formatDate } from "@/lib/utils";
import Link from "next/link";

export default function DashboardPage() {
  const { user } = useAuth();
  const [recentMatches, setRecentMatches] = useState<Match[]>([]);

  useEffect(() => {
    if (user) {
      api
        .get<Match[]>(`/matches?player_id=${user.id}&per_page=5`)
        .then((res) => {
          if (res.data) setRecentMatches(res.data);
        });
    }
  }, [user]);

  if (!user) return null;

  const winRate =
    user.matches_played > 0
      ? Math.round((user.matches_won / user.matches_played) * 100)
      : 0;

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-foreground">
          Welcome back, {user.name}
        </h1>
      </div>

      {/* Stats Cards */}
      <div className="mb-8 grid grid-cols-2 gap-4 sm:grid-cols-4">
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Matches Played</p>
          <p className="text-2xl font-bold text-foreground">
            {user.matches_played}
          </p>
        </div>
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Wins</p>
          <p className="text-2xl font-bold text-green-600 dark:text-green-400">
            {user.matches_won}
          </p>
        </div>
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Losses</p>
          <p className="text-2xl font-bold text-red-600 dark:text-red-400">
            {user.matches_played - user.matches_won}
          </p>
        </div>
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Win Rate</p>
          <p className="text-2xl font-bold text-foreground">{winRate}%</p>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="mb-8 flex gap-4">
        <Link
          href="/matches/new"
          className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          Record a Match
        </Link>
        <Link
          href="/matches"
          className="rounded-lg border bg-card px-4 py-2 text-sm font-medium text-foreground hover:bg-muted"
        >
          View All Matches
        </Link>
      </div>

      {/* Recent Matches */}
      <div>
        <h2 className="mb-4 text-lg font-semibold text-foreground">
          Recent Matches
        </h2>
        {recentMatches.length === 0 ? (
          <p className="text-muted-foreground">No matches yet. Record your first match!</p>
        ) : (
          <div className="space-y-3">
            {recentMatches.map((match) => (
              <Link
                key={match.id}
                href={`/matches/${match.id}`}
                className="block rounded-lg border bg-card p-4 hover:border-ring transition-colors"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <span
                      className={`font-medium ${match.winner_id === match.player1_id ? "text-green-600 dark:text-green-400" : "text-foreground"}`}
                    >
                      {match.player1.name}
                    </span>
                    <span className="text-muted-foreground">vs</span>
                    <span
                      className={`font-medium ${match.winner_id === match.player2_id ? "text-green-600 dark:text-green-400" : "text-foreground"}`}
                    >
                      {match.player2.name}
                    </span>
                  </div>
                  <div className="text-right">
                    <p className="font-mono text-sm text-foreground">
                      {match.player1_score}-{match.player2_score}
                    </p>
                  </div>
                </div>
                <p className="mt-1 text-xs text-muted-foreground">
                  {formatDate(match.played_at)}
                </p>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
