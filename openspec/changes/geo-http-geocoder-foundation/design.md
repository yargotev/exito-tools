# Design: Geo HTTP geocoder foundation

## Approach

Keep the client near the Geo Domain because the provider DTO mapping is domain-specific:

- Add `geo.HTTPGeocoder` implementing the existing `geo.Geocoder` interface.
- Build the endpoint by appending `/geocode-address` to `EXITO_GEO_BASE_URL`.
- Send a `POST` JSON request with `city` and `address`, plus `Authorization: Bearer <token>`.
- Decode only the selected fields documented in `docs/capabilities/geo.geocode-address.md`.
- Map non-2xx provider responses to `GEO_PROVIDER_UNAVAILABLE` and malformed JSON to `GEO_PROVIDER_INVALID_RESPONSE`.
- Keep missing/invalid configuration mapped to `GEO_NOT_CONFIGURED`.
- In application boot, choose `HTTPGeocoder` only when the resolved `GeoProvider` is configured; otherwise keep `UnavailableGeocoder` so startup never requires credentials.

This slice intentionally uses the standard `net/http` client directly. Future shared HTTP infrastructure can wrap timeouts, retries, request metadata, and authentication policies without changing the `geo.Geocoder` interface.
