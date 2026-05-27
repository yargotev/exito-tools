# Archive Report

**Change**: `2026-05-27-replace-orders-provider-with-geoms`  
**Archived on**: 2026-05-27  
**Archived to**: `openspec/changes/archive/2026-05-27-2026-05-27-replace-orders-provider-with-geoms/`

## Summary

Replaced the placeholder Orders provider behavior behind stable `orders.get` with GEOMS-backed lookup behavior, including Azure AD token acquisition, `findOrders`, `getOrder`, and `findItemsByOrder` enrichment. The public Capability ID and `exito orders get --id` command remain stable, with optional `--order-type` support.

## Verification

- `go test ./internal/domain/orders -run 'TestHTTPGetter(PostsRequestAndMapsProviderResponse|UsesProvidedGEOMSOrderType|FetchesGEOMSToken)'` ✅
- `go test ./internal/surface/cli -run TestOrdersGetCommandPassesOrderTypeFlag` ✅
- `go test ./internal/config -run 'TestResolveOrdersProviderConfiguration/prod_profile_reads_PDN_GEOMS_credential_bundle'` ✅
- `make test` ✅
- `go build ./cmd/exito` ✅
- `go test ./... -cover` ✅ 79.2% total coverage
- `go vet ./...` ✅
- `make lint` ✅ 0 issues
- `go test ./internal/domain/orders -coverprofile=/tmp/orders-cover.out` ✅ 80.0% package coverage

`make precommit` passed after staging the intentional archive/test/spec changes.

## Specs Synced

| Domain | Action | Details |
|--------|--------|---------|
| `capability-contract-foundation` | Updated | Replaced `orders.get domain execution` with GEOMS token/findOrders/details/items requirements and retained not-found behavior. |
| `cli-root` | Updated | Updated Orders get command requirement to document optional `--order-type` and Carulla scenario while retaining ID-required behavior. |
| `configuration-resolver` | Updated | Updated Orders provider configuration to include GEOMS client credentials, QA/PDN bundles, token fallback, and missing-credential behavior. |

## Archive Contents

- `proposal.md` ✅
- `design.md` ✅
- `tasks.md` ✅ 7/7 complete
- `apply-progress.md` ✅
- `verify-report.md` ✅ PASS
- `specs/` ✅

## Source of Truth Updated

- `openspec/specs/capability-contract-foundation/spec.md`
- `openspec/specs/cli-root/spec.md`
- `openspec/specs/configuration-resolver/spec.md`

## Notes

The verification report records no critical issues, warnings, or suggestions. Required scenarios are directly covered by tests, and GEOMS HTTP getter coverage was raised with focused helper/error-path tests.
