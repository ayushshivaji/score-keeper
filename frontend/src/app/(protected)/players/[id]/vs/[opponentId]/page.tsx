"use client";

import { useEffect, useState, use } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import { HeadToHead } from "@/types/match";
import { MatchCard } from "@/components/match-card";

export default function HeadToHeadPage({
  params,
}: {
  params: Promise<{ id: string; opponentId: string }>;
}) {
  const { id, opponentId } = use(params);
  const [h2h, setH2h] = useState<HeadToHead | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    api
      .get<HeadToHead>(`/users/${id}/head-to-head/${opponentId}`)
      .then((res) => {
        if (res.data) setH2h(res.data);
        setLoading(false);
      });
  }, [id, opponentId]);

  if (loading) return <p className="text-muted-foreground">Loading…</p>;
  if (!h2h) {
    return (
      <div className="text-center">
        <p className="text-muted-foreground">Head-to-head not found.</p>
        <Link href="/players" className="text-sm text-blue-600 dark:text-blue-400">
          Back to players
        </Link>
      </div>
    );
  }

  const p1Leading = h2h.player1_wins > h2h.player2_wins;
  const p2Leading = h2h.player2_wins > h2h.player1_wins;

  return (
    <div className="mx-auto max-w-3xl">
      <Link
        href={`/players/${h2h.player1.id}`}
        className="mb-4 inline-block text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
      >
        &larr; Back to {h2h.player1.name}
      </Link>

      {/* Side-by-side header */}
      <div className="mb-6 grid grid-cols-3 items-center gap-4 rounded-lg border bg-card p-6">
        <div className="text-center">
          {h2h.player1.avatar_url ? (
            <img
              src={h2h.player1.avatar_url}
              alt={h2h.player1.name}
              className="mx-auto mb-2 h-16 w-16 rounded-full"
            />
          ) : (
            <div className="mx-auto mb-2 h-16 w-16 rounded-full bg-muted" />
          )}
          <Link
            href={`/players/${h2h.player1.id}`}
            className={`font-semibold ${p1Leading ? "text-green-600 dark:text-green-400" : "text-foreground"}`}
          >
            {h2h.player1.name}
          </Link>
        </div>
        <div className="text-center">
          <p className="text-4xl font-bold font-mono text-foreground">
            {h2h.player1_wins} - {h2h.player2_wins}
          </p>
          <p className="mt-1 text-xs text-muted-foreground">{h2h.total_matches} matches</p>
        </div>
        <div className="text-center">
          {h2h.player2.avatar_url ? (
            <img
              src={h2h.player2.avatar_url}
              alt={h2h.player2.name}
              className="mx-auto mb-2 h-16 w-16 rounded-full"
            />
          ) : (
            <div className="mx-auto mb-2 h-16 w-16 rounded-full bg-muted" />
          )}
          <Link
            href={`/players/${h2h.player2.id}`}
            className={`font-semibold ${p2Leading ? "text-green-600 dark:text-green-400" : "text-foreground"}`}
          >
            {h2h.player2.name}
          </Link>
        </div>
      </div>

      {/* Point differential */}
      <div className="mb-6 grid grid-cols-2 gap-4">
        <div className="rounded-lg border bg-card p-4 text-center">
          <p className="text-sm text-muted-foreground">{h2h.player1.name} Points</p>
          <p className="text-2xl font-bold text-foreground">{h2h.player1_points}</p>
        </div>
        <div className="rounded-lg border bg-card p-4 text-center">
          <p className="text-sm text-muted-foreground">{h2h.player2.name} Points</p>
          <p className="text-2xl font-bold text-foreground">{h2h.player2_points}</p>
        </div>
      </div>

      {/* Matches */}
      <h2 className="mb-3 text-lg font-semibold text-foreground">All Matches</h2>
      {h2h.matches.length === 0 ? (
        <p className="text-muted-foreground">These players haven&apos;t played each other yet.</p>
      ) : (
        <div className="space-y-3">
          {h2h.matches.map((match) => (
            <MatchCard key={match.id} match={match} />
          ))}
        </div>
      )}
    </div>
  );
}
