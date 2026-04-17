"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { ReactNode } from "react";
import { useAuth } from "@/context/auth-context";
import { ThemeToggle } from "@/components/theme-toggle";

interface NavItem {
  href: string;
  label: string;
  icon: ReactNode;
}

const iconProps = {
  xmlns: "http://www.w3.org/2000/svg",
  width: 22,
  height: 22,
  viewBox: "0 0 24 24",
  fill: "none",
  stroke: "currentColor",
  strokeWidth: 2,
  strokeLinecap: "round" as const,
  strokeLinejoin: "round" as const,
  "aria-hidden": true,
};

const items: NavItem[] = [
  {
    href: "/dashboard",
    label: "Home",
    icon: (
      <svg {...iconProps}>
        <path d="M3 12 12 3l9 9" />
        <path d="M5 10v10a1 1 0 0 0 1 1h4v-6h4v6h4a1 1 0 0 0 1-1V10" />
      </svg>
    ),
  },
  {
    href: "/matches",
    label: "Matches",
    icon: (
      <svg {...iconProps}>
        <line x1="8" y1="6" x2="21" y2="6" />
        <line x1="8" y1="12" x2="21" y2="12" />
        <line x1="8" y1="18" x2="21" y2="18" />
        <line x1="3" y1="6" x2="3.01" y2="6" />
        <line x1="3" y1="12" x2="3.01" y2="12" />
        <line x1="3" y1="18" x2="3.01" y2="18" />
      </svg>
    ),
  },
  {
    href: "/matches/new",
    label: "Record",
    icon: (
      <svg {...iconProps}>
        <circle cx="12" cy="12" r="9" />
        <line x1="12" y1="8" x2="12" y2="16" />
        <line x1="8" y1="12" x2="16" y2="12" />
      </svg>
    ),
  },
  {
    href: "/players",
    label: "Players",
    icon: (
      <svg {...iconProps}>
        <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
        <circle cx="9" cy="7" r="4" />
        <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
        <path d="M16 3.13a4 4 0 0 1 0 7.75" />
      </svg>
    ),
  },
  {
    href: "/leaderboard",
    label: "Standings",
    icon: (
      <svg {...iconProps}>
        <path d="M8 21h8" />
        <path d="M12 17v4" />
        <path d="M7 4h10v6a5 5 0 0 1-10 0V4Z" />
        <path d="M17 4h3v2a3 3 0 0 1-3 3" />
        <path d="M7 4H4v2a3 3 0 0 0 3 3" />
      </svg>
    ),
  },
];

function isActive(pathname: string, href: string): boolean {
  if (href === "/dashboard") return pathname === "/dashboard";
  if (href === "/matches") {
    return pathname === "/matches" || (pathname.startsWith("/matches/") && !pathname.startsWith("/matches/new"));
  }
  if (href === "/matches/new") return pathname.startsWith("/matches/new");
  return pathname === href || pathname.startsWith(href + "/");
}

export function Sidebar() {
  const pathname = usePathname();
  const { user, logout } = useAuth();

  return (
    <aside className="fixed inset-y-0 left-0 z-20 flex w-24 flex-col items-center justify-between border-r bg-card py-6">
      {/* Brand + nav */}
      <div className="flex flex-col items-center gap-1">
        <Link
          href="/dashboard"
          className="mb-4 flex h-10 w-10 items-center justify-center rounded-xl bg-blue-600 text-sm font-bold text-white"
          aria-label="Score Keeper"
        >
          SK
        </Link>
        <nav className="flex flex-col items-center gap-1">
          {items.map((item) => {
            const active = isActive(pathname, item.href);
            return (
              <Link
                key={item.href}
                href={item.href}
                aria-current={active ? "page" : undefined}
                className="group flex flex-col items-center gap-1 px-2 py-2"
              >
                <span
                  className={`flex h-10 w-14 items-center justify-center rounded-2xl transition-colors ${
                    active
                      ? "bg-blue-500/20 text-blue-600 dark:bg-blue-400/20 dark:text-blue-300"
                      : "text-muted-foreground group-hover:bg-muted group-hover:text-foreground"
                  }`}
                >
                  {item.icon}
                </span>
                <span
                  className={`text-xs font-medium tracking-tight ${
                    active
                      ? "text-blue-600 dark:text-blue-300"
                      : "text-muted-foreground group-hover:text-foreground"
                  }`}
                >
                  {item.label}
                </span>
              </Link>
            );
          })}
        </nav>
      </div>

      {/* Footer: theme + user */}
      <div className="flex flex-col items-center gap-3">
        <ThemeToggle />
        {user && (
          <>
            {user.avatar_url ? (
              <img
                src={user.avatar_url}
                alt={user.name}
                title={user.name}
                className="h-9 w-9 rounded-full"
              />
            ) : (
              <div className="h-9 w-9 rounded-full bg-muted" />
            )}
            <button
              onClick={logout}
              className="text-[11px] font-medium text-muted-foreground hover:text-foreground"
            >
              Logout
            </button>
          </>
        )}
      </div>
    </aside>
  );
}
