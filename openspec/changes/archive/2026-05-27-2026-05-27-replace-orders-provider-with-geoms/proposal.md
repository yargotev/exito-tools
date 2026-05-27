# Replace Orders provider with GEOMS

## Why

`orders.get` currently targets a placeholder Orders provider contract. The real Orders integration uses GEOMS in two steps: obtain an Azure AD bearer token with client credentials, then call GEOMS `findOrders` with an order-number filter.

## Scope

- Replace the Orders HTTP getter internals with GEOMS token acquisition and `findOrders` calls.
- Keep the public `orders.get` Capability ID and `exito orders get --id` command stable.
- Add an optional GEOMS `orderType` input/CLI flag for Exito and Carulla lookups.
- Resolve GEOMS client credentials from environment/dotenv without committing secrets.
- Update non-sensitive profile base URLs for QA/staging and prod.

## Out of Scope

- New Orders capabilities beyond `orders.get`.
- Committing GEOMS secrets.
- Archiving this change before verification and explicit user request.
