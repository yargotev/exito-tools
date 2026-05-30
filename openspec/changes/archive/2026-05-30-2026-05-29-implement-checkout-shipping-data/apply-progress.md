# Apply Progress: Implement Checkout Shipping Data

## Status

All tasks complete.

## Completed Tasks

- [x] 1.1 Domain tests for shipping normalization, redaction, invalid input.
- [x] 1.2 Domain models/use case/capability definition/input mapping/validation.
- [x] 2.1 HTTP tests for `/attachments/shippingData`, body mapping, coordinates, diagnostics.
- [x] 2.2 Brand routing and HTTP `UpdateShippingData` implementation.
- [x] 2.3 App registration test and wiring.
- [x] 3.1 CLI confirmation and JSON/redaction tests.
- [x] 3.2 CLI `checkout update-shipping-data` command.
- [x] 4.1 Focused tests, `make test`, `make precommit`, and `go build ./cmd/exito`.
- [x] 4.2 Apply progress saved.

## TDD Cycle Evidence

| Task | Test File | Layer | Safety Net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|-----------|-------|------------|-----|-------|-------------|----------|
| 1.1/1.2 | `internal/domain/checkout/checkout_test.go` | Unit | ✅ `go test ./internal/domain/checkout` baseline passed | ✅ Missing shipping types/use case failed compile | ✅ `TestUpdateShippingData*` passed | ✅ Happy path + invalid input | ✅ gofumpt/precommit |
| 2.1/2.2 | `internal/domain/checkout/http_client_test.go` | Unit | ✅ Checkout package baseline passed | ✅ Missing `HTTPClient.UpdateShippingData` failed compile | ✅ HTTP shipping tests passed | ✅ Request body + existing redaction regression | ✅ Extracted shipping total/set helpers |
| 2.3 | `internal/app/app_test.go` | Unit | ✅ `go test ./internal/app` baseline passed | ✅ Capability registration expectation failed before wiring | ✅ App tests passed | ➖ Structural registration | ✅ gofumpt/precommit |
| 3.1/3.2 | `internal/surface/cli/root_test.go` | Unit | ✅ `go test ./internal/surface/cli` baseline passed | ✅ Unknown command/flag and no JSON envelope | ✅ CLI shipping tests passed | ✅ Confirmation failure + confirmed JSON adapter/redaction | ✅ Mirrored client-profile command pattern |
| 4.1/4.2 | Full repo | Verification | ✅ Focused tests passed | N/A | ✅ `make test`, `make precommit`, `go build ./cmd/exito` passed | N/A | ✅ No lint issues |

## Verification

- `go test ./internal/domain/checkout` ✅
- `go test ./internal/app ./internal/surface/cli` ✅
- `make test` ✅
- `make precommit` ✅
- `go build ./cmd/exito` ✅

## Deviations

None.

## Remaining

- Archive/sync this OpenSpec change after review.
- Commit staged changes when requested.
