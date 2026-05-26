# Design: Orders provider configuration foundation

## Approach

Mirror the existing Geo provider configuration contract for Orders. `config.Effective` gains an `OrdersProvider` field with base URL, source metadata, token source metadata, token presence, and configured state. The token value remains available to application wiring but is excluded from JSON serialization.

The resolver uses the same profile-aware credential layers already defined for secrets:

1. process environment;
2. `.env.<profile>`;
3. `.env`.

Orders uses `EXITO_ORDERS_BASE_URL` and `EXITO_ORDERS_TOKEN`, matching the existing `EXITO_<DOMAIN>_*` environment naming convention established for Geo.

## Dependency Direction

The configuration resolver remains independent of Orders domain HTTP implementation. A later slice can consume `Effective.OrdersProvider` to construct an Orders HTTP getter.
