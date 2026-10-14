# Generic resource core instead of one use case file per entity

The blueprint asks for one use case struct per action with `Execute`. Six
inventory entities share the same six actions, so writing them per entity
would mean 36 near-identical files. `internal/core/resource` keeps the
blueprint shape (one struct per action, depending only on a port) but is
generic over `TRequest`, `TInput`, `TResponse`. Entities differ only in data:
filters, sort columns, mapping, and the not-found error, which
`database.ResourceConfig` carries. Employee deactivation needs extra work in
the same transaction, expressed as `AfterDeactivate` instead of a subclass.

Inputs validate through `Validated() (T, error)` (value receiver returning the
normalized copy) so the generic constraint needs a single type parameter.

Kept on purpose:

- Gin, because asset-control-auth and the existing middleware are Gin-based;
  the blueprint's `net/http` examples translate one to one.
- One shared MongoDB database. The blueprint gives each module its own database;
  the inventory data is already imported into one database, so modules only
  own their collections.
- snake_case resource fields (the blueprint suggests camelCase) to avoid
  breaking stored data and clients beyond the route change.

Swagger: swag cannot resolve a generic type across several packages named
`dto`, so list responses are documented with `dto.PageDoc{content=[]...}` and
adapters blank-import `core/dto` and `core/exception` (the swag-documented way
to reference unused packages).

Text sorting uses a `pt`, strength 2 collation on both count and find, so
"ana" sorts next to "Ana Maria" instead of after every capitalized name.
