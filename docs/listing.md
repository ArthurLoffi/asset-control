# Resource routes and the listing contract

Every inventory entity (`employees`, `machines`, `monitors`, `phones`,
`printers`, `emails`) exposes the same routes; `branches` and `departments`
expose only the two reads. All require a session.

| Route | Response |
| --- | --- |
| `GET /api/<resources>` | 200 `{content, paging}` (contract below) |
| `GET /api/<resources>/{id}` | 200 record, or 404 `<ENTITY>_001` |
| `POST /api/<resources>` | 201 saved record with generated `id` and `created_at` |
| `PUT /api/<resources>/{id}` | 200 saved record; audited (see `audit-history.md`) |
| `DELETE /api/<resources>/{id}` | 200 `{id, deactivated_at}` soft delete (see `deactivation.md`) |
| `GET /api/<resources>/inactive` | 200 `{content: [{id, deactivated_at}], paging}` |
| `GET /api/changes` | 200 `{content, paging}` of audit summaries |

Records are returned directly, without a `{status, data}` wrapper. Errors use
`{message, code, logs}` (`errors.md`).

## Query parameters

| Param | Default | Rule |
| --- | --- | --- |
| `page` | `0` | Integer `>= 0`, zero-based. |
| `page_size` | `0` | Integer `0..100`; `0` means the default `20`. |
| `sort_by` | `id` | One of the entity's sortable fields; anything else is 400 `CORE_001`. The id is always the tie-breaker. Text sorts case-insensitively (`pt` collation). |
| `direction` | `ASC` | `ASC` or `DESC`, case-insensitive. |

Filters are optional; an absent parameter means no filter. Text filters match
partially, case-insensitively, and literally (regular expression characters
are escaped). Other filters are exact. Unknown parameters are ignored.

| Resource | Filters | `sort_by` |
| --- | --- | --- |
| employees | `name`, `branch_id`, `department_id` | `id`, `name`, `created_at` |
| machines | `employee_id` (repeat for any of several), `unassigned=true` | `id`, `user`, `model`, `created_at` |
| monitors, phones | `employee_id`, `unassigned=true` | `id`, `model`, `created_at` |
| printers | `branch_id`, `departments_id` | `id`, `model`, `created_at` |
| emails | `user` | `id`, `user`, `created_at` |
| branches, departments | `name` | `id`, `name` |
| changes | — | `time` (ordered by time, then record id) |

`employee_id` matches the legacy `employee_id` or any entry of `employee_ids`.
`unassigned=true` keeps records with neither. Lists never include deactivated
records.

## Response

```json
{
  "content": [{ "id": 42, "name": "Ana" }],
  "paging": { "total_elements": 42, "size": 20, "page": 0, "total_pages": 3,
              "has_content": true, "has_next": true, "has_previous": false }
}
```

- `content` is never `null`; `size` is the page size actually used.
- A page past the end is 200 with `content: []`, never 404.
- Count and find run in parallel (`core/database.FindPage`).

## Adding an entity

Declare the domain DTOs (`XxxRequestDTO` embedding `coredto.PaginatedRequestDTO`,
`XxxInputDTO` with `Validated()`, `XxxResponseDTO`, `XxxSortableFields`), the
schema (`XxxDocument`, `XxxCollection`, `XxxSortColumns`), the mapper, a data
adapter built with `database.NewResourceAdapter`, an API adapter, and one line
in `module.go`. The ports, use cases, handlers, and pagination are shared from
`internal/core`.
