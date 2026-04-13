# Backend Tests

All backend test files live in this directory as a single `tests` package within the Go module. They access the `internal/*` packages directly because they're inside the same module.

## Files

| File | What it tests |
|------|---------------|
| `validator_test.go` | Table tennis score validation (set rules, deuce, best-of-N, edge cases) |
| `auth_service_test.go` | JWT generation, validation, expiry, signing method enforcement |
| `dto_test.go` | API response envelope helpers (Success, ErrorResponse, Meta) |
| `config_test.go` | Environment variable loading and required-field checks |
| `middleware_test.go` | Auth middleware (cookie, header, precedence) and CORS middleware |

## Coverage of Phase 1 Features

- **Table tennis rules**: 21-point sets, win-by-2, deuce scenarios, best-of-3/5/7
- **Match validation**: too few/many sets, extra sets after match decided, no winner, invalid scores
- **JWT auth**: token generation, validation, expiry, wrong secret, wrong signing method
- **Auth middleware**: cookie + header support, missing/invalid token, cookie precedence
- **CORS**: headers set, OPTIONS preflight returns 204, custom origin
- **Config**: all required env vars, defaults, custom overrides

## Running Tests

```bash
cd backend
go test ./tests/... -v
```

Or with coverage:

```bash
go test ./tests/... -cover
```
