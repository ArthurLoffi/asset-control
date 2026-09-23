# Repository guidelines

Use Go 1.26.5, Gin, and the MongoDB v2 driver. Preserve the existing hexagonal
layout: adapter for HTTP/persistence, domain for DTOs/ports, usecase for business
orchestration. Preserve existing inventory routes, response envelopes, IDs,
validation, and BSON mappings. Keep database queries cancellable.

## Authentication

Reuse the sibling ../asset-control-auth module via go.mod; do not duplicate JWT,
bcrypt, cookie, CORS, or guessing-protection implementations. Wire dependencies
in src/main/main.go.
asset-control-auth is a library without an executable server. This API owns
the HTTP listener, configuration, and shared MongoDB connection lifecycle.
Auth routes register before session middleware; every inventory route and
Swagger register after it. Never add an auth bypass for
production or silently start without AUTH_JWT_SECRET. Keep .env.example and
README.md aligned with authentication configuration. Never commit real secrets.
Use gin.Logger() for HTTP request logs; retain the shared request-ID middleware
for response headers and login auditing without duplicating request logs.
Apply request rate limiting after public auth routes and before protected
inventory routes, with deterministic middleware coverage.

## Tests

Run go test ./... and go vet ./.... Integration tests use the integration tag
and MONGODB_TEST_URI, create/drop only fresh asset_control_test_* databases,
and must never mutate a real inventory database. Keep the testRouter helper
using actual signed cookies instead of disabling authentication.

## Changes

Use focused changes and short Conventional Commit subjects when committing.
Update AGENTS.md in every commit, consistent with the sibling projects.

## Deactivation

Employees, machines, phones, monitors, printers, and emails support deactivation.
POST /api/<resource>/deactivate?id=<id> sets the database timestamp once and
preserves the record. Missing/null deactivated_at means active (including legacy
imports). All ordinary lists, searches, assignment lookups, and updates use the
active scope. Deactivating an employee atomically unassigns active machines,
monitors, and phones; it does not deactivate the equipment.
`GET /api/machines/unassigned` and `GET /api/monitors/unassigned` list active
records with null or missing `employee_id`, sorted by ID.
GET /api/<resource>/inactive?after=0&limit=50 returns only deactivation IDs and
timestamps, with next_after cursor pagination (limit 1-100). Keep history requests
separate from active inventory and dashboard reads. Never reset a timestamp on
retry or let ordinary writes change lifecycle metadata.

## Audit history

Successful PUT updates for employees, machines, phones, monitors, printers, and
emails must write an atomic audit record to the `changes` collection in the same
MongoDB transaction. Audit identity comes only from the validated auth session,
never from request headers or bodies. Store before/after snapshots, resource,
resource ID, action, UTC time, and user ID. Redact passwords, password hashes,
tokens, and similar secrets from audit snapshots.
GET /api/changes/recent is a protected dashboard feed limited to ten updates,
ordered by time and ID. Resolve current editor names from the users collection
in one query. Return changed scalar fields and their before/after values, omitting
secrets and lifecycle metadata; never return full audit snapshots.
