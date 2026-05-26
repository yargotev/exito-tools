## ADDED Requirements

### Requirement: Application boot wires configured Geo provider

The application boot process SHALL use the configured Geo HTTP provider for `geo.geocode-address` when the resolved Geo provider configuration is complete.

#### Scenario: Configured Geo provider succeeds through registry execution

- **GIVEN** `EXITO_GEO_BASE_URL` and `EXITO_GEO_TOKEN` resolve to non-blank values
- **WHEN** the registered `geo.geocode-address` Capability is executed
- **THEN** the execution is routed to the HTTP geocoder instead of the unavailable geocoder
- **AND** the standard envelope contains the mapped Geo result
