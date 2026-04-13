# Score Keeper — Table Tennis Score Tracking Application

## Design Document

A web application for tracking table tennis match scores, player statistics, ELO ratings, and tournaments.

---

## 1. System Architecture

### Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 15 (React 19, TypeScript, Tailwind CSS, shadcn/ui) |
| Backend | Go 1.22+ (Gin framework) |
| Database | PostgreSQL 16 |
| Auth | Google OAuth 2.0 + JWT |
| Deployment | Docker Compose |

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                      Docker Compose                         │
│                                                             │
│  ┌────────────┐    ┌────────────┐    ┌────────────────┐    │
│  │  Frontend   │    │  Backend   │    │  PostgreSQL    │    │
│  │  Next.js    │───▶│  Go (Gin)  │───▶│  Database      │    │
│  │  :3000      │    │  :8080     │    │  :5432         │    │
│  └────────────┘    └────────────┘    └────────────────┘    │
│        │                 │                                   │
│        │                 ▼                                   │
│        │          ┌─────────────┐                           │
│        └─────────▶│ Google OAuth│                           │
│                   │ (external)  │                           │
│                   └─────────────┘                           │
└─────────────────────────────────────────────────────────────┘
```

### Auth Flow

1. User clicks "Sign in with Google" on the frontend.
2. Frontend redirects to `GET /api/v1/auth/google` on the Go backend.
3. Backend redirects to Google OAuth consent screen.
4. Google redirects back to `GET /api/v1/auth/google/callback` with an authorization code.
5. Backend exchanges code for tokens, fetches user profile, upserts user in the database.
6. Backend creates a JWT access token (15 min expiry) and a refresh token (7-day expiry, stored in DB).
7. Backend redirects to frontend with tokens set as HTTP-only secure cookies.
8. Frontend includes cookies automatically on subsequent API requests.
9. Backend middleware validates JWT on protected routes; refresh endpoint issues new JWTs.

### API Convention

- All endpoints prefixed with `/api/v1`.
- In development, Next.js `rewrites` in `next.config.ts` proxy `/api/v1/*` to `http://backend:8080/api/v1/*`.
- All responses use a consistent envelope:

```json
{
  "data": { ... },
  "error": null,
  "meta": { "page": 1, "per_page": 20, "total": 150 }
}
```

Error responses:

```json
{
  "data": null,
  "error": { "code": "VALIDATION_ERROR", "message": "Player 1 and Player 2 must be different" },
  "meta": null
}
```

---

## 2. Database Schema

### Table: `users`

| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() |
| google_id | VARCHAR(255) | UNIQUE, NOT NULL |
| email | VARCHAR(255) | UNIQUE, NOT NULL |
| name | VARCHAR(255) | NOT NULL |
| avatar_url | TEXT | NULLABLE |
| elo_rating | INTEGER | NOT NULL, DEFAULT 1200 |
| elo_peak | INTEGER | NOT NULL, DEFAULT 1200 |
| matches_played | INTEGER | NOT NULL, DEFAULT 0 |
| matches_won | INTEGER | NOT NULL, DEFAULT 0 |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |
| updated_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Indexes:** `users(elo_rating DESC)` for leaderboard queries.

### Table: `refresh_tokens`

| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() |
| user_id | UUID | FK → users(id) ON DELETE CASCADE |
| token_hash | VARCHAR(255) | UNIQUE, NOT NULL |
| expires_at | TIMESTAMPTZ | NOT NULL |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

### Table: `matches`

| Column | Type | Constraints | Notes |
|--------|------|-------------|-------|
| id | UUID | PK, DEFAULT gen_random_uuid() | |
| player1_id | UUID | FK → users(id), NOT NULL | |
| player2_id | UUID | FK → users(id), NOT NULL | |
| winner_id | UUID | FK → users(id), NOT NULL | |
| match_format | SMALLINT | NOT NULL, CHECK (match_format IN (3, 5, 7)) | Best of 3, 5, or 7 |
| player1_sets_won | SMALLINT | NOT NULL | |
| player2_sets_won | SMALLINT | NOT NULL | |
| player1_elo_before | INTEGER | NOT NULL | Snapshot for history |
| player2_elo_before | INTEGER | NOT NULL | |
| player1_elo_after | INTEGER | NOT NULL | |
| player2_elo_after | INTEGER | NOT NULL | |
| tournament_match_id | UUID | FK → tournament_matches(id), NULLABLE | NULL if casual match |
| played_at | TIMESTAMPTZ | NOT NULL | When the match was played |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() | When it was recorded |
| created_by | UUID | FK → users(id), NOT NULL | Who entered the score |

**CHECK:** `player1_id <> player2_id`
**Indexes:** `matches(player1_id)`, `matches(player2_id)`, `matches(played_at DESC)`

### Table: `match_sets`

| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() |
| match_id | UUID | FK → matches(id) ON DELETE CASCADE |
| set_number | SMALLINT | NOT NULL, CHECK (set_number BETWEEN 1 AND 7) |
| player1_score | SMALLINT | NOT NULL, CHECK (player1_score >= 0) |
| player2_score | SMALLINT | NOT NULL, CHECK (player2_score >= 0) |

**UNIQUE:** `(match_id, set_number)`
**Validation (application-level):** The set winner must have >= 21 points and win by >= 2.

### Table: `tournaments`

| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() |
| name | VARCHAR(255) | NOT NULL |
| description | TEXT | NULLABLE |
| format | VARCHAR(50) | NOT NULL, CHECK (format IN ('single_elimination', 'double_elimination', 'round_robin')) |
| match_format | SMALLINT | NOT NULL, DEFAULT 5 |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'registration', CHECK (status IN ('registration', 'in_progress', 'completed', 'cancelled')) |
| max_players | INTEGER | NULLABLE (NULL = unlimited) |
| created_by | UUID | FK → users(id), NOT NULL |
| starts_at | TIMESTAMPTZ | NULLABLE |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |
| updated_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

### Table: `tournament_participants`

| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() |
| tournament_id | UUID | FK → tournaments(id) ON DELETE CASCADE |
| user_id | UUID | FK → users(id) |
| seed | INTEGER | NULLABLE |
| registered_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**UNIQUE:** `(tournament_id, user_id)`

### Table: `tournament_matches`

| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() |
| tournament_id | UUID | FK → tournaments(id) ON DELETE CASCADE |
| round | SMALLINT | NOT NULL |
| position | SMALLINT | NOT NULL |
| player1_id | UUID | FK → users(id), NULLABLE |
| player2_id | UUID | FK → users(id), NULLABLE |
| winner_id | UUID | FK → users(id), NULLABLE |
| match_id | UUID | FK → matches(id), NULLABLE |
| bracket_type | VARCHAR(20) | DEFAULT 'winners', CHECK (bracket_type IN ('winners', 'losers', 'grand_final')) |
| status | VARCHAR(20) | NOT NULL, DEFAULT 'pending', CHECK (status IN ('pending', 'ready', 'completed', 'bye')) |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**UNIQUE:** `(tournament_id, round, position, bracket_type)`
**Indexes:** `tournament_matches(tournament_id, round, position)`

### Table: `elo_history`

| Column | Type | Constraints |
|--------|------|-------------|
| id | UUID | PK, DEFAULT gen_random_uuid() |
| user_id | UUID | FK → users(id) ON DELETE CASCADE |
| match_id | UUID | FK → matches(id) ON DELETE CASCADE |
| rating_before | INTEGER | NOT NULL |
| rating_after | INTEGER | NOT NULL |
| created_at | TIMESTAMPTZ | NOT NULL, DEFAULT NOW() |

**Indexes:** `(user_id, created_at)` for rating chart queries.

---

## 3. API Endpoints

All endpoints prefixed with `/api/v1`. Protected endpoints require a valid JWT.

### Auth

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /auth/google | No | Redirect to Google OAuth consent |
| GET | /auth/google/callback | No | OAuth callback, set cookies, redirect to frontend |
| POST | /auth/refresh | No (cookie) | Refresh access token |
| POST | /auth/logout | Yes | Invalidate refresh token, clear cookies |
| GET | /auth/me | Yes | Return current authenticated user |

### Users

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /users | Yes | List users (search by name, paginated) |
| GET | /users/:id | Yes | Get user profile with stats |
| GET | /users/:id/matches | Yes | Get user's match history (paginated, filterable) |
| GET | /users/:id/elo-history | Yes | Get user's ELO rating over time |
| GET | /users/:id/head-to-head/:opponentId | Yes | Head-to-head stats between two players |

### Matches

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /matches | Yes | Create a new match with set scores |
| GET | /matches | Yes | List matches (paginated, filterable by player, date) |
| GET | /matches/:id | Yes | Get match details with all set scores |
| DELETE | /matches/:id | Yes | Delete match (creator only, within 24h, recalculates ELO) |

**POST /matches — Request Body:**

```json
{
  "player1_id": "uuid",
  "player2_id": "uuid",
  "match_format": 5,
  "played_at": "2026-04-08T14:00:00Z",
  "sets": [
    { "player1_score": 21, "player2_score": 17 },
    { "player1_score": 19, "player2_score": 21 },
    { "player1_score": 21, "player2_score": 15 },
    { "player1_score": 21, "player2_score": 18 }
  ]
}
```

**Validation rules:**
- Number of sets: minimum `ceil(format/2)`, maximum `format`.
- Each set: winner must have >= 21 points and win by >= 2.
- Exactly one player wins `ceil(format/2)` sets.
- No sets played after a player has clinched the match.
- Winner is auto-determined from set scores.

### Leaderboard

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | /leaderboard | Yes | Ranked players by ELO (paginated, with rank numbers) |
| GET | /leaderboard/stats | Yes | Global stats (total matches, players, most active) |

### Tournaments

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /tournaments | Yes | Create a tournament |
| GET | /tournaments | Yes | List tournaments (filterable by status) |
| GET | /tournaments/:id | Yes | Get tournament details with bracket |
| PUT | /tournaments/:id | Yes | Update tournament (organizer only) |
| POST | /tournaments/:id/join | Yes | Join tournament (during registration) |
| DELETE | /tournaments/:id/leave | Yes | Leave tournament (during registration) |
| POST | /tournaments/:id/start | Yes | Start tournament, generate bracket (organizer only) |
| POST | /tournaments/:id/matches/:matchId/result | Yes | Submit match result within tournament |
| GET | /tournaments/:id/bracket | Yes | Get bracket data for rendering |

**GET /tournaments/:id/bracket — Response:**

```json
{
  "format": "single_elimination",
  "rounds": [
    {
      "round": 1,
      "name": "Quarter-Finals",
      "matches": [
        {
          "id": "uuid",
          "position": 1,
          "player1": { "id": "uuid", "name": "Alice", "seed": 1 },
          "player2": { "id": "uuid", "name": "Bob", "seed": 8 },
          "winner_id": "uuid",
          "score": "3-1",
          "status": "completed"
        }
      ]
    }
  ]
}
```

---

## 4. UI Pages

### 4.1 `/` — Landing Page

- If not authenticated: hero section with app description, "Sign in with Google" button.
- If authenticated: redirect to `/dashboard`.

### 4.2 `/dashboard` — Personal Dashboard

- Welcome message with user's name and avatar.
- Quick stats cards: current ELO, rank, W/L record, win percentage, current streak.
- ELO rating chart over time (line chart using recharts).
- Recent matches list (last 5).
- Quick action buttons: "Record a Match", "View Leaderboard".
- Upcoming tournament matches (if any).

### 4.3 `/matches/new` — Record a Match

- **Player selection:** two searchable comboboxes for Player 1 and Player 2. Current user pre-selected as Player 1.
- **Match format selector:** best of 3 / 5 / 7 radio buttons.
- **Date/time picker:** defaults to now.
- **Set entry area:** for each set, two **slider inputs** (range 0–40, typical range shown 0–25) for each player's score. Number inputs beside sliders for precise entry.
- **Dynamic sets:** start with set 1 visible. Show next set when current set is filled. Stop when a player clinches the match.
- **Real-time validation:** highlight invalid set scores (winner < 21, margin < 2, etc.).
- **Match summary preview** before submission showing who won and the final set score line.
- **Submit button.**

### 4.4 `/matches` — Match History

- Paginated card list of all matches.
- Filter controls: player name search, date range, match format.
- Each card shows: Player 1 vs Player 2, set scores (e.g., 3-1: 21-17, 19-21, 21-15, 21-18), date, ELO changes (+15 / -15).
- Click to navigate to match detail.

### 4.5 `/matches/:id` — Match Detail

- Full match info: players with avatars, date, match format.
- Set-by-set score breakdown in a visual table.
- ELO change for each player (before → after).
- Link to head-to-head between the two players.

### 4.6 `/players` — Player Directory

- Searchable list of all players.
- Each card: avatar, name, ELO, W/L record.
- Click to navigate to player profile.

### 4.7 `/players/:id` — Player Profile

- Header: avatar, name, ELO rating, rank badge.
- Stats: matches played, wins, losses, win rate, average sets per match, longest win streak.
- ELO chart over time.
- Recent matches (paginated).
- Frequent opponents list with head-to-head records.

### 4.8 `/players/:id/vs/:opponentId` — Head-to-Head

- Side-by-side player comparison.
- Overall record (e.g., Player A leads 7-3).
- Set record (Player A won 18 sets vs 12).
- Point differential.
- All matches between them with results.
- Dual ELO chart showing both players' ratings over time.

### 4.9 `/leaderboard` — Leaderboard

- Ranked table by ELO.
- Columns: rank, player (avatar + name), ELO, W/L, win%, matches played, trend (up/down arrow with recent change).
- Highlight current user's row.
- Time range toggle: All Time / Last 30 Days / Last 7 Days.

### 4.10 `/tournaments` — Tournament List

- Tabs: Upcoming (registration), In Progress, Completed.
- Each card: name, format, participant count, status, start date.
- "Create Tournament" button.

### 4.11 `/tournaments/new` — Create Tournament

- Form: name, description, format (single/double elimination, round robin), match format (best of N), max players, start date.
- Submit creates tournament in "registration" status.

### 4.12 `/tournaments/:id` — Tournament Detail

- Header: name, format, status, organizer, dates.
- Participants list with seeds.
- **Registration phase:** "Join" / "Leave" buttons. Organizer sees "Start Tournament" (enabled when >= 2 participants).
- **In progress / completed:** interactive bracket visualization.
  - Single elimination: standard bracket tree.
  - Double elimination: winners bracket + losers bracket + grand final.
  - Round robin: grid/matrix showing all matchups and results.
- Click on a bracket match to enter/view result.

### 4.13 `/tournaments/:id/matches/:matchId` — Tournament Match Entry

- Same UI as `/matches/new` but pre-populated with the two players from the bracket.
- On submission, updates bracket and advances winner.

### Shared Components

| Component | Description |
|-----------|-------------|
| `Navbar` | Logo, nav links (Dashboard, Matches, Leaderboard, Tournaments), user avatar dropdown (profile, logout) |
| `ScoreSlider` | Range slider (0–30) + number input for a single player's set score. Visual feedback for valid/invalid |
| `MatchCard` | Compact match result display (players, score, date, ELO delta) |
| `MatchForm` | Full match entry form with dynamic sets and validation |
| `PlayerSearchCombobox` | Searchable dropdown for selecting a player |
| `EloChart` | Recharts line chart for rating history |
| `BracketView` | Tournament bracket tree visualization |
| `RoundRobinGrid` | Round robin matrix view |
| `LeaderboardTable` | Ranked table with sorting |
| `StatsCard` | Stat display card (label, value, trend) |
| `Pagination` | Reusable pagination controls |

---

## 5. ELO Rating Algorithm

### Formula

**Expected score:**

```
E_A = 1 / (1 + 10^((R_B - R_A) / 400))
E_B = 1 - E_A
```

**New rating:**

```
R_A_new = R_A + K * (S_A - E_A)
R_B_new = R_B + K * (S_B - E_B)
```

Where:
- `R_A`, `R_B` = current ratings
- `S_A` = 1 if A wins, 0 if A loses (no draws in table tennis)
- `K` = K-factor (see below)

### K-Factor (Tiered)

| Condition | K | Rationale |
|-----------|---|-----------|
| Player has < 30 matches | 40 | New player, rating converges quickly |
| 30+ matches, ELO < 2000 | 24 | Intermediate, moderate adjustments |
| 30+ matches, ELO >= 2000 | 16 | Established, stable rating |

Each player uses their own K-factor in a match.

### Initial Rating

All new players start at **1200**.

### Margin of Victory Modifier (Phase 4 enhancement)

```
MOV = ln(set_difference + 1) * (2.2 / ((elo_difference * 0.001) + 2.2))
```

Multiply K by MOV. Gives more rating change for dominant wins (3-0) vs narrow ones (3-2), dampened when higher-rated player dominates.

### Implementation Notes

- ELO is calculated atomically during match creation inside a DB transaction.
- `elo_history` records every change for charting.
- Denormalized `matches_played` and `matches_won` on `users` are updated in the same transaction.
- `elo_peak` is updated if new rating exceeds it.
- Match deletion (within 24h, by creator): revert ELO using `elo_before`/`elo_after` snapshots. If subsequent matches exist for either player, refuse deletion or run a full re-computation chain.

---

## 6. Tournament System

### Single Elimination

1. Sort participants by seed (ELO descending or random).
2. Bracket size = next power of 2 >= participant count.
3. BYEs = bracket size - participant count. Higher seeds get BYEs.
4. Standard seeding: seed 1 vs seed N, seed 2 vs seed N-1, etc.
5. When match result is submitted, winner advances to the next round's match.

### Double Elimination

- Initial "winners" bracket same as single elimination.
- Losers drop to a "losers" bracket.
- Losers bracket matches generated dynamically as results arrive.
- Grand final: winners bracket champion vs losers bracket champion.
- If losers bracket champion wins the grand final, a reset match is played.

### Round Robin

- Every participant plays every other participant once.
- All `tournament_matches` pre-generated at start.
- Standings determined by: (1) wins, (2) head-to-head, (3) set differential, (4) point differential.

---

## 7. Docker Compose Setup

### docker-compose.yml

```yaml
version: "3.9"
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: scorekeeper
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: scorekeeper
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U scorekeeper"]
      interval: 5s
      timeout: 5s
      retries: 5

  migrate:
    build:
      context: ./backend
      dockerfile: Dockerfile.migrate
    depends_on:
      db:
        condition: service_healthy
    environment:
      DATABASE_URL: postgres://scorekeeper:${DB_PASSWORD}@db:5432/scorekeeper?sslmode=disable

  backend:
    build:
      context: ./backend
    depends_on:
      migrate:
        condition: service_completed_successfully
    environment:
      DATABASE_URL: postgres://scorekeeper:${DB_PASSWORD}@db:5432/scorekeeper?sslmode=disable
      GOOGLE_CLIENT_ID: ${GOOGLE_CLIENT_ID}
      GOOGLE_CLIENT_SECRET: ${GOOGLE_CLIENT_SECRET}
      JWT_SECRET: ${JWT_SECRET}
      FRONTEND_URL: http://localhost:3000
    ports:
      - "8080:8080"

  frontend:
    build:
      context: ./frontend
    depends_on:
      - backend
    environment:
      NEXT_PUBLIC_API_URL: http://localhost:8080/api/v1
    ports:
      - "3000:3000"

volumes:
  pgdata:
```

### Environment Variables (`.env.example`)

```
DB_PASSWORD=changeme
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
JWT_SECRET=a-random-64-char-secret
```

---

## 8. Directory Structure

### Backend (`/backend`)

```
backend/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point
├── internal/
│   ├── config/
│   │   └── config.go                # Env var loading
│   ├── middleware/
│   │   ├── auth.go                  # JWT validation
│   │   ├── cors.go                  # CORS config
│   │   └── logging.go              # Request logging
│   ├── handler/
│   │   ├── auth.go                  # Google OAuth, refresh, logout
│   │   ├── user.go                  # User profile & stats
│   │   ├── match.go                 # Match CRUD
│   │   ├── leaderboard.go          # Leaderboard
│   │   └── tournament.go           # Tournament management
│   ├── service/
│   │   ├── auth.go                  # OAuth exchange, JWT creation
│   │   ├── user.go                  # User business logic
│   │   ├── match.go                 # Match creation, validation, ELO
│   │   ├── elo.go                   # ELO algorithm
│   │   ├── leaderboard.go          # Leaderboard queries
│   │   └── tournament.go           # Bracket gen, advancement
│   ├── repository/
│   │   ├── user.go                  # User DB queries
│   │   ├── match.go                 # Match DB queries
│   │   ├── tournament.go           # Tournament DB queries
│   │   └── elo_history.go          # ELO history DB queries
│   ├── model/
│   │   ├── user.go                  # User struct
│   │   ├── match.go                 # Match, MatchSet structs
│   │   ├── tournament.go           # Tournament structs
│   │   └── elo.go                   # ELO history struct
│   ├── dto/
│   │   ├── request.go               # Request DTOs
│   │   └── response.go             # Response DTOs, pagination
│   ├── router/
│   │   └── router.go               # Route registration
│   └── validator/
│       └── match.go                 # Table tennis score validation
├── migrations/
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_matches.up.sql
│   ├── 000002_create_matches.down.sql
│   ├── 000003_create_tournaments.up.sql
│   ├── 000003_create_tournaments.down.sql
│   ├── 000004_create_elo_history.up.sql
│   └── 000004_create_elo_history.down.sql
├── Dockerfile
├── Dockerfile.migrate
├── go.mod
└── go.sum
```

### Frontend (`/frontend`)

```
frontend/
├── src/
│   ├── app/
│   │   ├── layout.tsx               # Root layout (providers, navbar)
│   │   ├── page.tsx                 # Landing page
│   │   ├── globals.css
│   │   ├── (auth)/
│   │   │   └── login/
│   │   │       └── page.tsx         # Login page
│   │   └── (protected)/
│   │       ├── layout.tsx           # Auth-guarded layout
│   │       ├── dashboard/
│   │       │   └── page.tsx
│   │       ├── matches/
│   │       │   ├── page.tsx         # Match history
│   │       │   ├── new/
│   │       │   │   └── page.tsx     # Record a match
│   │       │   └── [id]/
│   │       │       └── page.tsx     # Match detail
│   │       ├── players/
│   │       │   ├── page.tsx         # Player directory
│   │       │   └── [id]/
│   │       │       ├── page.tsx     # Player profile
│   │       │       └── vs/
│   │       │           └── [opponentId]/
│   │       │               └── page.tsx  # Head-to-head
│   │       ├── leaderboard/
│   │       │   └── page.tsx
│   │       └── tournaments/
│   │           ├── page.tsx         # Tournament list
│   │           ├── new/
│   │           │   └── page.tsx     # Create tournament
│   │           └── [id]/
│   │               └── page.tsx     # Tournament detail + bracket
│   ├── components/
│   │   ├── ui/                      # shadcn/ui base components
│   │   ├── navbar.tsx
│   │   ├── score-slider.tsx
│   │   ├── match-card.tsx
│   │   ├── match-form.tsx
│   │   ├── player-search.tsx
│   │   ├── elo-chart.tsx
│   │   ├── bracket-view.tsx
│   │   ├── round-robin-grid.tsx
│   │   ├── leaderboard-table.tsx
│   │   ├── stats-card.tsx
│   │   └── pagination.tsx
│   ├── lib/
│   │   ├── api.ts                   # Fetch wrapper with auth
│   │   ├── utils.ts
│   │   └── constants.ts
│   ├── hooks/
│   │   ├── use-auth.ts
│   │   ├── use-matches.ts
│   │   ├── use-players.ts
│   │   ├── use-leaderboard.ts
│   │   └── use-tournaments.ts
│   ├── types/
│   │   ├── user.ts
│   │   ├── match.ts
│   │   ├── tournament.ts
│   │   └── api.ts
│   └── context/
│       └── auth-context.tsx
├── public/
├── next.config.ts
├── tailwind.config.ts
├── tsconfig.json
├── package.json
└── Dockerfile
```

---

## 9. Implementation Phases

### Phase 1: Foundation — Auth + Match Recording

**Goal:** Users can sign in, record matches, and view history. No ELO calculations yet.

1. Project scaffolding: Go module, Next.js app, Docker Compose, PostgreSQL.
2. DB migrations: `users`, `refresh_tokens`, `matches`, `match_sets`.
3. Backend: config loading, DB connection, Gin router, CORS, logging middleware.
4. Backend: Google OAuth flow (handlers, JWT middleware, refresh tokens).
5. Backend: Match CRUD with table tennis score validation.
6. Frontend: Next.js setup, Tailwind + shadcn/ui, auth context.
7. Frontend: Login page, OAuth redirect, protected route layout.
8. Frontend: Match entry form with score sliders.
9. Frontend: Match history + match detail pages.

### Phase 2: Profiles + Leaderboard

**Goal:** Player profiles with stats, head-to-head, and leaderboard (win-based, no ELO yet).

1. Backend: User profile with aggregated stats (wins, losses, streaks).
2. Backend: Head-to-head endpoint.
3. Backend: Leaderboard with ranking by win count/win rate.
4. Frontend: Dashboard with stats cards.
5. Frontend: Player directory and profile pages.
6. Frontend: Head-to-head comparison page.
7. Frontend: Leaderboard page.

### Phase 3: Tournaments

**Goal:** Full tournament management with bracket visualization.

1. Backend: Tournament CRUD, registration, start logic.
2. Backend: Bracket generation (single elimination with BYEs).
3. Backend: Match result submission with bracket advancement.
4. Backend: Double elimination and round robin support.
5. Frontend: Tournament list, create form, detail page.
6. Frontend: Bracket visualization and round robin grid.
7. Frontend: Tournament match entry.

### Phase 4: ELO Ratings + Polish

**Goal:** Add ELO rating system and production polish.

1. DB migration: `elo_history` table, add ELO columns to `users`.
2. Backend: ELO calculation service (tiered K-factor).
3. Backend: Backfill ELO from existing match history (process matches in chronological order).
4. Backend: ELO history endpoint for charting.
5. Backend: Update leaderboard to rank by ELO instead of wins.
6. Frontend: ELO chart on dashboard and player profiles.
7. Frontend: ELO deltas on match cards and match detail.
8. Margin-of-victory ELO modifier.
9. Match deletion with ELO recalculation.
10. Rate limiting on API endpoints.
11. Loading states, error boundaries, empty states.
12. Mobile-responsive design pass.
13. Production Docker builds (multi-stage, optimized).

---

## 10. Key Go Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/gin-gonic/gin` | HTTP router |
| `github.com/jackc/pgx/v5` | PostgreSQL driver |
| `github.com/golang-migrate/migrate/v4` | DB migrations |
| `github.com/golang-jwt/jwt/v5` | JWT creation/validation |
| `github.com/google/uuid` | UUID generation |
| `golang.org/x/oauth2` | Google OAuth client |
| `github.com/rs/zerolog` | Structured logging |

## 11. Key Frontend Dependencies

| Package | Purpose |
|---------|---------|
| `@tanstack/react-query` | Server state management |
| `recharts` | Charts (ELO history) |
| `shadcn/ui` + `tailwindcss` | UI components + styling |
| `date-fns` | Date formatting |
| `cmdk` | Command-style combobox for player search |
