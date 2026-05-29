# Change: Implement Checkout add-items

## Why

The Checkout roadmap has a working orderForm base, but purchase assembly needs the next safe-write step: adding selected Catalog SKU IDs to an existing VTEX orderForm.

## Scope

- Add `checkout.add-items` capability metadata, input validation, use case, and result model.
- Add VTEX Checkout HTTP client support for adding cart items to an orderForm.
- Register the capability in application boot wiring.
- Expose a CLI command: `exito checkout add-items --brand <brand> --order-form-id <id> --item sku=<sku>,quantity=<qty>[,seller=<seller>] --confirm`.
- Add tests for validation, confirmation gating, provider request mapping, CLI parsing, and registration.

## Out of scope

- Client profile attachment.
- Shipping data/logistics attachment.
- Guided purchase assembly workflow orchestration.
- Place-order, process-order, or payment endpoints.

## Safety

`checkout.add-items` is a confirmation-required safe-write capability. Missing CLI confirmation must return `CONFIRMATION_REQUIRED` before the provider is called. The response remains a redacted orderForm summary.
