# Contacts, assignments, and deactivation

## Multiple contacts and assignees

Employee contacts accept `emails` (one or more distinct addresses) and
`phones` (distinct numbers). Machine, monitor, and phone assignments accept
`employee_ids` (distinct positive IDs). Preserve the singular `email`, `phone`,
and `employee_id` fields as the first value for legacy clients and imported
records. `?employee_id=` matches any associated employee; `?unassigned=true`
requires no associated employee. Employee deactivation removes only that
employee from each active equipment record.

## Deactivation

Employees, machines, phones, monitors, printers, and emails support deactivation.
`DELETE /api/<resource>/{id}` sets the database timestamp once and preserves the
record, returning `{id, deactivated_at}`; retries return the original timestamp
and missing IDs return 404 `<ENTITY>_001`. Missing/null `deactivated_at` means
active (including legacy imports). All ordinary lists, filters, lookups, and
updates use the active scope; ordinary writes never change lifecycle metadata
or reactivate a record. Deactivating an employee atomically (one transaction,
`ResourceConfig.AfterDeactivate`) unassigns active machines, monitors, and
phones; it does not deactivate the equipment.

`GET /api/<resource>/inactive` returns only `{id, deactivated_at}` with the
listing contract (`listing.md`), sorted by ID. Keep history requests separate
from active inventory and dashboard reads.
