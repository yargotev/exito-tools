# Design: Add Master Data Read Tools

## Technical Approach

Add a new Master Data Operational Domain following existing `orders` VTEX credential and `checkout` brand-client patterns. The first slice exposes only read-only capabilities through the existing registry and generic `exito run` path; no write commands or mutation endpoints are added.

## Architecture Decisions

| Decision | Choice | Alternatives considered | Rationale |
|----------|--------|-------------------------|-----------|
| Domain ownership | Create `internal/domain/masterdata` | Fold into Orders/Catalog | Master Data has distinct document, schema, index, and pagination semantics. |
| Config model | Add `config.VTEXMasterDataProvider` with Exito/Carulla brand providers using app key/token | Reuse OMS provider directly | Base URL differs by API/account; credentials can reuse env names but resolved model stays provider-specific. |
| Capability exposure | Register six executable read-only capabilities; rely on `exito run` first | Add explicit Cobra commands now | Existing generic run already supports JSON inputs; this keeps the slice narrow. |
| Results | Return domain summaries plus `map[string]any` data/schema under controlled fields | Return raw HTTP DTOs | Domain-owned wrappers keep diagnostics safe and contracts stable. |
| Pagination | Use envelope pagination for scroll token; result diagnostics for search range/total | Auto-fetch all pages | Product rules forbid unbounded work and hidden pagination. |

## Data Flow

```text
config.Resolve ──→ app.New ──→ masterdata.NewBrandClient
                                  │
exito run ──→ execution.Pipeline ──→ masterdata.UseCase ──→ HTTPClient ──→ VTEX Master Data
                                  │                         │
                                  └──── domain result + warnings/pagination
```

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/config/config.go` | Modify | Parse `vtexMasterData`, env keys, secret-omitting provider structs. |
| `internal/config/config_test.go` | Modify | TDD coverage for profile/brand config, env override, serialization safety. |
| `internal/domain/masterdata/masterdata.go` | Create | Capability IDs, inputs/results, use cases, validation. |
| `internal/domain/masterdata/brand_client.go` | Create | Exito/Carulla routing and unavailable client. |
| `internal/domain/masterdata/http_client.go` | Create | VTEX document/search/scroll/schema/index GET client and DTO mapping. |
| `internal/domain/masterdata/*_test.go` | Create | Use case and fake-server HTTP tests. |
| `internal/app/app.go` / `app_test.go` | Modify | Wire clients and assert registry entries. |
| `docs/configuration.md` | Modify | Document YAML and env variables. |

## Interfaces / Contracts

```go
type Client interface {
    GetDocument(context.Context, GetDocumentInput) (DocumentResult, error)
    SearchDocuments(context.Context, SearchDocumentsInput) (DocumentsPage, error)
    ScrollDocuments(context.Context, ScrollDocumentsInput) (DocumentsPage, error)
    ListSchemas(context.Context, EntityInput) (SchemasResult, error)
    GetSchema(context.Context, GetSchemaInput) (SchemaResult, error)
    ListIndices(context.Context, EntityInput) (IndicesResult, error)
}
```

Inputs normalize `brand` to `exito` by default. Required strings: `entity`; plus `documentId` for get-document and `schema` for get-schema. Search enforces range size ≤100; scroll enforces size ≤1000.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|--------------|----------|
| Unit | Config resolution, input validation, unavailable client, capability metadata | Table-driven Go tests first. |
| HTTP | Paths, query params, REST-Range, `X-VTEX-MD-TOKEN`, app key/token headers, error mapping | `httptest.Server` fake VTEX. |
| App/Pipeline | Registry contains six IDs and execution returns JSON-safe data | Existing app/CLI run test patterns. |

## Migration / Rollout

No migration required. Users add optional `profiles.<profile>.vtexMasterData.<brand>.baseUrl` and existing or future VTEX app key/token env vars.

## Open Questions

- [ ] Confirm exact Master Data base URLs and whether existing VTEX app credentials have Master Data permissions.
