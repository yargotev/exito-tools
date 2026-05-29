# Apply Progress: Checkout orderForm base

## Implementation commit

- `d0ccf4c feat: add checkout orderform base`

## Completed implementation

- Added `internal/domain/checkout` with domain-owned use cases, inputs, result summaries, validation, and capability metadata.
- Added `checkout.get-order-form` and `checkout.create-order-form` capability definitions and handlers.
- Added brand dispatch for `exito` and `carulla` providers.
- Added VTEX Checkout HTTP client for `GET /api/checkout/pub/orderForm/{orderFormId}` and current-cart creation with `forceNewCart=true` and `sc`.
- Added fake-server tests for provider request path/query, request metadata propagation, DTO mapping, and safe summary behavior.
- Added configuration resolver support for `vtexCheckout` profile/brand base URLs and environment overrides.
- Added minimal CLI commands under `exito checkout`.
- Registered Checkout capabilities in application boot wiring.

## Scope guard

No code was added for:

- `checkout.add-items`
- `checkout.update-client-profile`
- `checkout.update-shipping-data`
- final order placement
- payment processing
