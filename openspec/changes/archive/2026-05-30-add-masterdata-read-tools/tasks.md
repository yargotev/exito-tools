# Tasks: Add Master Data Read Tools

## Phase 1: Configuration Foundation

- [x] 1.1 RED: Add `internal/config` tests for `profiles.<profile>.vtexMasterData.<brand>.baseUrl`, env override, credential presence, and JSON secret omission.
- [x] 1.2 GREEN: Add `VTEXMasterDataProvider`, env keys, YAML parsing, and profile/brand credential resolution in `internal/config/config.go`.
- [x] 1.3 Update `docs/configuration.md` with `vtexMasterData` YAML and environment variable names.

## Phase 2: Domain Contracts and Validation

- [x] 2.1 RED: Add `internal/domain/masterdata` tests for capability definitions, required inputs, default brand, and VTEX limit validation.
- [x] 2.2 GREEN: Create `masterdata.go` with capability IDs, inputs/results, use cases, error codes, schemas, and validation.
- [x] 2.3 Create `brand_client.go` and unavailable client returning `MASTERDATA_NOT_CONFIGURED`.

## Phase 3: HTTP Provider Client

- [x] 3.1 RED: Add fake-server tests for get-document, search, scroll, list-schemas, get-schema, and list-indices paths/query/headers.
- [x] 3.2 GREEN: Create `http_client.go` using shared `httpclient`, VTEX app key/token headers, safe diagnostics, and DTO-to-domain mapping.
- [x] 3.3 Implement search `REST-Range` handling, resources parsing, no-sort warning, scroll token pagination, and VTEX limit failures.

## Phase 4: Application Wiring and Execution

- [x] 4.1 RED: Extend `internal/app/app_test.go` to expect six `masterdata.*` capabilities in the finalized registry.
- [x] 4.2 GREEN: Modify `internal/app/app.go` to build Master Data brand clients from resolved configuration and register all six capabilities.
- [x] 4.3 Add/extend execution or CLI run tests proving a Master Data capability emits envelope-safe data via `exito run`.

## Phase 5: Verification

- [x] 5.1 Run focused package tests for `internal/config`, `internal/domain/masterdata`, and `internal/app`.
- [x] 5.2 Run `make test`.
- [x] 5.3 Confirm no write, delete, schema lifecycle, or index mutation endpoints were introduced.
