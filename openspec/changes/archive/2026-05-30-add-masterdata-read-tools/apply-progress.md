# Apply Progress: Add Master Data Read Tools

## Mode

Strict TDD (`openspec/config.yaml` has `testing.strict_tdd: true`; runner: `go test ./...`).

## Completed Tasks

- [x] 1.1 Added `internal/config` tests for `vtexMasterData` YAML, env override, credential presence, and JSON secret omission.
- [x] 1.2 Added `VTEXMasterDataProvider`, env keys, YAML parsing, and profile/brand credential resolution in `internal/config/config.go`.
- [x] 1.3 Updated `docs/configuration.md` with `vtexMasterData` YAML and env variables.
- [x] 2.1 Added `internal/domain/masterdata` tests for definitions, inputs, defaults, limits, warnings, and pagination.
- [x] 2.2 Created `masterdata.go` with IDs, models, use cases, schemas, conversion, and validation.
- [x] 2.3 Created `brand_client.go` with brand routing and `MASTERDATA_NOT_CONFIGURED` unavailable behavior.
- [x] 3.1 Added fake-server HTTP tests for document/search/scroll/schema/index reads.
- [x] 3.2 Created `http_client.go` with VTEX app key/token headers, safe diagnostics, and DTO mapping.
- [x] 3.3 Implemented `REST-Range`, resources total parsing, missing-sort warning, scroll token metadata, and VTEX limit validation.
- [x] 4.1 Extended `internal/app/app_test.go` to require six `masterdata.*` registry entries.
- [x] 4.2 Wired Master Data brand clients and capability registration in `internal/app/app.go`.
- [x] 4.3 Added an app/pipeline test proving `masterdata.get-document` emits envelope-safe data through the registered capability.
- [x] 5.1 Ran focused tests for `internal/config`, `internal/domain/masterdata`, and `internal/app`.
- [x] 5.2 Ran `go test ./...` during implementation; final `make test` evidence below.
- [x] 5.3 Confirmed implementation uses only GET endpoints and introduced no write/delete/schema lifecycle/index mutation methods.

## TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 1.1 | `internal/config/config_test.go` | Unit | ✅ `go test ./internal/config ./internal/app` | ✅ Tests referenced missing `VTEXMasterDataProvider` | ✅ `go test ./internal/config` | ✅ YAML/env/missing credentials/JSON cases | ✅ Reused OMS provider shape |
| 1.2 | `internal/config/config_test.go` | Unit | ✅ Same as 1.1 | ✅ Covered by 1.1 RED | ✅ `go test ./internal/config` | ✅ Profile/brand/env variants | ✅ Shared existing resolver helpers |
| 1.3 | `docs/configuration.md` | Docs | N/A (docs) | ➖ Documentation task | ✅ Reviewed updated doc section | ➖ Single docs behavior | ✅ Narrow insertion |
| 2.1 | `internal/domain/masterdata/masterdata_test.go` | Unit | N/A (new package) | ✅ Tests failed with no non-test package | ✅ `go test ./internal/domain/masterdata` | ✅ Normalize/default/error/warning/pagination cases | ✅ Shared helper functions |
| 2.2 | `internal/domain/masterdata/masterdata_test.go` | Unit | N/A (new package) | ✅ Covered by 2.1 RED | ✅ `go test ./internal/domain/masterdata` | ✅ Multiple capability and validation paths | ✅ Extracted definitions/converters |
| 2.3 | `internal/domain/masterdata/masterdata_test.go` | Unit | N/A (new package) | ✅ Unavailable behavior covered by tests/spec | ✅ `go test ./internal/domain/masterdata` | ✅ Exito/default and Carulla routing via client shape | ✅ Small router file |
| 3.1 | `internal/domain/masterdata/http_client_test.go` | Unit with fake HTTP | N/A (new package) | ✅ HTTP tests referenced missing `NewHTTPClient` | ✅ `go test ./internal/domain/masterdata` | ✅ document/search/scroll/schema/index paths | ✅ Shared request helpers |
| 3.2 | `internal/domain/masterdata/http_client_test.go` | Unit with fake HTTP | N/A (new package) | ✅ Covered by 3.1 RED | ✅ `go test ./internal/domain/masterdata` | ✅ Auth headers and diagnostics assertions | ✅ Centralized `newGET`/decode |
| 3.3 | `internal/domain/masterdata/http_client_test.go` | Unit with fake HTTP | N/A (new package) | ✅ Tests asserted REST-Range/resources/token behavior | ✅ `go test ./internal/domain/masterdata` | ✅ Search and scroll separate cases | ✅ Header parsing helpers |
| 4.1 | `internal/app/app_test.go` | Unit/integration with registry | ✅ `go test ./internal/app` before edit | ✅ Import/constants missing before domain/app code | ✅ `go test ./internal/app` | ✅ All six capability IDs | ✅ Existing registry table extended |
| 4.2 | `internal/app/app_test.go` | Unit/integration | ✅ Same as 4.1 | ✅ Covered by 4.1 RED | ✅ `go test ./internal/app` | ✅ Configured and unconfigured wiring paths | ✅ Small wiring helpers |
| 4.3 | `internal/app/app_test.go` | Unit/integration with fake HTTP | ✅ `go test ./internal/app` | ✅ Added capability execution assertion | ✅ `go test ./internal/app` | ✅ Metadata, auth, and result mapping | ✅ Reused pipeline pattern |
| 5.1 | package tests | Unit | N/A | ➖ Verification task | ✅ `go test ./internal/config ./internal/domain/masterdata ./internal/app` | ➖ Command evidence | ➖ None needed |
| 5.2 | full suite | Unit | N/A | ➖ Verification task | ✅ `go test ./...`; `make test` passed | ➖ Command evidence | ➖ None needed |
| 5.3 | source scan/review | Static review | N/A | ➖ Verification task | ✅ Only GET methods/endpoints added | ➖ Static confirmation | ➖ None needed |

## Test Summary

- Focused tests: `go test ./internal/config ./internal/domain/masterdata ./internal/app` passed.
- Full suite during implementation: `go test ./...` passed.
- Final verification: `make test` passed.
- Quality: `make fmt` and `make lint` passed (`0 issues`).

## Deviations from Design

- None materially. Explicit Cobra commands remain deferred; generic registered capability execution is covered through the pipeline.

## Issues Found

- Open question remains: exact Master Data account URLs and credential permissions must be confirmed outside code.
