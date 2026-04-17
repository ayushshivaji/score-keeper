"use client";

import { useAuth } from "@/context/auth-context";
import { api } from "@/lib/api";
import { User } from "@/types/user";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

// Same-origin; goes through Next's rewrite to the backend.
const API_URL = "/api/v1";

export default function Home() {
  const { user, loading } = useAuth();
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (!loading && user) {
      router.replace("/dashboard");
    }
  }, [user, loading, router]);

  const handleStaticLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    const res = await api.post<User>("/auth/login", { username, password });
    setSubmitting(false);
    if (res.error) {
      setError(res.error.message || "Login failed");
      return;
    }
    window.location.href = "/dashboard";
  };

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-8 px-4">
      <div className="text-center">
        <h1 className="text-4xl font-bold text-foreground">Score Keeper</h1>
        <p className="mt-2 text-lg text-muted-foreground">
          Track your table tennis matches, stats, and rankings
        </p>
      </div>

      <a
        href={`${API_URL}/auth/google`}
        className="rounded-lg bg-blue-600 px-6 py-3 text-white font-medium hover:bg-blue-700 transition-colors"
      >
        Sign in with Google
      </a>

      <div className="flex items-center gap-4 w-full max-w-xs">
        <div className="h-px flex-1 bg-border" />
        <span className="text-xs uppercase text-muted-foreground">or</span>
        <div className="h-px flex-1 bg-border" />
      </div>

      <form
        onSubmit={handleStaticLogin}
        className="flex w-full max-w-xs flex-col gap-3"
      >
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
          autoComplete="username"
          className="rounded-lg border border-border bg-background px-4 py-2 text-foreground focus:outline-none focus:ring-2 focus:ring-blue-600"
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          autoComplete="current-password"
          className="rounded-lg border border-border bg-background px-4 py-2 text-foreground focus:outline-none focus:ring-2 focus:ring-blue-600"
        />
        {error && <p className="text-sm text-red-600">{error}</p>}
        <button
          type="submit"
          disabled={submitting}
          className="rounded-lg bg-foreground px-4 py-2 text-background font-medium hover:opacity-90 transition disabled:opacity-50"
        >
          {submitting ? "Signing in..." : "Sign in"}
        </button>
      </form>
    </div>
  );
}
