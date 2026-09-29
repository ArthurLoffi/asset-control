# Asset control backend

Go/Gin API backed by MongoDB using the official Go driver v2. Existing `/api`
URLs, snake_case JSON fields, numeric IDs, and response envelopes are preserved.
POST responses now include the saved record's generated `id` and `created_at`.

## Run

Requires Go 1.26+ and MongoDB 4.4+ on a replica set or sharded cluster
for transactional updates and employee deactivation (integration tests use MongoDB 8.0).

```sh
cp .env.example .env
# Edit .env for your deployment.
go run ./src/main
```

The server listens on port 8080; Swagger UI is at `/swagger/index.html`.
Run commands from the `asset-control` directory. Startup automatically reads
`.env` from the current working directory; exported environment variables take
precedence. The file is optional when configuration is supplied through the
environment. Do not commit credentials.

| Variable | Meaning |
| --- | --- |
| `MONGODB_URI` | Required MongoDB connection URI, including authentication/TLS options when applicable. |
| `MONGODB_DATABASE` | Required database containing the migrated collections. |
| `MONGODB_TIMEOUT` | Positive Go duration for MongoDB operations and startup checks; defaults to `10s`. |
| `REDIS_URL` | Optional Redis URL for caching protected inventory GET responses. If set, Redis must be reachable at startup. |
| `REDIS_CACHE_TTL` | Positive Go duration for cached responses; defaults to `1m`. |

Startup verifies connectivity, creates indexes, and initializes numeric ID
counters before accepting requests. SIGINT/SIGTERM drains HTTP requests and
closes the MongoDB connection pool. Request cancellations propagate to queries.
The connection behavior follows the [MongoDB Go driver context documentation](https://www.mongodb.com/docs/drivers/go/current/context/).

When `REDIS_URL` is configured, successful employee, machine, monitor, phone,
printer, and organization GET responses are cached in Redis. Each successful
create, update, or deactivation advances that resource's cache version, so later
reads use fresh data. Failed writes do not invalidate the cache. Email responses
are excluded because they contain password fields. Redis errors during requests
fall back to MongoDB; response entries expire after `REDIS_CACHE_TTL`. If Redis
is unavailable during a successful write, an older entry can remain until that
expiry, so keep the TTL short for deployments that need a tighter bound.

## Migrated data contract

This backend expects the original numeric `id`, `employee_id`, `branch_id`, and
`department_id`/`departments_id` fields, with the original snake_case field names.
MongoDB's `_id` is independent and is neither exposed nor used as the public ID.
IDs must be positive integers within BSON's signed 64-bit range. Dates should
be BSON dates, and nullable values should retain their BSON null/boolean types.

Collections remain `employees`, `machines`, `monitors`, `phones`, `printers`,
`branches`, and `departments`. The renamed email collection is **`usableEmails`**,
recorded as `database.UsableEmailsCollection`. Email lookups are available through
`GET /api/emails/id?id=42` and `GET /api/emails/user?user=Ana`.
`GET /api/emails` lists every email in ascending ID order using the response
`{ "status": "emails retrieved successfully", "data": [...] }`; an empty
collection returns `data: []`. Listed records include all five fields, including
the unchanged plaintext `password_hash`.
The ID must be a positive int32 (1 through 2147483647). User searches trim outer
whitespace and use literal, case-insensitive partial matching, returning the
lowest-ID match, consistent with employee name searches. Missing matches return
404; missing/invalid search values return 400.
Both return `{ "status": "email found successfully", "data": { ... } }`, with
`id`, `user`, `email`, `password_hash`, and `created_at` (RFC 3339 timestamp).
`password_hash` contains the plaintext password exactly as stored; lookups do
not hash, mask, or modify it. There is no email delete route.

Create emails with `POST /api/emails` and edit them with
`PUT /api/emails?id=42`. Both require a JSON body containing `user`, `email`
(a valid email address), and `password_hash` (a nonempty plaintext password).
PUT replaces all three fields; it preserves `id`, `created_at`, and unknown
imported fields, returning 404 for missing IDs. Passwords are not trimmed or
hashed. POST returns 201 with status `created`; PUT returns 200 with status
`email updated successfully`. Both include the saved record in `data`.
Startup creates a unique `id` index and seeds the `usableEmails` ID counter
from imported records. New IDs are stored as BSON int32; creation fails with
500 when the int32 ID range is exhausted. Restart the backend before using
the write routes so the new counter and index are initialized.
This application does not copy data or rename existing collections.

The additional `_counters` collection allocates numeric IDs atomically. On
startup, each counter is advanced to at least the largest imported ID. Counters
never move backwards on restart; failed writes can leave gaps. Unique `id`
indexes protect imported IDs and concurrent inserts. `employee_id` indexes
support asset lookups. The database user therefore needs read/write access and
permission to create indexes. Duplicate imported IDs cause startup to fail.

Finish imports before starting the backend. If further records are imported
later, stop application writes and restart the backend to reseed counters. Other
writers must use the same counter protocol. Do not delete/reset `_counters`.

PUT replaces the API's mutable fields: omitted nullable fields become null;
`false` remains false. It preserves `_id`, `id`, `created_at`, and fields unknown
to this backend. Missing IDs return 404 without inserting a document.
Machine, monitor, and phone situations default to `Em uso`; the original SQL
enums and required fields are validated in the application. No additional
foreign-key or email-uniqueness rules are assumed without the original schema.
For phone POST and PUT, `model`, `imei`, and `serial_number` may be empty when
`equipament` is `Só chip` or `Sem chip`. They remain required for `Celular` or
unspecified equipment, and retain the 100-character maximum in every case.
Existing identification values are preserved when supplied; no placeholders
are generated. Restart the API after updating it to apply this validation rule.
Database failures return 500, missing records 404, conflicts 409, and operation
timeouts 504. Internal database details are not exposed in responses.

## Validation tests

```sh
make test
make vet
```

Run integration tests against a disposable local MongoDB:

```sh
docker run --detach --rm --name asset-control-mongo-test \
  -p 127.0.0.1:27028:27017 mongo:8.0 --replSet rs0 --bind_ip_all
docker exec asset-control-mongo-test mongosh --quiet --eval 'rs.initiate()'
MONGODB_TEST_URI='mongodb://127.0.0.1:27028/?directConnection=true' make test-integration
docker stop asset-control-mongo-test
```

Each integration test creates and drops only a fresh `asset_control_test_*`
database. `MONGODB_TEST_URI` is required for the `integration` build tag;
explicit integration runs fail if it is missing or unreachable. The suite
covers all asset create/read/update/list routes, migrated documents, relationship
filters, literal case-insensitive name searches, defaults, null/false updates,
validation, error mapping, cancellation, concurrent IDs, and restart behavior.

Regenerate API documentation with `make swagger` (pinned generator version).

## Authentication

Requires Go 1.26.5 and sibling `../asset-control-auth` at build time (local Go
module replacement). The API mounts the auth module directly; no separate auth
server is needed. Configure `AUTH_JWT_SECRET` (random, at least 32 bytes),
`AUTH_TOKEN_TTL` (default `8h`), `AUTH_COOKIE_SECURE` (true for production HTTPS),
and `AUTH_ALLOWED_ORIGINS` (exact frontend origins) using `.env.example`.
Startup fails without a valid signing key.

Users live in the configured database's `users` collection, with MongoDB
ObjectId `_id`, unique case-sensitive `name`, and bcrypt `password`. Login audits
update only audit fields; credentials and inventory data are preserved. Startup
creates a unique `users.name` index; existing duplicates must be resolved first.

- `POST /api/auth/login`: JSON `{user, password}`; sets an HttpOnly JWT cookie
  and returns `{status, data: {id, name, expires_at}}`.
- `GET /api/auth/session`: returns the current session or 401.
- `POST /api/auth/logout`: clears the cookie, including after expiry.
- All inventory routes and Swagger require the session cookie. Missing/invalid
  tokens return 401; expired tokens include `code: "session_expired"`.
- Every write requires `X-CSRF-Protection: 1`. Requests from browsers must use
  an origin in `AUTH_ALLOWED_ORIGINS`; CORS credentials are enabled only there.

Use the frontend at the same origin/site as the API. Login and inventory requests
share the `/api` base URL. Read `../asset-control-auth/README.md` for the full
contract, deployment layout, user provisioning, and process-local guessing-guard
limitations. Logout clears the browser cookie; it does not revoke copied JWTs.

For an HTTPS frontend behind Cloudflare Tunnel, set `AUTH_COOKIE_SECURE=true`
and add the frontend's exact HTTPS origin to `AUTH_ALLOWED_ORIGINS`. Browser
login sends an `Origin` header even when `/api` is same-origin; an origin absent
from that list receives 403 before a session cookie is created. Restart the API
after changing these settings. Gin trusts `CF-Connecting-IP` only when the
immediate peer is local cloudflared (`127.0.0.1` or `::1`). Keep the API listener
private to the tunnel; forwarded headers from remote direct clients are ignored.

Integration inventory tests now supply signed test cookies and CSRF headers.

## Inventory deactivation

Employee create and update requests accept `emails` and `phones` arrays. At
least one distinct, valid email is required; phone numbers are optional. Machine,
monitor, and phone requests accept distinct positive IDs in `employee_ids`.
Responses retain `email`, `phone`, and `employee_id` as the first value and
include the full arrays. Existing singular-only requests and imported records
remain supported. Employee-ID lookups match any assignment.

For employees, machines, phones, monitors, printers, and emails:

- `POST /api/<resource>/deactivate?id=<id>` requires a session and browser CSRF
  protection. It preserves the document and sets `deactivated_at` using the MongoDB
  server clock. The `{status, data}` response contains `id` and the UTC timestamp.
  Retries return the original timestamp; missing IDs return 404.
- Ordinary lists, searches, employee assignment lookups, and PUT operations only
  access active records. Legacy documents with no timestamp (or null) remain active.
  Deactivating an employee atomically clears their assignment from active machines,
  monitors, and phones. The equipment remains active; inactive equipment is untouched.
- `GET /api/<resource>/inactive?after=0&limit=50` returns a separate history page:
  `{status, data: {records: [{id, deactivated_at}], next_after}}`. Pass `next_after`
  as the next `after`; null means the last page. Limit is 1-100. The endpoint returns
  lifecycle metadata only; the original document remains stored unchanged.
- Startup creates indexes for active lists, employee assignments, and inactive
  history. Dashboard requests never load inactive history. No migration of legacy
  records is required, and ordinary writes cannot reactivate records.

`GET /api/machines/unassigned` lists only active machines with no employee in
`employee_id` or `employee_ids`, in ascending ID order. An empty result is `data: []`.
`GET /api/monitors/unassigned` applies the same filter to monitors.
