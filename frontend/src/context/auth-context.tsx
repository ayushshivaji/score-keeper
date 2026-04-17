"use client";

import { createContext, useContext, useEffect, useState, useCallback, ReactNode } from "react";
import { User } from "@/types/user";
import { api } from "@/lib/api";

interface AuthContextType {
  user: User | null;
  loading: boolean;
  logout: () => Promise<void>;
  refresh: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: true,
  logout: async () => {},
  refresh: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchUser = useCallback(async () => {
    const res = await api.get<User>("/auth/me");
    if (res.data) {
      setUser(res.data);
    } else {
      setUser(null);
    }
    setLoading(false);
  }, []);

  useEffect(() => {
    fetchUser();
  }, [fetchUser]);

  const logout = async () => {
    await api.post("/auth/logout", {});
    setUser(null);
    window.location.href = "/";
  };

  const refresh = async () => {
    const res = await api.post("/auth/refresh", {});
    if (res.data) {
      await fetchUser();
    }
  };

  return (
    <AuthContext.Provider value={{ user, loading, logout, refresh }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
