# Leaderboard API

## GET /api/v1/leaderboard

Returns the standings: all players who have played at least one match, ranked by wins.

**Auth required:** Yes

**Description:** Players are sorted by `matches_won` descending, then by `total_points` descending, then by win rate (`matches_won / matches_played`) descending, then by `name` ascending. Players with zero matches played are excluded. Rank is implicit in the returned order (first item is rank 1).

`total_points` is the sum of points the player has scored across all their recorded matches (their side of each `player1_score` / `player2_score`).

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number (1-indexed) |
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
    },
    {
      "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "google_id": "109876543210987654321",
      "email": "bob@example.com",
      "name": "Bob",
      "avatar_url": null,
      "matches_played": 12,
      "matches_won": 5,
      "total_points": 198,
      "created_at": "2026-04-02T08:00:00Z",
      "updated_at": "2026-04-08T12:00:00Z"
    }
  ],
  "error": null,
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 2
  }
}
```

**Notes:**

- `losses` and `win_rate` are not returned by the API — clients derive them from `matches_played` and `matches_won`. See `computeStandingsRow` in `frontend/src/lib/utils.ts` for the reference implementation.
- `matches_played` / `matches_won` on each user are denormalized counters maintained atomically by `POST /api/v1/matches` and `DELETE /api/v1/matches/:id`, so the leaderboard reflects the live state immediately after a match is recorded or deleted.
