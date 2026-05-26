# Proposal: Geo HTTP geocoder foundation

## Why

`geo.geocode-address` currently returns `GEO_NOT_CONFIGURED` even when Geo provider credentials are resolved. The next narrow slice should wire configured Geo provider settings into a real HTTP-backed geocoder while preserving the domain-owned result contract.

## What

- Add a Geo Domain HTTP geocoder that posts city/address input to the configured provider endpoint.
- Map provider DTO fields into the existing `GeocodeAddressResult` shape.
- Return stable structured errors for missing configuration, provider failures, and invalid provider responses.
- Wire application boot to use the HTTP geocoder only when `Config.GeoProvider.Configured` is true; otherwise keep the unavailable geocoder.

## Out of Scope

- Shared retry/backoff policies.
- Orders provider wiring.
- TUI integration beyond existing capability metadata.
