# Users API

## GET /api/v1/users

Lists all registered users with optional search.

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

---

## GET /api/v1/users/:id

Returns a single user's profile.

**Auth required:** Yes

**Path Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | User ID |

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
    "created_at": "2026-04-01T10:00:00Z",
    "updated_at": "2026-04-08T14:30:00Z"
  },
  "error": null
}
```

**Error Response (400):**

```json
{
  "data": null,
  "error": {
    "code": "BAD_REQUEST",
    "message": "invalid user id"
  }
}
```

**Error Response (404):**

```json
{
  "data": null,
  "error": {
    "code": "NOT_FOUND",
    "message": "user not found"
  }
}
```
