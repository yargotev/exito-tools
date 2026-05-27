# Apply Progress

## Change

`2026-05-27-add-vtex-intelligent-search-products`

## Status

Implemented Phase 1 VTEX Intelligent Search product search as a new read-only Catalog capability and CLI command.

## TDD / Verification Evidence

| Task | Evidence |
| --- | --- |
| Configuration resolution | `internal/config/config_test.go` covers YAML and prod env overrides for `vtexIntelligentSearch` base URLs. |
| Domain capability and validation | `internal/domain/catalog/intelligent_search_test.go` covers typed multi-ID lookup, ambiguous query rejection, request construction, cookie redaction, and response mapping. |
| HTTP client | Fake-server test asserts `product_search/trade-policy/...` path facets, query params, cookie forwarding, and redacted output. |
| Application wiring | `go test ./internal/app` passes with the new capability registered during boot. |
| CLI command | `internal/surface/cli/root_test.go` covers `catalog intelligent-search products`, repeated values/facets, metadata, and required `--trade-policy`. |
| Full suite | `make test` passed. |

## Files Changed

- `internal/domain/catalog/intelligent_search.go` — input/result models, capability definition, validation, typed query construction.
- `internal/domain/catalog/http_intelligent_searcher.go` — VTEX Intelligent Search HTTP client and response mapper.
- `internal/domain/catalog/intelligent_brand_searcher.go` — brand dispatch for Exito/Carulla.
- `internal/domain/catalog/intelligent_unavailable.go` — stable not-configured behavior.
- `internal/config/config.go` — `vtexIntelligentSearch` provider resolution.
- `internal/app/app.go` — explicit boot registration.
- `internal/surface/cli/root.go` — nested Cobra command and flags.
- `docs/configuration.md`, `docs/capabilities/catalog.intelligent-search-products.md`, `exito.yaml` — non-sensitive docs/config.
