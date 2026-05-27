# Archive Report

**Change**: `2026-05-27-add-vtex-intelligent-search-products`  
**Archived on**: 2026-05-27  
**Archived to**: `openspec/changes/archive/2026-05-27-add-vtex-intelligent-search-products/`

## Summary

Implemented Phase 1 VTEX Intelligent Search product search as a new read-only Catalog capability, preserving the existing Legacy Catalog search contract. The change adds `catalog.intelligent-search-products`, the explicit CLI command `exito catalog intelligent-search products`, non-sensitive per-brand Intelligent Search base URL configuration, and documentation.

## Verification

- `./exito catalog intelligent-search products --trade-policy 1 --text tv --count 3` ✅ returned `ok:true` with TV products.
- `go test ./internal/domain/catalog ./internal/config ./internal/app ./internal/surface/cli` ✅
- `go test ./...` ✅
- `make test` ✅
- `go build ./cmd/exito` ✅
- `./exito capabilities | grep -n "intelligent-search"` ✅
- `./exito catalog intelligent-search products --trade-policy 1 --by sku-id --value 912350 --count 1` ✅

## Specs Synced

| Spec | Action |
| --- | --- |
| `openspec/specs/catalog/spec.md` | Added VTEX Intelligent Search product capability requirements. |
| `openspec/specs/configuration-resolver/spec.md` | Added VTEX Intelligent Search provider configuration requirements. |
| `openspec/specs/cli-root/spec.md` | Added Intelligent Search products command requirements. |

## Archive Contents

- `proposal.md` ✅
- `design.md` ✅
- `tasks.md` ✅ 8/8 complete
- `apply-progress.md` ✅
- `verify-report.md` ✅ PASS
- `archive-report.md` ✅
- `specs/` ✅

## Issues

None.
