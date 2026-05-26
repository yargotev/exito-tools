# CLI Capabilities Delta Specification

## MODIFIED Requirements

### Requirement: Capabilities command emits machine-readable inventory

The system MUST expose `exito capabilities` as a machine-readable CLI command that returns the finalized Capability Registry inventory using the standard JSON Envelope shape with request metadata.

#### Scenario: Empty registry is returned successfully

- GIVEN no real business capabilities are registered yet
- WHEN a user runs `exito capabilities`
- THEN stdout contains JSON with `ok: true`
- AND `data.capabilities` is an empty array
- AND `meta.requestId` is present
- AND `meta.durationMs` is present
- AND the command does not render root help text

#### Scenario: Correlation ID is propagated

- GIVEN a `--correlation-id` flag is provided
- WHEN a user runs `exito capabilities`
- THEN `meta.correlationId` matches the supplied value
