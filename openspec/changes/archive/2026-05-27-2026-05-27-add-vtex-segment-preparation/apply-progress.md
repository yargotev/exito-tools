# Apply Progress

## Implemented

- Added `catalog.create-vtex-segment` as a confirmation-required safe-write Catalog capability.
- Added domain input/result/use case, brand dispatcher, unavailable creator, and HTTP VTEX Sessions adapter.
- Wired the capability in `internal/app` using existing public VTEX Catalog brand base URLs.
- Added CLI command `exito catalog create-vtex-segment` with required `--region-id`, required `--sales-channel`, `--confirm`, and optional `--include-cookie`.
- Added docs at `docs/capabilities/catalog.create-vtex-segment.md`.
- Added unit tests for metadata/confirmation, input validation, HTTP request payload, token redaction, app registration, and CLI confirmation behavior.

## Verification

- `go test ./...` passed.
