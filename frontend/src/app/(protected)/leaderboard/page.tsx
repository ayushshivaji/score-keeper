"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { User } from "@/types/user";
import { computeStandingsRow } from "@/lib/utils";
import { useAuth } from "@/context/auth-context";

export default function LeaderboardPage() {
  const { user } = useAuth();
  const [players, setPlayers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get<User[]>("/leaderboard?per_page=50").then((res) => {
      if (res.data) setPlayers(res.data);
      setLoading(false);
    });
  }, []);

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold text-foreground">Standings</h1>

      {loading ? (
        <p className="text-muted-foreground">Loading…</p>
      ) : players.length === 0 ? (
        <p className="text-muted-foreground">No matches recorded yet.</p>
      ) : (
        <div className="overflow-hidden rounded-lg border bg-card">
          <table className="min-w-full divide-y divide-border text-sm">
            <thead className="bg-muted">
              <tr>
                <th className="px-4 py-3 text-left font-medium text-muted-foreground">#</th>
                <th className="px-4 py-3 text-left font-medium text-muted-foreground">Player</th>
                <th className="px-4 py-3 text-right font-medium text-muted-foreground">Played</th>
                <th className="px-4 py-3 text-right font-medium text-muted-foreground">Wins</th>
                <th className="px-4 py-3 text-right font-medium text-muted-foreground">Losses</th>
                <th className="px-4 py-3 text-right font-medium text-muted-foreground">Win %</th>
                <th className="px-4 py-3 text-right font-medium text-muted-foreground">Points</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {players.map((p, idx) => {
                const { losses, winRate } = computeStandingsRow(p.matches_played, p.matches_won);
                const isMe = user?.id === p.id;
                return (
                  <tr key={p.id} className={isMe ? "bg-blue-50 dark:bg-blue-950/30" : ""}>
                    <td className="px-4 py-3 font-mono text-muted-foreground">{idx + 1}</td>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-2">
                        {p.avatar_url && (
                          <img
                            src={p.avatar_url}
                            alt={p.name}
                            className="h-7 w-7 rounded-full"
                          />
                        )}
                        <span className={`font-medium ${isMe ? "text-blue-700 dark:text-blue-300" : "text-foreground"}`}>
                          {p.name}
                        </span>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-right text-foreground">{p.matches_played}</td>
                    <td className="px-4 py-3 text-right font-medium text-green-600 dark:text-green-400">{p.matches_won}</td>
                    <td className="px-4 py-3 text-right text-red-600 dark:text-red-400">{losses}</td>
                    <td className="px-4 py-3 text-right text-foreground">{winRate}%</td>
                    <td className="px-4 py-3 text-right font-mono text-foreground">{p.total_points}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
