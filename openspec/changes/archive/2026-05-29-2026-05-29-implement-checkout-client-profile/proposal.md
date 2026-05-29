# Change: Implement Checkout client profile attachment

## Why

The Checkout roadmap has orderForm base and add-items support. Purchase assembly now needs the next safe-write step: attaching customer client profile data to an existing VTEX orderForm while keeping PII out of command output and diagnostics.

## Scope

- Add `checkout.update-client-profile` capability metadata, input validation, use case, and result model.
- Add VTEX Checkout HTTP client support for the client profile attachment endpoint.
- Register the capability in application boot wiring.
- Expose a CLI command: `exito checkout update-client-profile --brand <brand> --order-form-id <id> --input-json <profile-json> --confirm`.
- Add tests for validation, confirmation gating, provider request mapping, CLI JSON parsing, registration, and PII redaction by result shape.

## Out of scope

- Shipping data/logistics attachment.
- Guided purchase assembly workflow orchestration.
- Place-order, process-order, or payment endpoints.

## Safety

`checkout.update-client-profile` is a confirmation-required safe-write capability. Missing CLI confirmation must return `CONFIRMATION_REQUIRED` before the provider is called. Customer profile values are request-only sensitive inputs and must not be echoed in stdout JSON, logs, or diagnostics by default.
