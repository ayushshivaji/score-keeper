# Matches API

Each match is a single standalone game (not a best-of-N series). The match is recorded as one pair of final scores.

## POST /api/v1/matches

Records a new table tennis match.

**Auth required:** Yes

**Description:** Creates a match record, validates the score against table tennis rules, determines the winner automatically, and updates both players' aggregate stats (`matches_played`, `matches_won`, `total_points`). All operations happen atomically in a single database transaction.

**Request Body:**

```json
{
  "player1_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "player2_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
  "player1_score": 21,
  "player2_score": 18,
  "played_at": "2026-04-08T14:00:00Z"
}
```

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `player1_id` | UUID | Yes | ID of player 1 |
| `player2_id` | UUID | Yes | ID of player 2 (must differ from player 1) |
| `player1_score` | int | Yes | Player 1's final score (>= 0) |
| `player2_score` | int | Yes | Player 2's final score (>= 0) |
| `played_at` | datetime | Yes | When the match was played (ISO 8601) |

**Validation Rules:**

- `player1_id` and `player2_id` must be different
- Both scores must be non-negative and not tied
- Winning score must be >= 21
- Winner must lead by >= 2
- In deuce (both >= 20): winner must win by exactly 2 (e.g., 22-20, 23-21)
- When opponent has < 20 points: winning score must be exactly 21

**Success Response (201):**

```json
{
  "data": {
    "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
    "player1_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "player2_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "winner_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "player1_score": 21,
    "player2_score": 18,
    "played_at": "2026-04-08T14:00:00Z",
    "created_at": "2026-04-08T14:30:00Z",
    "created_by": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "player1": {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "google_id": "108234567890123456789",
      "email": "alice@example.com",
      "name": "Alice",
      "avatar_url": "https://lh3.googleusercontent.com/a/photo",
      "matches_played": 16,
      "matches_won": 11,
      "total_points": 312,
      "created_at": "2026-04-01T10:00:00Z",
      "updated_at": "2026-04-08T14:30:00Z"
    },
    "player2": {
      "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "google_id": "109876543210987654321",
      "email": "bob@example.com",
      "name": "Bob",
      "avatar_url": null,
      "matches_played": 13,
      "matches_won": 5,
      "total_points": 198,
      "created_at": "2026-04-02T08:00:00Z",
      "updated_at": "2026-04-08T14:30:00Z"
    }
  },
  "error": null
}
```

**Error Response — validation failure (400):**

```json
{
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "winner must win by at least 2 points"
  }
}
```

**Error Response — same player (400):**

```json
{
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "player 1 and player 2 must be different"
  }
}
```

---

## GET /api/v1/matches

Lists matches with optional filtering and pagination. Sorted by `played_at` descending.

**Auth required:** Yes

**Query Parameters:**

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `player_id` | UUID | (none) | Filter to matches involving this player |
| `page` | int | 1 | Page number |
| `per_page` | int | 20 | Items per page (max 50) |

**Success Response (200):**

```json
{
  "data": [
    {
      "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
      "player1_id": "...",
      "player2_id": "...",
      "winner_id": "...",
      "player1_score": 21,
      "player2_score": 18,
      "played_at": "2026-04-08T14:00:00Z",
      "created_at": "2026-04-08T14:30:00Z",
      "created_by": "...",
      "player1": { "id": "...", "name": "Alice", "total_points": 312, "...": "..." },
      "player2": { "id": "...", "name": "Bob", "total_points": 198, "...": "..." }
    }
  ],
  "error": null,
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 42
  }
}
```

---

## GET /api/v1/matches/:id

Returns full details for a single match.

**Auth required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Match ID |

**Success Response (200):** Same shape as a single item in `GET /matches`.

**Error Response (404):**

```json
{
  "data": null,
  "error": {
    "code": "NOT_FOUND",
    "message": "match not found"
  }
}
```

---

## DELETE /api/v1/matches/:id

Deletes a match and reverts player statistics.

**Auth required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Match ID |

**Description:** Only the user who created the match can delete it, and only within 24 hours of creation. Deleting a match decrements both players' `matches_played` counters, the winner's `matches_won` counter, and each player's `total_points` by the points they scored in that match. All operations are atomic.

**Constraints:**
- Only the match creator can delete
- Must be within 24 hours of creation

**Success Response (200):**

```json
{
  "data": {
    "message": "match deleted"
  },
  "error": null
}
```

**Error Response — not creator (400):**

```json
{
  "data": null,
  "error": {
    "code": "BAD_REQUEST",
    "message": "only the creator can delete a match"
  }
}
```

**Error Response — too old (400):**

```json
{
  "data": null,
  "error": {
    "code": "BAD_REQUEST",
    "message": "matches can only be deleted within 24 hours of creation"
  }
}
```
