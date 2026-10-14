# Repository guidelines

## Stack

Go 1.26.5, Gin, the MongoDB v2 driver, and Redis for inventory caching.

## Architecture

Modular monolith with hexagonal modules (ports and adapters):

- `cmd/api/main.go`: reads configuration, connects MongoDB/Redis, calls
  `app.New`, and serves HTTP. Nothing else.
- `internal/app/`: composition root; registers every module and global
  middleware. Generated Swagger lives in `internal/app/swagger`.
- `internal/core/`: code shared by modules (config, exception, audit port,
  database, cache, httpx and its middleware). Core never imports a module.
- `internal/modules/<module>/`: `inventory` (employees, machines, monitors,
  phones, printers, emails), `organization` (branches, departments), and
  `audit` (change records and feed). Each has `domain/` (dto, exception),
  `application/` (port, usecase), `infrastructure/` (inbound/http/adapter,
  outbound/data/{adapter,schema}), and `module.go` with `Register`.

Dependencies point inward (infrastructure → application → domain). Use cases
depend only on ports. A module never imports another module: shared contracts
go to `core` as interfaces and the owning module registers the implementation
(for example `audit.Register` returns the `core/audit.Recorder`). File names use
snake_case; tests sit next to the tested file.

Preserve existing inventory routes, response envelopes, IDs, validation, and
BSON mappings. Keep database queries cancellable.

## Authentication

Reuse the sibling `../asset-control-auth` module via go.mod; do not duplicate
JWT, bcrypt, cookie, CORS, or guessing-protection implementations. It is a
library without a server: this API owns the HTTP listener, configuration, and
the shared MongoDB connection lifecycle, wired in `cmd/api/main.go` and
`internal/app/app.go`. The API reads `AUTH_*` variables and passes `auth.Config`.

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
- [docs/errors.md](docs/errors.md): the `{message, code, logs}` envelope and
  error codes. Report every handler failure through `httpx.WriteError`.

## Changes and commits

- Update CHANGELOG.md in every commit.
- Keep changes focused; use short Conventional Commit subjects with scope
  (`type(scope): imperative description`). Name branches `feat-<name>`,
  `fix-<name>`, `refactor-<name>`, or `docs-<name>`.
- Record non-obvious decisions in `agents/decisions/<branch>/`, failures in
  `agents/failures/<branch>/`, and findings in `agents/findings/<branch>/`.
- Do not add an AI co-author trailer.
- Ask the human for approval before committing or opening a PR.
