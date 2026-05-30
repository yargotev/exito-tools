# Tasks: Implement Checkout Shipping Data

## Phase 1: Domain Contract

- [x] 1.1 Add failing domain tests in `internal/domain/checkout/checkout_test.go` for shipping input normalization, redaction, and invalid input.
- [x] 1.2 Add `checkout.update-shipping-data` domain models, updater interface, use case, definition, capability input mapping, and validation in `internal/domain/checkout/checkout.go`.

## Phase 2: Provider and Wiring

- [x] 2.1 Add failing HTTP tests in `internal/domain/checkout/http_client_test.go` for POST `/attachments/shippingData`, body mapping, coordinate order, and diagnostics.
- [x] 2.2 Implement shipping updater routing in `internal/domain/checkout/brand_client.go` and `internal/domain/checkout/http_client.go`.
- [x] 2.3 Add failing app registration expectation in `internal/app/app_test.go` and register the capability in `internal/app/app.go`.

## Phase 3: CLI Surface

- [x] 3.1 Add failing CLI tests in `internal/surface/cli/root_test.go` for confirmation gating and `--input-json` adaptation/redacted output.
- [x] 3.2 Implement `exito checkout update-shipping-data` in `internal/surface/cli/root.go` using the shared Pipeline and JSON envelope.

## Phase 4: Verification

- [x] 4.1 Run focused Checkout/app/CLI tests, then `make test`.
- [x] 4.2 Save apply progress with TDD evidence and report remaining archive/commit steps.
