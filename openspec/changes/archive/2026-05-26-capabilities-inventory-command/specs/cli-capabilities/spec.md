# CLI Capabilities Delta Specification

## ADDED Requirements

### Requirement: Capabilities command emits machine-readable inventory

The system MUST expose `exito capabilities` as a machine-readable CLI command that returns the finalized Capability Registry inventory using the standard JSON Envelope shape.

#### Scenario: Empty registry is returned successfully

- GIVEN no real business capabilities are registered yet
- WHEN a user runs `exito capabilities`
- THEN stdout contains JSON with `ok: true`
- AND `data.capabilities` is an empty array
- AND the command does not render root help text

#### Scenario: Boot flags affect capability inventory context

- GIVEN explicit `--config` and `--profile` flags are provided
- WHEN a user runs `exito capabilities`
- THEN the application is booted with those configuration inputs before rendering the inventory

### Requirement: Root help reflects only implemented commands

The system MUST keep bare `exito` as human-readable help while allowing implemented subcommands to appear in help.

#### Scenario: Deferred commands remain absent

- GIVEN `capabilities` is implemented but Orders, Geo, and generic run are deferred
- WHEN root help is rendered
- THEN `capabilities` may appear
- AND `orders`, `geo`, and `run` are not advertised as working commands
