# capability-registry Delta

## ADDED Requirements

### Requirement: Geo explicit command uses registered Capability

The explicit Geo CLI command SHALL invoke the registered `geo.geocode-address` Capability rather than duplicating use-case logic in the CLI surface.

#### Scenario: Metadata is propagated

- **WHEN** the command runs with `--profile` and `--correlation-id`
- **THEN** the output envelope metadata includes the selected profile, supplied correlation ID, and `geo.geocode-address` Capability ID
