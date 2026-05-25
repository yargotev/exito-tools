# Capability Contract Metadata Delta Specification

## ADDED Requirements

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
