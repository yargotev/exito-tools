# Geo Specification

## Purpose

Define Geo Domain capabilities for geocoding and read-only regional coverage diagnostics.

## Requirements

### Requirement: VTEX region coverage diagnostics

Exito Tools MUST expose a read-only `geo.resolve-vtex-region` capability for VTEX Checkout Regions coverage diagnostics from known coordinates.

#### Scenario: Resolve coverage from coordinates

- **Given** VTEX public checkout regions is configured for the requested brand
- **When** a caller executes `geo.resolve-vtex-region` with `country`, `salesChannel`, `longitude`, and `latitude`
- **Then** the capability MUST call `GET /api/checkout/pub/regions`
- **And** the request query MUST include `country={country}` and `sc={salesChannel}`
- **And** the request query MUST include `geoCoordinates={longitude};{latitude}` preserving longitude before latitude
- **And** the result MUST include returned sellers and region diagnostics.

#### Scenario: Coverage is true when any seller is present

- **Given** VTEX Checkout Regions returns sellers for the requested coordinates
- **When** at least one returned seller has a non-empty `id`
- **Then** `hasCoverage` MUST be `true`.

#### Scenario: Coverage is false for no sellers

- **Given** VTEX Checkout Regions returns no sellers
- **When** `geo.resolve-vtex-region` succeeds
- **Then** `hasCoverage` MUST be `false`.

#### Scenario: Coverage rule preserves historical context

- **Given** previous Exito storefront logic treated coverage as true only when a returned seller ID differed from the requested brand/account
- **When** `geo.resolve-vtex-region` reports coverage diagnostics
- **Then** Exito Tools MUST document that this CLI uses the broader VTEX Regions interpretation where any returned seller counts as coverage
- **And** the returned sellers MUST remain available for future business-specific interpretation.

#### Scenario: Region diagnostics are read-only

- **Given** a caller executes `geo.resolve-vtex-region`
- **When** the capability contacts VTEX
- **Then** it MUST NOT write Checkout orderForm shipping data
- **And** it MUST NOT patch Master Data addresses
- **And** it MUST NOT create or update VTEX sessions or segments.

#### Scenario: CLI command emits JSON envelope

- **Given** the CLI is available
- **When** a caller runs `exito geo resolve-vtex-region --brand exito --country COL --sales-channel 1 --longitude -74.160580822 --latitude 4.598090587`
- **Then** the command MUST execute `geo.resolve-vtex-region`
- **And** the command MUST emit the standard JSON envelope on stdout.
