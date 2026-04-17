"use client";

import { useEffect, useState, use } from "react";
import { api } from "@/lib/api";
import { Match } from "@/types/match";
import { formatDateTime } from "@/lib/utils";
import { useAuth } from "@/context/auth-context";
import { useRouter } from "next/navigation";
import Link from "next/link";

export default function MatchDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const { user } = useAuth();
  const router = useRouter();
  const [match, setMatch] = useState<Match | null>(null);
  const [loading, setLoading] = useState(true);
  const [deleting, setDeleting] = useState(false);
  const [confirmingDelete, setConfirmingDelete] = useState(false);
  const [deleteError, setDeleteError] = useState("");

  useEffect(() => {
    api.get<Match>(`/matches/${id}`).then((res) => {
      if (res.data) setMatch(res.data);
      setLoading(false);
    });
  }, [id]);

  const handleDelete = async () => {
    if (!match) return;
    setDeleting(true);
    setDeleteError("");
    const res = await api.delete(`/matches/${match.id}`);
    if (res.error) {
      setDeleteError(res.error.message);
      setDeleting(false);
    } else {
      router.push("/matches");
    }
  };

  if (loading) {
    return <p className="text-muted-foreground">Loading...</p>;
  }

  if (!match) {
    return (
      <div className="text-center">
        <p className="text-muted-foreground">Match not found.</p>
        <Link href="/matches" className="text-sm text-blue-600 dark:text-blue-400">
          Back to matches
        </Link>
      </div>
    );
  }

  const canDelete = user && match.created_by === user.id;

  return (
    <div className="mx-auto max-w-2xl">
      <Link href="/matches" className="mb-4 inline-block text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300">
        &larr; Back to matches
      </Link>

      <div className="rounded-lg border bg-card p-6">
        {/* Players and Result */}
        <div className="mb-6 flex items-center justify-center gap-6">
          <div className="text-center">
            {match.player1.avatar_url && (
              <img src={match.player1.avatar_url} alt="" className="mx-auto mb-1 h-12 w-12 rounded-full" />
            )}
            <p
              className={`font-medium ${match.winner_id === match.player1_id ? "text-green-600 dark:text-green-400" : "text-foreground"}`}
            >
              {match.player1.name}
            </p>
          </div>
          <div className="text-center">
            <p className="text-3xl font-bold font-mono text-foreground">
              {match.player1_score} - {match.player2_score}
            </p>
          </div>
          <div className="text-center">
            {match.player2.avatar_url && (
              <img src={match.player2.avatar_url} alt="" className="mx-auto mb-1 h-12 w-12 rounded-full" />
            )}
            <p
              className={`font-medium ${match.winner_id === match.player2_id ? "text-green-600 dark:text-green-400" : "text-foreground"}`}
            >
              {match.player2.name}
            </p>
          </div>
        </div>

        {/* Meta info */}
        <div className="flex items-center justify-between border-t pt-4 text-xs text-muted-foreground">
          <p>Played {formatDateTime(match.played_at)}</p>
          <p>Recorded {formatDateTime(match.created_at)}</p>
        </div>

        {/* Delete */}
        {canDelete && (
          <div className="mt-4 border-t pt-4">
            {deleteError && (
              <p className="mb-2 text-xs text-red-600 dark:text-red-400">{deleteError}</p>
            )}
            {!confirmingDelete ? (
              <button
                onClick={() => setConfirmingDelete(true)}
                className="text-sm text-red-600 hover:text-red-700 dark:text-red-400 dark:hover:text-red-300"
              >
                Delete this match
              </button>
            ) : (
              <div className="flex items-center gap-3">
                <span className="text-sm text-muted-foreground">Are you sure?</span>
                <button
                  onClick={handleDelete}
                  disabled={deleting}
                  className="rounded-md bg-red-600 px-3 py-1 text-xs font-medium text-white hover:bg-red-700 disabled:opacity-50"
                >
                  {deleting ? "Deleting..." : "Yes, delete"}
                </button>
                <button
                  onClick={() => { setConfirmingDelete(false); setDeleteError(""); }}
                  disabled={deleting}
                  className="rounded-md border px-3 py-1 text-xs font-medium text-foreground hover:bg-muted disabled:opacity-50"
                >
                  Cancel
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
