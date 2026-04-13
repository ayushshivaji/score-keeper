# Matches API

## POST /api/v1/matches

Records a new table tennis match with set-by-set scores.

**Auth required:** Yes

**Description:** Creates a match record, validates all set scores against table tennis rules, determines the winner automatically, and updates both players' match statistics. All operations happen atomically in a single database transaction.

**Request Body:**

```json
{
  "player1_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "player2_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
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

**Request Fields:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `player1_id` | UUID | Yes | ID of player 1 |
| `player2_id` | UUID | Yes | ID of player 2 (must differ from player 1) |
| `match_format` | int | Yes | Best-of format: `3`, `5`, or `7` |
| `played_at` | datetime | Yes | When the match was played (ISO 8601) |
| `sets` | array | Yes | Array of set scores |
| `sets[].player1_score` | int | Yes | Player 1's score in this set (>= 0) |
| `sets[].player2_score` | int | Yes | Player 2's score in this set (>= 0) |

**Validation Rules:**

- `player1_id` and `player2_id` must be different
- `match_format` must be 3, 5, or 7
- Number of sets: minimum `ceil(format/2)`, maximum `format`
- Each set: winning score must be >= 21, winner must lead by >= 2
- In deuce (both >= 20): winner must win by exactly 2 (e.g., 22-20, 23-21)
- When opponent has < 20 points: winning score must be exactly 21
- Exactly one player must win `ceil(format/2)` sets
- No sets played after a player has clinched the match

**Success Response (201):**

```json
{
  "data": {
    "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
    "player1_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "player2_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "winner_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "match_format": 5,
    "player1_sets_won": 3,
    "player2_sets_won": 1,
    "tournament_match_id": null,
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
      "created_at": "2026-04-02T08:00:00Z",
      "updated_at": "2026-04-08T14:30:00Z"
    },
    "sets": [
      { "id": "...", "match_id": "...", "set_number": 1, "player1_score": 21, "player2_score": 17 },
      { "id": "...", "match_id": "...", "set_number": 2, "player1_score": 19, "player2_score": 21 },
      { "id": "...", "match_id": "...", "set_number": 3, "player1_score": 21, "player2_score": 15 },
      { "id": "...", "match_id": "...", "set_number": 4, "player1_score": 21, "player2_score": 18 }
    ]
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
    "message": "set 2: winning score must be at least 11"
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

Lists matches with optional filtering and pagination.

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
      "match_format": 5,
      "player1_sets_won": 3,
      "player2_sets_won": 1,
      "tournament_match_id": null,
      "played_at": "2026-04-08T14:00:00Z",
      "created_at": "2026-04-08T14:30:00Z",
      "created_by": "...",
      "player1": { "id": "...", "name": "Alice", "..." : "..." },
      "player2": { "id": "...", "name": "Bob", "..." : "..." },
      "sets": [
        { "set_number": 1, "player1_score": 11, "player2_score": 7 },
        { "set_number": 2, "player1_score": 9, "player2_score": 11 },
        { "set_number": 3, "player1_score": 11, "player2_score": 5 },
        { "set_number": 4, "player1_score": 11, "player2_score": 8 }
      ]
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

Returns full details for a single match including set scores and player info.

**Auth required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Match ID |

**Success Response (200):**

```json
{
  "data": {
    "id": "c3d4e5f6-a7b8-9012-cdef-123456789012",
    "player1_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "player2_id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
    "winner_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "match_format": 5,
    "player1_sets_won": 3,
    "player2_sets_won": 1,
    "tournament_match_id": null,
    "played_at": "2026-04-08T14:00:00Z",
    "created_at": "2026-04-08T14:30:00Z",
    "created_by": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "player1": {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "name": "Alice",
      "email": "alice@example.com",
      "avatar_url": "https://lh3.googleusercontent.com/a/photo",
      "matches_played": 16,
      "matches_won": 11
    },
    "player2": {
      "id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      "name": "Bob",
      "email": "bob@example.com",
      "avatar_url": null,
      "matches_played": 13,
      "matches_won": 5
    },
    "sets": [
      { "id": "...", "match_id": "...", "set_number": 1, "player1_score": 21, "player2_score": 17 },
      { "id": "...", "match_id": "...", "set_number": 2, "player1_score": 19, "player2_score": 21 },
      { "id": "...", "match_id": "...", "set_number": 3, "player1_score": 21, "player2_score": 15 },
      { "id": "...", "match_id": "...", "set_number": 4, "player1_score": 21, "player2_score": 18 }
    ]
  },
  "error": null
}
```

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

**Description:** Only the user who created the match can delete it, and only within 24 hours of creation. Deleting a match decrements both players' `matches_played` counters and the winner's `matches_won` counter. All operations are atomic.

**Constraints:**
- Only the match creator can delete
- Must be within 24 hours of creation
- Cascades deletion to `match_sets`

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
