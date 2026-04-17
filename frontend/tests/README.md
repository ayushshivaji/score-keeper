# Frontend Tests

All frontend test files live in this directory. Tests use Jest with `ts-jest` and `jsdom`.

## Files

| File | What it tests |
|------|---------------|
| `validate-set.test.ts` | Table tennis set validation logic (mirrors backend rules) |
| `match-form-logic.test.ts` | Match form business logic: visible sets, match result, submission guards |
| `utils.test.ts` | Date and score formatting utilities |
| `api.test.ts` | API client wrapper: GET/POST/DELETE, credentials, error handling |

## Coverage of Phase 1 Features

- **Set score validation**: standard wins, deuce, ties, win-by-less-than-2, invalid combinations
- **Match form logic**: dynamic set visibility based on match state, winner determination, set count for best-of-3/5/7
- **Submission guards**: same-player check, missing players, incomplete match
- **Formatters**: date formatting, datetime formatting, set score string formatting
- **API client**: URL construction, JSON body, credentials cookie, error envelope parsing

## Running Tests

```bash
cd frontend
npm test
```

Or directly:

```bash
npx jest
```

Watch mode:

```bash
npx jest --watch
```
