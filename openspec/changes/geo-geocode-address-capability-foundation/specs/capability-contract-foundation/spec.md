# capability-contract-foundation Delta

## ADDED Requirements

### Requirement: Geo geocode-address Capability contract

The system SHALL define a neutral executable Capability with ID `geo.geocode-address` for geocoding a city/address pair.

#### Scenario: Capability metadata is exposed

- **WHEN** the Geo geocode-address Capability definition is inspected
- **THEN** it exposes domain `geo`, version `1.0.0`, read-only risk, agents and people audiences, CLI/TUI/command-palette visibility, and required string input fields `city` and `address`

#### Scenario: Missing Geo provider configuration is structured

- **WHEN** `geo.geocode-address` is executed before a real Geo provider client is configured
- **THEN** the execution returns a structured `GEO_NOT_CONFIGURED` error envelope
