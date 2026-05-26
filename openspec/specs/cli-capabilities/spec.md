# CLI Capabilities Specification

## Requirements

### Requirement: Capabilities command emits machine-readable inventory

The system MUST expose `exito capabilities` as a machine-readable CLI command that returns the finalized Capability Registry inventory using the standard JSON Envelope shape with request metadata.

#### Scenario: Registered inventory is returned successfully

- GIVEN application boot has registered implemented capabilities
- WHEN a user runs `exito capabilities`
- THEN stdout contains JSON with `ok: true`
- AND `data.capabilities` contains the finalized capability definitions
- AND `meta.requestId` is present
- AND `meta.durationMs` is present
- AND the command does not render root help text

#### Scenario: Boot flags affect capability inventory context

- GIVEN explicit `--config` and `--profile` flags are provided
- WHEN a user runs `exito capabilities`
- THEN the application is booted with those configuration inputs before rendering the inventory

#### Scenario: Correlation ID is propagated

- GIVEN a `--correlation-id` flag is provided
- WHEN a user runs `exito capabilities`
- THEN `meta.correlationId` matches the supplied value

### Requirement: Root help reflects implemented commands

The system MUST keep bare `exito` as human-readable help while allowing implemented subcommands to appear in help.

#### Scenario: Implemented commands are advertised

- GIVEN `capabilities`, `run`, `orders`, and `geo` are implemented
- WHEN root help is rendered
- THEN those implemented commands may appear in help
- AND root help remains human-readable text rather than a JSON Envelope
