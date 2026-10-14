# Audit history

Successful PUT updates for employees, machines, phones, monitors, printers, and
emails must write an atomic audit record to the `changes` collection in the same
MongoDB transaction. Audit identity comes only from the validated auth session,
never from request headers or bodies. Store before/after snapshots, resource,
resource ID, action, UTC time, and user ID. Redact passwords, password hashes,
tokens, and similar secrets from audit snapshots.

## Recent changes feed

`GET /api/changes/recent` is a protected dashboard feed limited to ten updates,
ordered by time and ID. Resolve current editor names from the `users` collection
in one query. Return changed scalar fields and their before/after values, omitting
secrets and lifecycle metadata; never return full audit snapshots.
