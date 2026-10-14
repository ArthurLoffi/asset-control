# Repository guidelines

## Stack

Go 1.26.5, Gin, the MongoDB v2 driver, and Redis for inventory caching.

## Architecture

Hexagonal layout under `src/main`; dependencies point inward
(infrastructure → application → domain):

- `domain/`: DTOs, validation, domain errors, and audit types. No HTTP or database
  dependencies.
- `application/ports/`: interfaces implemented by infrastructure.
- `application/usecase/`: business orchestration that depends only on ports.
- `infrastructure/inbound/http/`: controllers, middleware, routes, and Swagger.
- `infrastructure/outbound/`: Redis cache and MongoDB data (database, entities,
  repository).
- `main.go`: configuration and dependency wiring only.

Preserve existing inventory routes, response envelopes, IDs, validation, and
BSON mappings. Keep database queries cancellable.

## Authentication

Reuse the sibling `../asset-control-auth` module via go.mod; do not duplicate
JWT, bcrypt, cookie, CORS, or guessing-protection implementations. It is a
library without a server: this API owns the HTTP listener, configuration, and
the shared MongoDB connection lifecycle, wired in `src/main/main.go`.

- Auth routes register before the session middleware; every inventory route and
  Swagger register after it.
- Never add a production auth bypass or start without `AUTH_JWT_SECRET`.
- Apply request rate limiting after public auth routes and before protected
  inventory routes, with deterministic middleware coverage.
- Trust `CF-Connecting-IP` only when the immediate peer is local cloudflared
  (`127.0.0.1` or `::1`); keep direct forwarded headers untrusted. Rate limiting
  and login attempt protection use the validated Gin client IP.
- For HTTPS, configure an exact frontend origin and a Secure cookie.
- Use `gin.Logger()` for request logs; keep the shared request-ID middleware for
  response headers and login auditing without duplicating request logs.
- Keep `.env.example` and README.md aligned with auth configuration. Never
  commit real secrets.

## Tests

- `make test` / `go test ./...` and `make vet` / `go vet ./...`.
- `make test-integration` uses the `integration` tag and `MONGODB_TEST_URI`.
  It creates/drops only fresh `asset_control_test_*` databases and must never
  mutate a real inventory database.
- Keep the `testRouter` helper using real signed cookies instead of disabling
  authentication.
- Regenerate Swagger with `make swagger` after changing route annotations.

## Feature documentation

Feature contracts live in `docs/`. Update the matching file when behavior
changes:

- [docs/deactivation.md](docs/deactivation.md): multiple contacts/assignees,
  deactivation, inactive history, unassigned lists.
- [docs/audit-history.md](docs/audit-history.md): audit records and the recent
  changes feed.

## Changes and commits

- Update CHANGELOG.md in every commit.
- Keep changes focused; use short Conventional Commit subjects
  (`type: imperative description`).
- Do not add an AI co-author trailer.
- Ask the human for approval before committing or opening a PR.
