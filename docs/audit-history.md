# Audit history

Successful `PUT /api/<resource>/{id}` updates for employees, machines, phones, monitors, printers, and
emails must write an atomic audit record to the `changes` collection in the same
MongoDB transaction. Audit identity comes only from the validated auth session,
never from request headers or bodies. Store before/after snapshots, resource,
resource ID, action, UTC time, and user ID. Redact passwords, password hashes,
tokens, and similar secrets from audit snapshots.

The shared `resource.UpdateUseCase` records the change through the
`core/audit.Recorder` port; the audit module implements it.

## Changes feed

`GET /api/changes` is a protected feed using the listing contract
(`listing.md`), ordered by time and record ID; the dashboard requests
`?direction=DESC&page_size=10`. Resolve current editor names from the `users` collection
in one query. Return changed scalar fields and their before/after values, omitting
secrets and lifecycle metadata; never return full audit snapshots.
