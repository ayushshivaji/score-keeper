# Auth API

## GET /api/v1/auth/google

Initiates Google OAuth login flow.

**Auth required:** No

**Description:** Redirects the user to Google's OAuth consent screen. After the user approves, Google redirects back to the callback URL.

**Response:** `307 Temporary Redirect` to Google OAuth consent URL.

---

## GET /api/v1/auth/google/callback

Handles the OAuth callback from Google.

**Auth required:** No

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `code` | string | Yes | Authorization code from Google |

**Description:** Exchanges the authorization code for Google tokens, fetches the user's profile (email, name, avatar), and upserts the user in the database. Sets `access_token` (15 min, HTTP-only cookie) and `refresh_token` (7 days, HTTP-only cookie).

**Response:** `307 Temporary Redirect` to `http://localhost:3000/dashboard`

**Error Response (400):**

```json
{
  "data": null,
  "error": {
    "code": "BAD_REQUEST",
    "message": "missing code"
  }
}
```

---

## POST /api/v1/auth/refresh

Refreshes the access token using the refresh token cookie.

**Auth required:** No (uses refresh_token cookie)

**Description:** Validates the refresh token, rotates it (deletes old, creates new), and issues a new access token. Both tokens are set as HTTP-only cookies.

**Request:** No body required. The `refresh_token` cookie is read automatically.

**Success Response (200):**

```json
{
  "data": {
    "message": "tokens refreshed"
  },
  "error": null
}
```

**Error Response (401):**

```json
{
  "data": null,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "refresh token expired"
  }
}
```

---

## POST /api/v1/auth/logout

Logs out the current user.

**Auth required:** Yes

**Description:** Deletes all refresh tokens for the user and clears both `access_token` and `refresh_token` cookies.

**Request:** No body required.

**Success Response (200):**

```json
{
  "data": {
    "message": "logged out"
  },
  "error": null
}
```

---

## GET /api/v1/auth/me

Returns the currently authenticated user's profile.

**Auth required:** Yes

**Description:** Looks up the user by the `user_id` claim in the JWT access token.

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
