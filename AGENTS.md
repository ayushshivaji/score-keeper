# Score Keeper — Agent Context

Quick orientation for future agents. Read this first before making auth or plumbing changes.

## Stack
- **Backend**: Go + Gin, pgx/pgxpool against Postgres 16. Entry point `backend/cmd/server/main.go`. Layered: `config → repository → service → handler → router`.
- **Frontend**: Next.js (non-standard fork — see `frontend/AGENTS.md`; always check `node_modules/next/dist/docs/` before using Next APIs you aren't sure about). App Router under `frontend/src/app/`.
- **Infra**: `docker-compose.yml` runs `db → migrate → backend → frontend`. Backend has no published port; frontend proxies `/api/v1/*` to `backend:8080` via `next.config.ts` rewrites driven by `BACKEND_URL`. Cookies stay same-origin — no CORS between browser and backend.

## Auth — how it actually works
- **Two login paths, both issue the same cookie pair**:
  - Google OAuth: `GET /api/v1/auth/google` → Google consent → `GET /api/v1/auth/google/callback` → upsert by `google_id` → set cookies → 302 to `/dashboard`.
  - Static login: `POST /api/v1/auth/login` with `{username, password}` → validated against `STATIC_LOGIN_USERNAME` / `STATIC_LOGIN_PASSWORD` env vars via `crypto/subtle.ConstantTimeCompare` → upserts a synthetic user with `google_id = "static:<username>"`, `email = "<username>@static.local"` → set cookies → 200 JSON.
- **Why the synthetic `google_id`**: the `users` table has `google_id VARCHAR(255) UNIQUE NOT NULL` (migration `000001_create_users.up.sql`). Rather than a schema migration, static logins reuse `UpsertByGoogleID` with a `static:` prefix. Everything downstream (matches, leaderboard, refresh tokens) just sees a regular user row.
- **Tokens**: `AuthService.GenerateAccessToken` issues a 15-min HS256 JWT with `user_id` claim. `GenerateRefreshToken` issues a 32-byte random hex token; only its SHA-256 hash is stored in `refresh_tokens` with a 7-day expiry. Refresh rotates (delete old, insert new) in `RefreshAccessToken`.
- **Cookies**: `access_token` (900s) and `refresh_token` (604800s), both `HttpOnly`, path `/`, `Secure=false` (dev). If you change domains or go HTTPS, flip `Secure` and reconsider `SameSite`.
- **Middleware**: `middleware.AuthRequired` reads `access_token` cookie first, falls back to `Authorization: Bearer`, validates, stashes `user_id` (uuid.UUID) in the Gin context. Protected handlers call `c.MustGet("user_id").(uuid.UUID)`.
- **Static login is optional**: if either env var is empty, `LoginStatic` returns `ErrStaticLoginDisabled` → 404. Google OAuth still works.

## Key files for auth changes
- `backend/internal/config/config.go` — env-var loading. All auth-related knobs live here.
- `backend/internal/service/auth.go` — token issuance, validation, refresh rotation, static cred check (`LoginStatic`), sentinel errors `ErrStaticLoginDisabled` / `ErrInvalidCredentials`.
- `backend/internal/handler/auth.go` — Gin handlers: `GoogleLogin`, `GoogleCallback`, `StaticLogin`, `Refresh`, `Logout`, `Me`. Cookie-setting lives only here.
- `backend/internal/router/router.go` — route table. Public auth routes are under `auth := v1.Group("/auth")`; everything else goes through `middleware.AuthRequired`.
- `backend/internal/repository/user.go` — `UpsertByGoogleID`, `GetByID`, refresh-token CRUD. `userCols` constant drives `scanUser`; keep in sync with `model.User`.
- `backend/internal/model/user.go` — user struct. Note `GoogleID` is a plain `string` (not nullable) because the DB column is `NOT NULL`.
- `backend/migrations/000001_create_users.up.sql` — `users` + `refresh_tokens` schema. `google_id` and `email` are both `UNIQUE NOT NULL`.
- `frontend/src/app/page.tsx` — login UI: Google button + static-login form. Posts to `/auth/login` via `api.post`, then `window.location.href = "/dashboard"` so the `AuthProvider` re-fetches `/auth/me` on a fresh page load.
- `frontend/src/context/auth-context.tsx` — holds the `user` state, calls `/auth/me` on mount, exposes `logout` and `refresh`.
- `frontend/src/lib/api.ts` — always `credentials: "include"`, base `/api/v1`. Don't call the backend directly from the browser — use this client.
- `docker-compose.yml` — backend env block. Passes `GOOGLE_*`, `JWT_SECRET`, `STATIC_LOGIN_USERNAME`, `STATIC_LOGIN_PASSWORD`. Static vars default to empty (= disabled).

## Environment variables (backend)
| Var | Required | Notes |
|---|---|---|
| `DATABASE_URL` | yes | pgx DSN |
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` | yes | OAuth app |
| `GOOGLE_REDIRECT_URL` | no | Defaults to `http://localhost:8080/api/v1/auth/google/callback`; in compose it's overridden to the frontend origin so `Set-Cookie` lands on the same host |
| `JWT_SECRET` | yes | HS256 signing key |
| `FRONTEND_URL` | no | Used for the post-OAuth redirect and CORS |
| `STATIC_LOGIN_USERNAME` | no | Leave unset to disable the `/auth/login` form |
| `STATIC_LOGIN_PASSWORD` | no | Same |
| `PORT` | no | Default 8080 |

## Gotchas
- **Don't drop the `static:` prefix** on the synthetic `google_id`. If you ever let user-supplied values collide with real Google `sub` IDs, a static user could squat an OAuth identity.
- **Don't compare credentials with `==`**. Use `crypto/subtle.ConstantTimeCompare` — pattern is already in `LoginStatic`.
- **Static login doesn't use bcrypt**. The password lives in env, not the DB, and is compared in memory. If you add per-user static accounts, you must introduce a `password_hash` column and bcrypt/argon2.
- **OAuth callback host**: compose sets `GOOGLE_REDIRECT_URL` to the *frontend* origin on purpose (`docker-compose.yml` has a comment). The frontend's Next rewrite forwards the callback to the backend so the `Set-Cookie` lands on the browser-visible host. Don't "fix" it to point at the backend directly.
- **Frontend is a non-standard Next fork**. `frontend/AGENTS.md` warns that APIs may differ from training data. When editing frontend code, consult `node_modules/next/dist/docs/` rather than assuming vanilla Next behavior.
- **Refresh token storage is hashed**. Never log the raw token; never try to look up by plaintext.
- **`UpsertByGoogleID` also updates `email`, `name`, `avatar_url`**. For static users these are derived from the username on every login — changing `STATIC_LOGIN_USERNAME` between restarts will create a new row (new `google_id`), orphaning the old one's match history. Document this if it ever matters.

## Running
```sh
# from repo root
docker compose up --build
# frontend on :3000, backend only reachable through the frontend proxy
```

Static login quick test:
```sh
STATIC_LOGIN_USERNAME=admin STATIC_LOGIN_PASSWORD=secret docker compose up --build
# then POST {"username":"admin","password":"secret"} to http://localhost:3000/api/v1/auth/login
```
