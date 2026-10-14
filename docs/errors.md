# Error envelope

Every error response, including auth, rate limiting, unknown routes, and
panics, has one shape:

```json
{ "message": "Request payload is invalid", "code": "CORE_001", "logs": { "validationErrors": ["name is required"] } }
```

- `message` is human-readable. `code` is stable; clients branch on it.
- `logs` carries debug context and is `{}` when there is none. It is public:
  never put secrets in it.
- Handlers report failures only through `httpx.WriteError`. Status ≥ 500 logs
  the original error with `slog.Error`; lower statuses log `slog.Warn`. Both
  include the request ID. Unknown errors stay opaque to clients.

| Code | Status | When |
| --- | --- | --- |
| `CORE_001` | 400 | Validation failed; `logs.validationErrors` lists messages. |
| `CORE_002` | varies | Generic HTTP error: unknown route (404), malformed JSON or invalid ID (400), missing record (404), rate limit (429). |
| `CORE_003` | 500 | Unhandled error or panic. |
| `CORE_004` | 409 | Duplicate record. |
| `CORE_005` | 504 | Database operation timed out. |
| `AUTH_001`–`AUTH_007` | 401/403/429 | Authentication; see `../asset-control-auth/README.md`. |

Domain errors use `exception.New(message, code, status, logs)` with a code
prefixed by the entity or module and a three-digit number (`MACHINE_001`).
