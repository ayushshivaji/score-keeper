"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { User } from "@/types/user";
import { QuickPair } from "@/types/quick-pair";
import { CreateMatchRequest } from "@/types/match";
import { validateMatchScore } from "@/lib/utils";

function clampScore(raw: string): number {
  const n = parseInt(raw, 10);
  if (Number.isNaN(n)) return 0;
  if (n < 0) return 0;
  if (n > 40) return 40;
  return n;
}

export default function NewMatchPage() {
  const router = useRouter();
  const [players, setPlayers] = useState<User[]>([]);
  const [quickPairs, setQuickPairs] = useState<QuickPair[]>([]);
  const [player1Id, setPlayer1Id] = useState("");
  const [player2Id, setPlayer2Id] = useState("");
  const [player1Score, setPlayer1Score] = useState(0);
  const [player2Score, setPlayer2Score] = useState(0);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  // Inline "add quick pair" form state
  const [showAddPair, setShowAddPair] = useState(false);
  const [newPairP1, setNewPairP1] = useState("");
  const [newPairP2, setNewPairP2] = useState("");
  const [pairError, setPairError] = useState<string | null>(null);
  const [savingPair, setSavingPair] = useState(false);

  useEffect(() => {
    api.get<User[]>(`/users?per_page=50`).then((res) => {
      if (res.data) setPlayers(res.data);
    });
  }, []);

  const loadQuickPairs = useCallback(() => {
    api.get<QuickPair[]>(`/quick-pairs`).then((res) => {
      if (res.data) setQuickPairs(res.data);
    });
  }, []);

  useEffect(() => {
    loadQuickPairs();
  }, [loadQuickPairs]);

  const handlePlayer1Change = (id: string) => {
    setPlayer1Id(id);
    if (id && player1Score === 0) {
      setPlayer1Score(21);
    }
  };

  const handleQuickPairClick = (pair: QuickPair) => {
    setPlayer1Id(pair.player1_id);
    setPlayer2Id(pair.player2_id);
  };

  const handleDeletePair = async (pairId: string) => {
    await api.delete(`/quick-pairs/${pairId}`);
    loadQuickPairs();
  };

  const handleSavePair = async () => {
    if (!newPairP1 || !newPairP2) return;
    setPairError(null);
    setSavingPair(true);
    const res = await api.post<QuickPair>("/quick-pairs", {
      player1_id: newPairP1,
      player2_id: newPairP2,
    });
    setSavingPair(false);
    if (res.error) {
      if (res.error.code === "CONFLICT") {
        setPairError("This pair is already saved");
      } else if (res.error.code === "BAD_REQUEST") {
        setPairError("Players must be different");
      } else {
        setPairError(res.error.message || "Failed to save pair");
      }
      return;
    }
    setNewPairP1("");
    setNewPairP2("");
    setShowAddPair(false);
    loadQuickPairs();
  };

  const scoreError =
    player1Score === 0 && player2Score === 0
      ? null
      : validateMatchScore(player1Score, player2Score);
  const hasValidScore = !scoreError && (player1Score > 0 || player2Score > 0);
  const scoreWinner = hasValidScore
    ? player1Score > player2Score
      ? "player1"
      : "player2"
    : null;

  const player1Name = players.find((p) => p.id === player1Id)?.name || "Player 1";
  const player2Name = players.find((p) => p.id === player2Id)?.name || "Player 2";

  const handleSubmit = async () => {
    if (!player1Id || !player2Id || player1Id === player2Id || !hasValidScore) return;

    setSubmitting(true);
    setError("");

    const req: CreateMatchRequest = {
      player1_id: player1Id,
      player2_id: player2Id,
      player1_score: player1Score,
      player2_score: player2Score,
      played_at: new Date().toISOString(),
    };

    const res = await api.post("/matches", req);
    if (res.error) {
      setError(res.error.message);
      setSubmitting(false);
    } else {
      router.push("/matches");
    }
  };

  const scoreInputClass = (isWinner: boolean) =>
    `w-full rounded-md border bg-card px-3 py-2 text-center font-mono text-lg text-foreground ${
      scoreError
        ? "border-red-400 bg-red-50 dark:bg-red-950/40"
        : isWinner
          ? "border-green-400 bg-green-50 dark:bg-green-950/40"
          : ""
    }`;

  return (
    <div className="mx-auto max-w-2xl">
      <h1 className="mb-6 text-2xl font-bold text-foreground">Record a Match</h1>

      {/* Quick pairs */}
      <div className="mb-6">
        <div className="mb-2 flex items-center justify-between">
          <span className="text-sm font-medium text-muted-foreground">Quick pairs</span>
          {!showAddPair && (
            <button
              type="button"
              onClick={() => {
                setShowAddPair(true);
                setPairError(null);
              }}
              className="text-xs text-blue-600 hover:underline"
            >
              + Add pair
            </button>
          )}
        </div>

        {quickPairs.length === 0 && !showAddPair && (
          <p className="text-xs text-muted-foreground">
            No quick pairs yet. Add one for faster match entry.
          </p>
        )}

        {quickPairs.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {quickPairs.map((pair) => (
              <div
                key={pair.id}
                className="flex items-center gap-1 rounded-full border bg-card pl-3 pr-1 py-1 text-xs"
              >
                <button
                  type="button"
                  onClick={() => handleQuickPairClick(pair)}
                  className="font-medium text-foreground hover:underline"
                >
                  {pair.player1.name} vs {pair.player2.name}
                </button>
                <button
                  type="button"
                  onClick={() => handleDeletePair(pair.id)}
                  aria-label={`Remove ${pair.player1.name} vs ${pair.player2.name}`}
                  className="ml-1 flex h-5 w-5 items-center justify-center rounded-full text-muted-foreground hover:bg-muted hover:text-foreground"
                >
                  ×
                </button>
              </div>
            ))}
          </div>
        )}

        {showAddPair && (
          <div className="mt-3 rounded-lg border bg-card p-3">
            <div className="grid grid-cols-2 gap-2">
              <select
                value={newPairP1}
                onChange={(e) => setNewPairP1(e.target.value)}
                className="rounded-md border bg-background px-2 py-1 text-sm"
              >
                <option value="">Player 1</option>
                {players
                  .filter((p) => p.id !== newPairP2)
                  .map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
              </select>
              <select
                value={newPairP2}
                onChange={(e) => setNewPairP2(e.target.value)}
                className="rounded-md border bg-background px-2 py-1 text-sm"
              >
                <option value="">Player 2</option>
                {players
                  .filter((p) => p.id !== newPairP1)
                  .map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
              </select>
            </div>
            {pairError && (
              <p className="mt-2 text-xs text-red-600">{pairError}</p>
            )}
            <div className="mt-3 flex gap-2">
              <button
                type="button"
                onClick={handleSavePair}
                disabled={!newPairP1 || !newPairP2 || savingPair}
                className="rounded-md bg-blue-600 px-3 py-1 text-xs font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {savingPair ? "Saving…" : "Save pair"}
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowAddPair(false);
                  setNewPairP1("");
                  setNewPairP2("");
                  setPairError(null);
                }}
                className="rounded-md border px-3 py-1 text-xs font-medium text-foreground hover:bg-muted"
              >
                Cancel
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Player Selection */}
      <div className="mb-6 grid grid-cols-2 gap-4">
        <div>
          <label className="mb-1 block text-sm font-medium text-foreground">Player 1</label>
          <select
            value={player1Id}
            onChange={(e) => handlePlayer1Change(e.target.value)}
            className="w-full rounded-md border bg-card px-3 py-2 text-sm text-foreground"
          >
            <option value="">Select player</option>
            {players
              .filter((p) => p.id !== player2Id)
              .map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
          </select>
        </div>
        <div>
          <label className="mb-1 block text-sm font-medium text-foreground">Player 2</label>
          <select
            value={player2Id}
            onChange={(e) => setPlayer2Id(e.target.value)}
            className="w-full rounded-md border bg-card px-3 py-2 text-sm text-foreground"
          >
            <option value="">Select player</option>
            {players
              .filter((p) => p.id !== player1Id)
              .map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
          </select>
        </div>
      </div>

      {/* Score entry */}
      <div className="mb-6 rounded-lg border bg-card p-4">
        <div className="mb-2 flex items-center justify-between">
          <span className="text-sm font-medium text-muted-foreground">Final Score</span>
          {scoreError && (
            <span className="text-xs text-red-600 dark:text-red-400">{scoreError}</span>
          )}
          {hasValidScore && (
            <span className="text-xs text-green-600 dark:text-green-400">
              {player1Score}-{player2Score}
            </span>
          )}
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium text-foreground">{player1Name}</label>
            <input
              type="number"
              min={0}
              max={40}
              value={player1Score}
              onChange={(e) => setPlayer1Score(clampScore(e.target.value))}
              onFocus={(e) => e.target.select()}
              className={scoreInputClass(scoreWinner === "player1")}
            />
          </div>
          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium text-foreground">{player2Name}</label>
            <input
              type="number"
              min={0}
              max={40}
              value={player2Score}
              onChange={(e) => setPlayer2Score(clampScore(e.target.value))}
              onFocus={(e) => e.target.select()}
              className={scoreInputClass(scoreWinner === "player2")}
            />
          </div>
        </div>
      </div>

      {/* Match Summary */}
      {hasValidScore && scoreWinner && (
        <div className="mb-6 rounded-lg border border-green-300 bg-green-50 p-4 dark:border-green-900 dark:bg-green-950/40">
          <p className="text-sm font-medium text-green-800 dark:text-green-300">
            {scoreWinner === "player1" ? player1Name : player2Name} wins{" "}
            {Math.max(player1Score, player2Score)}-{Math.min(player1Score, player2Score)}
          </p>
        </div>
      )}

      {error && (
        <div className="mb-6 rounded-lg border border-red-300 bg-red-50 p-4 dark:border-red-900 dark:bg-red-950/40">
          <p className="text-sm text-red-800 dark:text-red-300">{error}</p>
        </div>
      )}

      {/* Submit */}
      <button
        onClick={handleSubmit}
        disabled={
          !hasValidScore ||
          !player1Id ||
          !player2Id ||
          player1Id === player2Id ||
          submitting
        }
        className="w-full rounded-lg bg-blue-600 px-4 py-3 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {submitting ? "Saving..." : "Save Match"}
      </button>
    </div>
  );
}
