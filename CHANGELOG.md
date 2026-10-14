# Changelog

All notable changes to this project are documented in this file.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added

- Paginated listing contract on every `GET /api/<resources>`: `page`,
  `page_size`, `sort_by`, `direction`, entity filters, and `{content, paging}`
  responses (`docs/listing.md`). Text sorts case-insensitively.
- Filters: employees `name`/`branch_id`/`department_id`; machines, monitors, and
  phones `employee_id` (repeatable) and `unassigned`; printers
  `branch_id`/`departments_id`; emails `user`; branches/departments `name`.
- Entity not-found codes (`EMPLOYEE_001`, `MACHINE_001`, … `DEPARTMENT_001`).

### Changed

- REST routes: `GET|PUT|DELETE /api/<resources>/{id}` replace `?id=` queries;
  `DELETE` replaces `POST /deactivate`; `GET /api/branches[/{id}]` and
  `/api/departments[/{id}]` replace `/api/utils/*`; `GET /api/changes` replaces
  `/api/changes/recent`. Records are returned without `{status, data}`.
- `/inactive` uses the listing contract instead of `after`/`next_after`.
- Generic core (`resource`, `database.ResourceAdapter`, `httpx` handlers)
  replaces the per-entity controllers, use cases, ports, and repositories.
- Each module creates its own indexes (`app.EnsureSchema`).
- Every error response uses `{message, code, logs}` (`CORE_001`–`CORE_005`,
  `AUTH_xxx`) instead of `{error}`; see `docs/errors.md`. Validation messages
  cite public field names and are listed in `logs.validationErrors`.
- JSON bodies are decoded once and validated once, after normalization, in the
  use case.
- Unknown routes return 404 `CORE_002` instead of Gin's plain-text page.
- Adopt the modular hexagonal layout: `cmd/api`, `internal/app`, `internal/core`,
  and `internal/modules/{inventory,organization,audit}`, each module with its own
  `domain`, `application`, `infrastructure`, and `module.go`.
- Read `AUTH_*` configuration in the API and use the single `auth` package of
  `asset-control-auth`; panic recovery moved to `internal/core/httpx`.
- Load `.env` once through `internal/core/config` instead of per component.
- Logout responds 204 (from `asset-control-auth`).
- Reorganize packages into `domain`, `application` (ports, use cases), and
  `infrastructure` (HTTP, cache, MongoDB) layers.
- Move feature contracts from AGENTS.md to `docs/`.

### Fixed

- Scope the rate limiter to each router instead of package globals, so separate
  routers (and integration tests) no longer share request budgets.
