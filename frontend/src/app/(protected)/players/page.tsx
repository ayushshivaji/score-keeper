"use client";

import { useCallback, useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import { User } from "@/types/user";
import { computeStandingsRow } from "@/lib/utils";

export default function PlayersPage() {
  const [players, setPlayers] = useState<User[]>([]);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [newName, setNewName] = useState("");
  const [creating, setCreating] = useState(false);
  const [createError, setCreateError] = useState<string | null>(null);

  const load = useCallback(() => {
    setLoading(true);
    api
      .get<User[]>(`/users?per_page=50&search=${encodeURIComponent(search)}`)
      .then((res) => {
        if (res.data) setPlayers(res.data);
        setLoading(false);
      });
  }, [search]);

  useEffect(() => {
    load();
  }, [load]);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    const name = newName.trim();
    if (!name) return;
    setCreating(true);
    setCreateError(null);
    const res = await api.post<User>("/users", { name });
    setCreating(false);
    if (res.error) {
      setCreateError(res.error.message || "Failed to add player");
      return;
    }
    setNewName("");
    load();
  };

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold text-foreground">Players</h1>

      <form
        onSubmit={handleCreate}
        className="mb-6 flex w-full max-w-md flex-col gap-2 sm:flex-row"
      >
        <input
          type="text"
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          placeholder="Add player by name..."
          className="flex-1 rounded-md border bg-card px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground"
        />
        <button
          type="submit"
          disabled={creating || !newName.trim()}
          className="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
        >
          {creating ? "Adding…" : "Add player"}
        </button>
      </form>
      {createError && (
        <p className="mb-4 text-sm text-red-600">{createError}</p>
      )}

      <input
        type="text"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Search by name..."
        className="mb-6 w-full max-w-md rounded-md border bg-card px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground"
      />

      {loading ? (
        <p className="text-muted-foreground">Loading…</p>
      ) : players.length === 0 ? (
        <p className="text-muted-foreground">No players found.</p>
      ) : (
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {players.map((p) => {
            const { losses, winRate } = computeStandingsRow(p.matches_played, p.matches_won);
            return (
              <Link
                key={p.id}
                href={`/players/${p.id}`}
                className="flex items-center gap-3 rounded-lg border bg-card p-4 hover:border-ring transition-colors"
              >
                {p.avatar_url ? (
                  <img src={p.avatar_url} alt={p.name} className="h-10 w-10 rounded-full" />
                ) : (
                  <div className="h-10 w-10 rounded-full bg-muted" />
                )}
                <div className="min-w-0 flex-1">
                  <p className="truncate font-medium text-foreground">{p.name}</p>
                  <p className="text-xs text-muted-foreground">
                    {p.matches_won}W-{losses}L · {winRate}% · {p.total_points} pts
                  </p>
                </div>
              </Link>
            );
          })}
        </div>
      )}
    </div>
  );
}
