# API Overview

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Protected endpoints require a valid JWT access token, sent automatically via HTTP-only cookies set during login. Alternatively, pass the token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

## Response Envelope

All responses follow this format:

### Success

```json
{
  "data": { ... },
  "error": null,
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 42
  }
}
```

`meta` is only present on paginated list endpoints.

### Error

```json
{
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human-readable error description"
  }
}
```

### Common Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing or invalid access token |
| `NOT_FOUND` | 404 | Resource does not exist |
| `BAD_REQUEST` | 400 | Malformed request or invalid parameter |
| `VALIDATION_ERROR` | 400 | Request body fails validation rules |
| `SERVER_ERROR` | 500 | Unexpected server error |

## Endpoint Groups

| Group | Doc |
|-------|-----|
| Auth | [api-auth.md](./api-auth.md) |
| Users | [api-users.md](./api-users.md) |
| Matches | [api-matches.md](./api-matches.md) |
| Leaderboard | [api-leaderboard.md](./api-leaderboard.md) |

## Pagination

List endpoints accept these query parameters:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | int | 1 | Page number (1-indexed) |
| `per_page` | int | 20 | Items per page (max 50) |
