# Design: Geo geocode-address Capability foundation

## Approach

Mirror the established Orders Domain capability pattern with a narrow Geo Domain package:

- `internal/domain/geo` owns `GeocodeAddressInput`, `GeocodeAddressResult`, `Location`, the `Geocoder` interface, and use case.
- `NewGeocodeAddressCapability` adapts the use case into the neutral executable Capability contract.
- The Capability definition uses ID `geo.geocode-address`, domain `geo`, version `1.0.0`, read-only risk, agents/people audiences, CLI/TUI/command-palette visibility, and required string fields `city` and `address`.
- `UnavailableGeocoder` is the default boot dependency and returns `GEO_NOT_CONFIGURED` so the capability is discoverable and executable through generic run without pretending an external client exists.
- `app.New` registers the capability explicitly after `orders.get`.

## Deferred Work

A future slice should add the explicit `exito geo geocode-address --city <city> --address <address>` command. A separate infrastructure slice should add provider configuration and HTTP mapping.
