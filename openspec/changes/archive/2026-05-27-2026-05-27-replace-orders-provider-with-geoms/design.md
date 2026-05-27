# Design

`orders.get` remains the public contract. `internal/domain/orders.HTTPGetter` now owns the GEOMS-specific DTOs and token flow so external API shapes do not leak to surfaces.

The getter first obtains a bearer token. If `EXITO_ORDERS_TOKEN` is present, it is treated as a pre-fetched fallback token. Otherwise, the getter posts Azure AD client-credentials form data (`client_id`, `client_secret`, `grant_type=client_credentials`, `scope=api://<scope>/.default`) to the configured token URL, defaulting to the GEOMS tenant URL. The token response `expires_in` controls the in-memory token cache, with a one-minute refresh margin.

The order lookup posts to `<orders.baseUrl>/findOrders` with GEOMS envelope fields and filters. `id` maps to `filters.orderNumber`; `orderType` defaults to `ExitoEcomm` and can be set to `CarullaEcomm`.

Configuration keeps non-sensitive GEOMS base URLs in YAML. Secrets remain in environment/dotenv via either `EXITO_ORDERS_CLIENT_ID`, `EXITO_ORDERS_CLIENT_SECRET`, `EXITO_ORDERS_SCOPE` or the existing GEOMS bundle variables (`GEOMS_CREDENTIALS_QA`/`GEOMS_CREDENTIALS_PDN`).
