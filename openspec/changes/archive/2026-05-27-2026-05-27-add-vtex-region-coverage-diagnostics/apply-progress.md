# Apply Progress

## Change

`2026-05-27-add-vtex-region-coverage-diagnostics`

## Status

Implemented Phase 2 VTEX Intelligent Search roadmap support as a read-only Geo Domain VTEX Checkout Regions coverage diagnostic capability and CLI command.

## Retrospective Note

This artifact was added after implementation because it was omitted during the manual apply work. It is therefore **retrospective evidence**, not a perfect chronological RED/GREEN log. The referenced tests and verification commands exist and pass now, but the original RED timing cannot be independently reconstructed from tool output.

## TDD Cycle Evidence

| Task | RED | GREEN | TRIANGULATE | SAFETY NET | REFACTOR |
| --- | --- | --- | --- | --- | --- |
| Add Geo OpenSpec delta for read-only VTEX region coverage diagnostics. | ➖ Spec artifact, no executable test. | ✅ Written in `openspec/changes/2026-05-27-add-vtex-region-coverage-diagnostics/specs/geo/spec.md`. | ➖ Artifact task. | ➖ No runtime files changed for this task. | ✅ Delta kept narrow to Geo coverage diagnostics. |
| Implement Geo domain input/result models, capability definition, validation, and unavailable resolver. | ⚠️ Retrospective: tests now exist in `internal/domain/geo/geo_test.go`. | ✅ `TestResolveVTEXRegionCapabilityExecutesUseCase` passes and proves capability input mapping/use-case execution. | ✅ Covered with capability execution and required-flag CLI validation. | ✅ Existing Geo tests still pass under `go test ./internal/domain/geo`. | ✅ Domain remains Cobra/Bubble Tea-free. |
| Implement VTEX Checkout Regions HTTP adapter and brand resolver. | ⚠️ Retrospective: fake-server tests now exist in `internal/domain/geo/geo_test.go`. | ✅ `TestHTTPVTEXRegionResolverBuildsCoordinatesAndCoverage` and `TestHTTPVTEXRegionResolverCoverageFalseForOnlyAccountSeller` pass. | ✅ Covers true coverage and account-only false coverage; no-sellers false is implied by empty seller extraction but not separately tested. | ✅ Catalog/Geo tests pass under `go test ./...`. | ✅ HTTP adapter keeps provider DTO mapping in Geo domain and uses shared HTTP infrastructure. |
| Wire the capability in `internal/app` from existing VTEX public base URLs. | ⚠️ Retrospective: app package covered by existing boot tests; no new focused app test added. | ✅ `go test ./internal/app` passes with the new registration path. | ➖ Wiring path is simple explicit registration. | ✅ Existing app tests pass. | ✅ Uses explicit application wiring; no side-effect registration introduced. |
| Add CLI command `geo resolve-vtex-region`. | ⚠️ Retrospective: CLI tests now exist in `internal/surface/cli/root_test.go`. | ✅ `TestGeoResolveVTEXRegionCommandRunsCapability` and `TestGeoResolveVTEXRegionCommandRequiresCoordinatesAndSalesChannel` pass. | ✅ Covers successful JSON envelope and missing required flag behavior. | ✅ Existing CLI tests pass under `go test ./internal/surface/cli`. | ✅ Cobra code uses `RunE`, `cobra.NoArgs`, required flags, and `cmd.OutOrStdout()`. |
| Add unit tests for URL construction, coverage rule, CLI wiring, and required flags. | ✅ Tests added in `internal/domain/geo/geo_test.go` and `internal/surface/cli/root_test.go`. | ✅ Focused tests and `go test ./...` pass. | ✅ URL construction, coverage true, coverage false, CLI success, and required flags covered. | ✅ Full suite passes. | ✅ Assertions check behavior, not implementation-only details. |
| Run `make test`. | ➖ Verification task. | ✅ `make test` passed. | ➖ Single verification command. | ✅ Full repository test suite passed. | ➖ No refactor for verification task. |

## Verification Evidence

| Check | Command | Result |
| --- | --- | --- |
| Focused Geo package | `go test ./internal/domain/geo` | ✅ Passed |
| Focused CLI package | `go test ./internal/surface/cli` | ✅ Passed |
| Full test suite | `go test ./...` | ✅ Passed |
| Make test | `make test` | ✅ Passed |
| Build | `go build ./cmd/exito` | ✅ Passed |
| Coverage | `go test ./... -cover` | ✅ Passed |
| Vet | `go vet ./...` | ✅ Passed |
| Formatting | `gofmt -l cmd internal` | ✅ No files listed |

## Files Changed

- `openspec/changes/2026-05-27-add-vtex-region-coverage-diagnostics/proposal.md` — change intent and non-goals.
- `openspec/changes/2026-05-27-add-vtex-region-coverage-diagnostics/design.md` — technical design for read-only VTEX Checkout Regions diagnostics.
- `openspec/changes/2026-05-27-add-vtex-region-coverage-diagnostics/specs/geo/spec.md` — Geo spec delta and scenarios.
- `openspec/changes/2026-05-27-add-vtex-region-coverage-diagnostics/tasks.md` — completed task checklist.
- `internal/domain/geo/vtex_region.go` — input/result models, capability definition, use case, validation.
- `internal/domain/geo/http_vtex_region_resolver.go` — VTEX Checkout Regions HTTP adapter, seller extraction, coverage rule.
- `internal/domain/geo/vtex_brand_region_resolver.go` — Exito/Carulla brand dispatch.
- `internal/domain/geo/unavailable.go` — unconfigured VTEX region resolver behavior.
- `internal/app/app.go` — explicit application wiring and capability registration.
- `internal/surface/cli/root.go` — Cobra command `geo resolve-vtex-region` and flags.
- `internal/domain/geo/geo_test.go` — Geo capability and HTTP adapter tests.
- `internal/surface/cli/root_test.go` — CLI command tests.

## Known Gaps

- This is retrospective evidence, so RED timing and pre-change safety-net execution cannot be independently proven.
- No dedicated test currently covers a successful provider response with an empty sellers array.
- Brand dispatch is simple and wired, but direct unit coverage for `VTEXBrandRegionResolver` is still low.
