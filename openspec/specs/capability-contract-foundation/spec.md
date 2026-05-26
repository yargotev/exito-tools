# Capability Contract Foundation Specification

## Purpose

Define shared foundation contracts.

## Requirements

### Requirement: Shared contract skeletons exist

The system MUST provide shared contract types for Capability metadata, Structured Errors, and JSON Envelope-shaped results so later slices can extend them without changing boundaries.

#### Scenario: Future capability can depend on shared types

- GIVEN a later slice adds a real capability
- WHEN it needs common metadata or failure shapes
- THEN it can use those shared contract types

#### Scenario: Runtime envelope emission remains deferred

- GIVEN the initial scaffold only implements root help
- WHEN no explicit machine-readable command runs
- THEN shared envelope types may exist without requiring JSON output

### Requirement: Visible contracts remain English-only

The system MUST keep user-facing labels and messages English-only.

#### Scenario: Shared CLI-facing text is English-only

- GIVEN the scaffold exposes root help or contract-facing messages
- WHEN a user reads them
- THEN the visible text is in English only

### Requirement: Capability definitions expose neutral input schemas

The system MUST model schema-first Capability input metadata as neutral contracts so interaction surfaces and agents can discover complete input object requirements without redefining them per surface.

#### Scenario: Capability input schema is serialized in inventory output

- GIVEN a Capability definition includes an input schema with fields
- WHEN the CLI Surface emits the machine-readable capabilities inventory
- THEN the JSON definition includes `inputSchema.fields` with field name, type, required marker, and description

### Requirement: Capability input schemas use shared primitive categories

The system MUST provide stable input type categories for common JSON-shaped values so future Operational Domains avoid ad-hoc type strings.

#### Scenario: Future domains use shared input type constants

- GIVEN a future domain registers a Capability with an input schema
- WHEN it declares a string, number, boolean, object, or array input field
- THEN it can use shared capability input type constants

### Requirement: Capability input schemas are enforceable at execution time

The system MUST be able to validate complete Capability input objects against neutral input schema metadata before a Capability handler runs.

#### Scenario: Required input field is missing

- GIVEN a Capability has an input schema with a required field
- WHEN execution receives an input object without that field
- THEN execution returns a structured `INVALID_INPUT` failure
- AND the Capability handler is not invoked

#### Scenario: Input field has the wrong type

- GIVEN a Capability has an input schema with a typed field
- WHEN execution receives an input object with an incompatible value type
- THEN execution returns a structured `INVALID_INPUT` failure
- AND the Capability handler is not invoked

#### Scenario: Capability has no input schema

- GIVEN a Capability does not define an input schema
- WHEN execution receives an input object
- THEN execution preserves existing behavior and invokes the Capability handler

### Requirement: Capability definitions expose discovery metadata

The system MUST model neutral Capability metadata for domain, version, risk, confirmation requirements, audiences, and visibility so interaction surfaces can discover and adapt capabilities consistently.

#### Scenario: Capability metadata is serialized in inventory output

- GIVEN a Capability definition includes domain, version, risk, audiences, and visibility
- WHEN the CLI Surface emits the machine-readable capabilities inventory
- THEN the JSON definition includes those metadata fields

### Requirement: Capability metadata uses documented categories

The system MUST provide stable metadata categories for read-only, safe-write, and destructive risk; agents and people audiences; and CLI, TUI, and command-palette visibility.

#### Scenario: Future domains use shared constants

- GIVEN a future domain registers a Capability
- WHEN it assigns risk, audience, or visibility metadata
- THEN it can use shared capability metadata constants instead of ad-hoc strings

### Requirement: Geo geocode-address Capability contract

The system SHALL define a neutral executable Capability with ID `geo.geocode-address` for geocoding a city/address pair.

#### Scenario: Capability metadata is exposed

- **WHEN** the Geo geocode-address Capability definition is inspected
- **THEN** it exposes domain `geo`, version `1.0.0`, read-only risk, agents and people audiences, CLI/TUI/command-palette visibility, and required string input fields `city` and `address`

#### Scenario: Missing Geo provider configuration is structured

- **WHEN** `geo.geocode-address` is executed before a real Geo provider client is configured
- **THEN** the execution returns a structured `GEO_NOT_CONFIGURED` error envelope

### Requirement: Geo provider HTTP mapping

The Geo Domain SHALL map configured provider geocode responses to the stable `geo.geocode-address` result contract without exposing provider DTOs to callers.

#### Scenario: Provider response is mapped to domain result

- **GIVEN** `EXITO_GEO_BASE_URL` and `EXITO_GEO_TOKEN` are configured
- **WHEN** `geo.geocode-address` receives a successful provider response containing `data.latitude`, `data.longitude`, `data.estado`, `data.dirtrad`, `data.barrio`, and `data.coddane`
- **THEN** the capability result contains `location.latitude`, `location.longitude`, `status`, `normalizedAddress`, `neighborhood`, and `daneCode`

### Requirement: Orders get capability exposes a neutral contract

The system MUST expose `orders.get` as a read-only Orders Domain Capability with schema-first input and domain-owned result models.

#### Scenario: Orders get definition is discoverable

- GIVEN `orders.get` is registered
- WHEN a surface inspects its Capability definition
- THEN the definition includes domain `orders`, version metadata, read-only risk, agents and people audiences, CLI/TUI/command-palette visibility, and a required string `id` input field

#### Scenario: Orders get default dependency is not configured

- GIVEN the Application has no real Orders API client yet
- WHEN `orders.get` is executed with valid input through the generic pipeline
- THEN execution returns a structured `ORDERS_NOT_CONFIGURED` failure envelope

### Requirement: orders.get domain execution

The system MUST execute `orders.get` through a domain-owned getter and return the stable `orders.GetResult` result shape.

#### Scenario: Configured Orders provider returns a mapped order

- **Given** the Orders provider is configured with a base URL and token
- **And** the provider returns `{"order":{"id":"A123","status":"created","createdAt":"2026-05-26T00:00:00Z"}}`
- **When** `orders.get` is executed with input `{"id":"A123"}`
- **Then** the capability succeeds
- **And** the data is an Orders domain result containing order ID `A123`, status `created`, and the provider `createdAt` value

#### Scenario: Orders provider returns not found

- **Given** the Orders provider is configured
- **And** the provider responds with HTTP 404
- **When** `orders.get` is executed for the missing ID
- **Then** the capability fails with structured error code `ORDER_NOT_FOUND`
