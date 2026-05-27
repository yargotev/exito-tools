# Design: VTEX region coverage diagnostics

## Approach

The capability belongs to the Geo Domain because it answers a coordinate coverage question. It uses VTEX Checkout Regions as a read-only provider, wired from the existing non-sensitive VTEX public base URL configuration per brand.

## Capability contract

- ID: `geo.resolve-vtex-region`
- Risk: read-only
- Inputs:
  - `brand`: `exito` or `carulla`, default `exito`
  - `country`: default `COL`
  - `salesChannel`: required
  - `longitude`: required string/number shaped as text
  - `latitude`: required string/number shaped as text
- Output:
  - brand, country, salesChannel, coordinates
  - `hasCoverage`
  - sellers with stable summary fields
  - diagnostics with request path/query and provider payload

## URL construction

The HTTP adapter must preserve VTEX's coordinate ordering exactly as `longitude;latitude` in `geoCoordinates`.

## Safety

The implementation only performs an HTTP GET and does not accept or emit secrets. No session token or cookie handling is introduced in this phase.
