# Orders Specification

## Purpose

Define Orders Domain capabilities and provider behavior.

## Requirements

### Requirement: Independent VTEX OMS order lookup

Exito Tools MUST expose an independent `orders.get-vtex` capability for VTEX OMS order detail lookups without changing the existing GEOMS-backed `orders.get` capability.

#### Scenario: Query VTEX OMS by order id

- **Given** VTEX OMS is configured for the selected profile and brand
- **When** a caller executes `orders.get-vtex` with an `id`
- **Then** Exito Tools MUST call VTEX OMS order detail for that id
- **And** the result MUST be returned in the standard JSON envelope
- **And** the capability metadata MUST use `orders.get-vtex`

#### Scenario: Orders TUI exposes both provider lookups

- **Given** the TUI Surface renders people-facing Orders Domain primary actions
- **When** Orders capabilities are discovered from the shared registry
- **Then** `orders.get` MUST appear as `Get GEOMS order`
- **And** `orders.get-vtex` MUST appear as `Get VTEX OMS order`

### Requirement: VTEX OMS credentials stay server-side

Exito Tools MUST resolve VTEX OMS app key and app token only from process environment or non-committed dotenv files, and MUST NOT expose their values in JSON-serialized effective configuration.

#### Scenario: Serialize effective configuration

- **Given** VTEX OMS app key and app token are configured
- **When** effective configuration is serialized
- **Then** the app key and app token values MUST be omitted
- **And** only credential presence/source metadata MAY be exposed
