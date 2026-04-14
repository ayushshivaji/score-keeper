# Users API

## GET /api/v1/users

Lists all registered users with optional search. Returns base `User` rows.

**Auth required:** Yes

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `search` | string | (none) | Filter users by name (case-insensitive partial match) |
| `page` | int | 1 | Page number |
| `per_page` | int | 20 | Items per page (max 50) |

**Success Response (200):**

```json
{
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "google_id": "108234567890123456789",
      "email": "alice@example.com",
      "name": "Alice",
      "avatar_url": "https://lh3.googleusercontent.com/a/photo",
      "matches_played": 15,
      "matches_won": 10,
      "total_points": 312,
      "created_at": "2026-04-01T10:00:00Z",
      "updated_at": "2026-04-08T14:30:00Z"
    }
  ],
  "error": null,
  "meta": { "page": 1, "per_page": 20, "total": 2 }
}
```

---

## GET /api/v1/users/:id

Returns a single user's profile with computed stats (streaks and recent form).

**Auth required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | User ID |

**Response fields beyond the base `User` row:**

| Field | Type | Description |
|-------|------|-------------|
| `losses` | int | `matches_played - matches_won` |
| `win_rate` | float | `matches_won / matches_played` as a decimal (0.0–1.0) |
| `current_streak` | int | Signed streak: positive = consecutive wins, negative = consecutive losses, 0 = no matches |
| `longest_win_streak` | int | Longest consecutive win run in their history |
| `longest_loss_streak` | int | Longest consecutive loss run in their history |
| `recent_form` | string[] | Last up-to-5 results as `"W"` / `"L"`, newest first |

**Success Response (200):**

```json
{
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "google_id": "108234567890123456789",
    "email": "alice@example.com",
    "name": "Alice",
    "avatar_url": "https://lh3.googleusercontent.com/a/photo",
    "matches_played": 15,
    "matches_won": 10,
    "total_points": 312,
    "created_at": "2026-04-01T10:00:00Z",
    "updated_at": "2026-04-08T14:30:00Z",
    "losses": 5,
    "win_rate": 0.6667,
    "current_streak": 3,
    "longest_win_streak": 5,
    "longest_loss_streak": 2,
    "recent_form": ["W", "W", "W", "L", "W"]
  },
  "error": null
}
```

**Error Response (404):**

```json
{
  "data": null,
  "error": { "code": "NOT_FOUND", "message": "user not found" }
}
```

---

## GET /api/v1/users/:id/head-to-head/:opponentId

Returns aggregated and per-match head-to-head history between two players.

**Auth required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | First player |
| `opponentId` | UUID | Second player (must differ from `id`) |

**Response fields:**

| Field | Type | Description |
|-------|------|-------------|
| `player1` | User | The player referenced by `:id` |
| `player2` | User | The player referenced by `:opponentId` |
| `total_matches` | int | Total matches played between them |
| `player1_wins` | int | How many of those matches `player1` won |
| `player2_wins` | int | How many `player2` won |
| `player1_points` | int | Points `player1` scored across these matches |
| `player2_points` | int | Points `player2` scored across these matches |
| `matches` | Match[] | All matches between them, newest first |

**Success Response (200):**

```json
{
  "data": {
    "player1": { "id": "...", "name": "Alice", "matches_played": 15, "matches_won": 10, "total_points": 312, "...": "..." },
    "player2": { "id": "...", "name": "Bob", "matches_played": 13, "matches_won": 5, "total_points": 198, "...": "..." },
    "total_matches": 7,
    "player1_wins": 5,
    "player2_wins": 2,
    "player1_points": 142,
    "player2_points": 98,
    "matches": [
      {
        "id": "...",
        "player1_id": "...",
        "player2_id": "...",
        "winner_id": "...",
        "player1_score": 21,
        "player2_score": 18,
        "played_at": "2026-04-08T14:00:00Z",
        "created_at": "2026-04-08T14:30:00Z",
        "created_by": "...",
        "player1": { "...": "..." },
        "player2": { "...": "..." }
      }
    ]
  },
  "error": null
}
```

**Error Response — same player (400):**

```json
{
  "data": null,
  "error": { "code": "BAD_REQUEST", "message": "players must be different" }
}
```

**Error Response — player not found (400):**

```json
{
  "data": null,
  "error": { "code": "BAD_REQUEST", "message": "player not found" }
}
```
