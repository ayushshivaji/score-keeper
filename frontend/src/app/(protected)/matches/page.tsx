"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import { Match } from "@/types/match";
import { Meta } from "@/types/api";
import { MatchCard } from "@/components/match-card";
import Link from "next/link";

export default function MatchesPage() {
  const [matches, setMatches] = useState<Match[]>([]);
  const [meta, setMeta] = useState<Meta | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    api.get<Match[]>(`/matches?page=${page}&per_page=20`).then((res) => {
      if (res.data) setMatches(res.data);
      if (res.meta) setMeta(res.meta);
      setLoading(false);
    });
  }, [page]);

  const totalPages = meta ? Math.ceil(meta.total / meta.per_page) : 1;

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <h1 className="text-2xl font-bold text-foreground">Match History</h1>
        <Link
          href="/matches/new"
          className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >
          Record Match
        </Link>
      </div>

      {loading ? (
        <p className="text-muted-foreground">Loading...</p>
      ) : matches.length === 0 ? (
        <div className="rounded-lg border bg-card p-8 text-center">
          <p className="text-muted-foreground">No matches recorded yet.</p>
          <Link
            href="/matches/new"
            className="mt-2 inline-block text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
          >
            Record your first match
          </Link>
        </div>
      ) : (
        <>
          <div className="space-y-3">
            {matches.map((match) => (
              <MatchCard key={match.id} match={match} />
            ))}
          </div>

          {totalPages > 1 && (
            <div className="mt-6 flex items-center justify-center gap-2">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
                className="rounded-md border bg-card px-3 py-1 text-sm text-foreground disabled:opacity-50"
              >
                Previous
              </button>
              <span className="text-sm text-muted-foreground">
                Page {page} of {totalPages}
              </span>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
                className="rounded-md border bg-card px-3 py-1 text-sm text-foreground disabled:opacity-50"
              >
                Next
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
