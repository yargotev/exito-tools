# Capability Registry Specification

## Purpose

Define boot-time registry behavior.

## Requirements

### Requirement: Registry accepts boot-time registration

The system MUST support explicit capability registration during application boot so Operational Domains can add capabilities through visible wiring.

#### Scenario: Orders get is available after finalization

- GIVEN application boot wires the Orders Domain
- WHEN the registry is finalized
- THEN `orders.get` is discoverable by stable Capability ID

### Requirement: Registry becomes immutable after boot

The system MUST treat the Capability Registry as immutable after finalization, including slice-backed metadata inside registered Capability definitions.

#### Scenario: Returned input schema fields cannot mutate registry state

- GIVEN a finalized registry contains a Capability with input schema fields
- WHEN a caller mutates the fields returned from `All`
- THEN a later `All` call returns the original input schema fields

#### Scenario: Post-boot mutation is rejected

- GIVEN the registry has been finalized
- WHEN code attempts another registration
- THEN the attempt is rejected with a stable failure outcome

#### Scenario: Returned definition metadata cannot mutate registry state

- GIVEN a finalized registry contains a Capability with audience or visibility metadata
- WHEN a caller mutates the slices returned from `All`
- THEN a later `All` call returns the original metadata

### Requirement: Registry inventory is discoverable by surfaces

The system MUST expose a stable snapshot of registered capabilities that surfaces can serialize without mutating the registry.

#### Scenario: CLI reads immutable inventory

- GIVEN application boot has finalized the registry
- WHEN the CLI capabilities command reads all definitions
- THEN it receives a defensive-copy snapshot suitable for JSON inventory output

### Requirement: Registry supports executable capability lookup

The system MUST allow surfaces and execution code to look up immutable registered Capability entries by stable Capability ID.

#### Scenario: Registered executable capability is found

- GIVEN application wiring registered an executable Capability before finalization
- WHEN the finalized registry is queried by Capability ID
- THEN the matching immutable entry is returned.

#### Scenario: Duplicate capability IDs are rejected

- GIVEN a Capability ID is already registered during boot
- WHEN wiring attempts to register another Capability with the same ID
- THEN the registry rejects the duplicate with a stable failure outcome.

### Requirement: Application boot registers Geo geocode-address

The application boot process SHALL explicitly register `geo.geocode-address` in the immutable Capability Registry.

#### Scenario: Generic run can route to Geo geocode-address

- **WHEN** the CLI runs `exito run geo.geocode-address` with complete `city` and `address` input
- **THEN** the shared execution Pipeline routes to the registered Capability and emits a standard JSON envelope with `meta.capabilityId` set to `geo.geocode-address`

### Requirement: Geo explicit command uses registered Capability

The explicit Geo CLI command SHALL invoke the registered `geo.geocode-address` Capability rather than duplicating use-case logic in the CLI surface.

#### Scenario: Metadata is propagated

- **WHEN** the command runs with `--profile` and `--correlation-id`
- **THEN** the output envelope metadata includes the selected profile, supplied correlation ID, and `geo.geocode-address` Capability ID

### Requirement: Application boot wires configured Geo provider

The application boot process SHALL use the configured Geo HTTP provider for `geo.geocode-address` when the resolved Geo provider configuration is complete.

#### Scenario: Configured Geo provider succeeds through registry execution

- **GIVEN** `EXITO_GEO_BASE_URL` and `EXITO_GEO_TOKEN` resolve to non-blank values
- **WHEN** the registered `geo.geocode-address` Capability is executed
- **THEN** the execution is routed to the HTTP geocoder instead of the unavailable geocoder
- **AND** the standard envelope contains the mapped Geo result
