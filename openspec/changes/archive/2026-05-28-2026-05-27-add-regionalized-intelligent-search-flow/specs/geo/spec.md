## MODIFIED Requirements

### Requirement: VTEX region coverage diagnostics

#### Scenario: Resolve coverage from coordinates

- **Given** VTEX public checkout regions is configured for the requested brand
- **When** a caller executes `geo.resolve-vtex-region` with `country`, `salesChannel`, `longitude`, and `latitude`
- **Then** the capability MUST call `GET /api/checkout/pub/regions`
- **And** the request query MUST include `country={country}` and `sc={salesChannel}`
- **And** the request query MUST include `geoCoordinates={longitude};{latitude}` preserving longitude before latitude
- **And** the result MUST include returned sellers, resolved region IDs when present, and region diagnostics.
