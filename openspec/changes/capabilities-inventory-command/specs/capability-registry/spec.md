# Capability Registry Delta Specification

## ADDED Requirements

### Requirement: Registry inventory is discoverable by surfaces

The system MUST expose a stable snapshot of registered capabilities that surfaces can serialize without mutating the registry.

#### Scenario: CLI reads immutable inventory

- GIVEN application boot has finalized the registry
- WHEN the CLI capabilities command reads all definitions
- THEN it receives a defensive-copy snapshot suitable for JSON inventory output
