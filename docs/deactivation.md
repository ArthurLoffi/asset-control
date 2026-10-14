# Contacts, assignments, and deactivation

## Multiple contacts and assignees

Employee contacts accept `emails` (one or more distinct addresses) and
`phones` (distinct numbers). Machine, monitor, and phone assignments accept
`employee_ids` (distinct positive IDs). Preserve the singular `email`, `phone`,
and `employee_id` fields as the first value for legacy clients and imported
records. Assignment searches match any associated employee; unassigned lists
require no associated employee. Employee deactivation removes only that
employee from each active equipment record.

## Deactivation

Employees, machines, phones, monitors, printers, and emails support deactivation.
`POST /api/<resource>/deactivate?id=<id>` sets the database timestamp once and
preserves the record. Missing/null `deactivated_at` means active (including
legacy imports). All ordinary lists, searches, assignment lookups, and updates
use the active scope. Deactivating an employee atomically unassigns active
machines, monitors, and phones; it does not deactivate the equipment.

`GET /api/machines/unassigned` and `GET /api/monitors/unassigned` list active
records with null or missing `employee_id`, sorted by ID.

`GET /api/<resource>/inactive?after=0&limit=50` returns only deactivation IDs and
timestamps, with `next_after` cursor pagination (limit 1-100). Keep history
requests separate from active inventory and dashboard reads. Never reset a
timestamp on retry or let ordinary writes change lifecycle metadata.
