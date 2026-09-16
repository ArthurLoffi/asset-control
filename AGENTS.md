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
active scope. Deactivating an employee does not cascade to assigned equipment.
GET /api/<resource>/inactive?after=0&limit=50 returns only deactivation IDs and
timestamps, with next_after cursor pagination (limit 1-100). Keep history requests
separate from active inventory and dashboard reads. Never reset a timestamp on
retry or let ordinary writes change lifecycle metadata.
