# Proposal: Implement Checkout Shipping Data

## Intent

Complete the next Checkout purchase-assembly slice by adding a confirmation-gated capability to attach delivery address data and selected logistics options to an existing VTEX orderForm.

## Scope

### In Scope
- Add `checkout.update-shipping-data` capability and CLI command.
- Map shipping address and logistics selections to VTEX Checkout `shippingData` attachment.
- Return only the redacted orderForm summary plus safe shipping diagnostics.
- Cover confirmation, validation, app wiring, HTTP mapping, and JSON envelope tests.

### Out of Scope
- Final order placement, process-order, payment, or payment attachment work.
- A full guided purchase macro across Catalog and Checkout.
- Persisting customer PII or exposing raw address data in stdout/logs.

## Capabilities

### New Capabilities
None

### Modified Capabilities
- `checkout`: implement the already-specified `checkout.update-shipping-data` behavior.

## Approach

Follow the existing Checkout pattern: add domain input/result/use case/definition, extend `BrandClient` and `HTTPClient`, register in `internal/app`, and add `exito checkout update-shipping-data --input-json ... --confirm` as a surface adapter over the shared Pipeline.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/domain/checkout` | Modified | Shipping data models, validation, HTTP attachment call, tests |
| `internal/app/app.go` | Modified | Capability registration and provider interface extension |
| `internal/surface/cli/root.go` | Modified | New Checkout subcommand and CLI tests |
| `openspec/changes/2026-05-29-implement-checkout-shipping-data` | New | SDD artifacts for this slice |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| PII leaks in CLI output | Medium | Return redacted orderForm summary only; tests assert raw address is not echoed |
| Incorrect VTEX coordinates order | Medium | Model input as VTEX `geoCoordinates` order: longitude, latitude; test request body |
| Parallel/unsafe Checkout writes | Low | Capability remains safe-write and confirmation-required |

## Rollback Plan

Revert the commit/files for this change; removing registration and CLI command returns Checkout to the previous client-profile slice.

## Dependencies

- Existing VTEX Checkout base URL configuration.
- Caller must know the desired address/logistics selections from prior orderForm inspection/simulation.

## Success Criteria

- [ ] `checkout.update-shipping-data` appears in capability discovery.
- [ ] CLI requires `--confirm` before provider mutation.
- [ ] HTTP client posts to `/attachments/shippingData` with address/logistics body.
- [ ] Tests pass with `make test`.
