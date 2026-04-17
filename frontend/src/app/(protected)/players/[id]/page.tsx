"use client";

import { useEffect, useState, use } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import { UserProfile } from "@/types/user";
import { Match } from "@/types/match";
import { MatchCard } from "@/components/match-card";
import { formatStreak } from "@/lib/utils";
import { useAuth } from "@/context/auth-context";

export default function PlayerProfilePage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const { user } = useAuth();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [matches, setMatches] = useState<Match[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    Promise.all([
      api.get<UserProfile>(`/users/${id}`),
      api.get<Match[]>(`/matches?player_id=${id}&per_page=10`),
    ]).then(([profileRes, matchesRes]) => {
      if (profileRes.data) setProfile(profileRes.data);
      if (matchesRes.data) setMatches(matchesRes.data);
      setLoading(false);
    });
  }, [id]);

  if (loading) return <p className="text-muted-foreground">Loading…</p>;
  if (!profile) {
    return (
      <div className="text-center">
        <p className="text-muted-foreground">Player not found.</p>
        <Link href="/players" className="text-sm text-blue-600 dark:text-blue-400">
          Back to players
        </Link>
      </div>
    );
  }

  const winPct = Math.round(profile.win_rate * 100);

  return (
    <div className="mx-auto max-w-3xl">
      <Link href="/players" className="mb-4 inline-block text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300">
        &larr; Back to players
      </Link>

      {/* Header */}
      <div className="mb-6 flex items-center gap-4 rounded-lg border bg-card p-6">
        {profile.avatar_url ? (
          <img src={profile.avatar_url} alt={profile.name} className="h-16 w-16 rounded-full" />
        ) : (
          <div className="h-16 w-16 rounded-full bg-muted" />
        )}
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-foreground">{profile.name}</h1>
          <p className="text-sm text-muted-foreground">{profile.email}</p>
        </div>
        {user && user.id !== profile.id && (
          <Link
            href={`/players/${profile.id}/vs/${user.id}`}
            className="rounded-md border bg-card px-3 py-2 text-xs font-medium text-foreground hover:bg-muted"
          >
            Head-to-head vs you
          </Link>
        )}
      </div>

      {/* Stats cards */}
      <div className="mb-6 grid grid-cols-2 gap-4 sm:grid-cols-4">
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Matches</p>
          <p className="text-2xl font-bold text-foreground">{profile.matches_played}</p>
        </div>
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Record</p>
          <p className="text-2xl font-bold">
            <span className="text-green-600 dark:text-green-400">{profile.matches_won}</span>
            <span className="text-muted-foreground">-</span>
            <span className="text-red-600 dark:text-red-400">{profile.losses}</span>
          </p>
        </div>
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Win Rate</p>
          <p className="text-2xl font-bold text-foreground">{winPct}%</p>
        </div>
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Total Points</p>
          <p className="text-2xl font-bold text-foreground">{profile.total_points}</p>
        </div>
      </div>

      {/* Streaks + form */}
      <div className="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-3">
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Current Streak</p>
          <p className="text-xl font-semibold text-foreground">{formatStreak(profile.current_streak)}</p>
        </div>
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Longest Win Streak</p>
          <p className="text-xl font-semibold text-green-600 dark:text-green-400">{profile.longest_win_streak}</p>
        </div>
        <div className="rounded-lg border bg-card p-4">
          <p className="text-sm text-muted-foreground">Recent Form</p>
          <div className="mt-1 flex gap-1">
            {profile.recent_form.length === 0 ? (
              <span className="text-sm text-muted-foreground">—</span>
            ) : (
              profile.recent_form.map((r, i) => (
                <span
                  key={i}
                  className={`inline-flex h-6 w-6 items-center justify-center rounded text-xs font-bold ${
                    r === "W"
                      ? "bg-green-100 text-green-700 dark:bg-green-900/50 dark:text-green-300"
                      : "bg-red-100 text-red-700 dark:bg-red-900/50 dark:text-red-300"
                  }`}
                >
                  {r}
                </span>
              ))
            )}
          </div>
        </div>
      </div>

      {/* Recent matches */}
      <h2 className="mb-3 text-lg font-semibold text-foreground">Recent Matches</h2>
      {matches.length === 0 ? (
        <p className="text-muted-foreground">No matches yet.</p>
      ) : (
        <div className="space-y-3">
          {matches.map((match) => (
            <MatchCard key={match.id} match={match} />
          ))}
        </div>
      )}
    </div>
  );
}
