# Asset control backend

Go/Gin API backed by MongoDB using the official Go driver v2. Existing `/api`
URLs, snake_case JSON fields, numeric IDs, and response envelopes are preserved.
POST responses now include the saved record's generated `id` and `created_at`.

## Run

Requires Go 1.26+ and MongoDB 4.4+ (integration tests use MongoDB 8.0).

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

Startup verifies connectivity, creates indexes, and initializes numeric ID
counters before accepting requests. SIGINT/SIGTERM drains HTTP requests and
closes the MongoDB connection pool. Request cancellations propagate to queries.
The connection behavior follows the [MongoDB Go driver context documentation](https://www.mongodb.com/docs/drivers/go/current/context/).

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
  -p 127.0.0.1:27028:27017 mongo:8.0
MONGODB_TEST_URI=mongodb://127.0.0.1:27028 make test-integration
docker stop asset-control-mongo-test
```

Each integration test creates and drops only a fresh `asset_control_test_*`
database. `MONGODB_TEST_URI` is required for the `integration` build tag;
explicit integration runs fail if it is missing or unreachable. The suite
covers all asset create/read/update/list routes, migrated documents, relationship
filters, literal case-insensitive name searches, defaults, null/false updates,
validation, error mapping, cancellation, concurrent IDs, and restart behavior.

Regenerate API documentation with `make swagger` (pinned generator version).
