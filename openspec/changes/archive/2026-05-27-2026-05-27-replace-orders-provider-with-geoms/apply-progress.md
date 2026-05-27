# Apply Progress

## Change

`2026-05-27-replace-orders-provider-with-geoms`

## Status

Implementation was already present when this verification pass began. This progress file records the implementation evidence available in the repository plus the verification-remediation tests added before archive.

## TDD Cycle Evidence

| Task | RED | GREEN | TRIANGULATE | SAFETY NET | REFACTOR |
|------|-----|-------|-------------|------------|----------|
| Replace Orders HTTP getter endpoint and DTOs with GEOMS `findOrders` | ✅ Written: `internal/domain/orders/http_getter_test.go` | ✅ Passed: `go test ./internal/domain/orders` | ✅ Default and explicit order-type request cases | ✅ Existing Orders domain tests pass | ✅ Domain-owned DTOs remain in `internal/domain/orders` |
| Enrich `orders.get` with GEOMS `getOrder` details and `findItemsByOrder` food/non-food items | ✅ Written: `internal/domain/orders/http_getter_test.go` | ✅ Passed: `go test ./internal/domain/orders` | ✅ Summary, details, food, and non-food endpoint assertions | ✅ Existing Orders domain tests pass | ✅ Enrichment remains behind stable `orders.GetResult` |
| Add GEOMS token acquisition and expiry-aware in-memory caching | ✅ Written: `internal/domain/orders/http_getter_test.go` | ✅ Passed: `go test ./internal/domain/orders` | ✅ Dynamic token and fallback token paths | ✅ Existing Orders domain tests pass | ✅ Token source isolated from surfaces |
| Extend Orders configuration for client credentials and GEOMS credential bundles | ✅ Written: `internal/config/config_test.go` | ✅ Passed: `go test ./internal/config` | ✅ QA bundle and prod PDN bundle cases | ✅ Existing config resolver tests pass | ✅ Secrets remain JSON-omitted |
| Preserve `orders.get` and `exito orders get --id` while adding optional `orderType` / `--order-type` | ✅ Written: `internal/domain/orders/orders_test.go`, `internal/surface/cli/root_test.go` | ✅ Passed: `go test ./internal/domain/orders ./internal/surface/cli` | ✅ Capability input forwarding and CLI flag forwarding | ✅ Existing CLI tests pass | ✅ CLI still executes through shared pipeline |
| Update non-sensitive docs and YAML base URLs | ✅ Written: spec/docs assertions covered by review | ✅ Passed: `go test ./...` | ➖ Documentation/config artifact update | ✅ Existing config tests pass | ✅ Secrets not added to committed YAML/docs |
| Run `make test` before handoff | ✅ Written: verification command evidence | ✅ Passed: `make test` during final verification | ➖ Single verification command | ✅ Full suite pass | ➖ Not applicable |

## Files Changed

| File | Change |
|------|--------|
| `internal/domain/orders/http_getter.go` | GEOMS HTTP getter, token source, DTO mapping, details/items enrichment. |
| `internal/domain/orders/http_getter_test.go` | GEOMS request/default order type, explicit Carulla order type, token caching, and error tests. |
| `internal/domain/orders/http_getter_internal_test.go` | GEOMS mapper/helper, empty details/items, downstream error, structured-code, and token error-path tests. |
| `internal/domain/orders/orders.go` | Stable capability definition with optional `orderType`. |
| `internal/domain/orders/orders_test.go` | Capability definition and input forwarding tests. |
| `internal/config/config.go` | Orders client credentials, token URL, and GEOMS QA/PDN bundle resolution. |
| `internal/config/config_test.go` | Orders provider, QA bundle, PDN bundle, dotenv precedence, and JSON omission tests. |
| `internal/surface/cli/root.go` | `orders get --order-type` flag and shared pipeline input. |
| `internal/surface/cli/root_test.go` | CLI `orders get` behavior and `--order-type` forwarding tests. |
| `docs/capabilities/orders.get.md` | GEOMS capability contract documentation. |
| `docs/configuration.md` | Non-sensitive GEOMS configuration documentation. |
| `exito.yaml` | Non-sensitive GEOMS base URL profile values. |
| `openspec/changes/2026-05-27-replace-orders-provider-with-geoms/*` | Proposal, design, specs, tasks, verification, and apply progress artifacts. |

## Verification Commands

- `go test ./internal/domain/orders -coverprofile=/tmp/orders-cover.out` ✅ 80.0% package coverage
- `go test ./internal/domain/orders -run 'TestHTTPGetter(PostsRequestAndMapsProviderResponse|UsesProvidedGEOMSOrderType|FetchesGEOMSToken)'`
- `go test ./internal/surface/cli -run TestOrdersGetCommandPassesOrderTypeFlag`
- `go test ./internal/config -run 'TestResolveOrdersProviderConfiguration/prod_profile_reads_PDN_GEOMS_credential_bundle'`
