# Orders HTTP getter foundation

## Summary

Wire `orders.get` to a real provider-backed HTTP getter when Orders provider configuration is present.

## Scope

- Add an Orders domain HTTP getter that calls the configured provider through shared HTTP infrastructure.
- Map provider DTOs to the domain-owned `orders.Order` model.
- Translate provider configuration, transport, non-success, not-found, and invalid-response failures into structured domain errors.
- Wire `app.New` to use the HTTP getter only when `Effective.OrdersProvider.Configured` is true.

## Out of Scope

- Changing the public `orders.get` input or result contract.
- Adding provider-specific business fields beyond the documented initial order shape.
- Adding retries, pagination, or write capabilities.
