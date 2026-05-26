# CLI Root Specification

## Purpose

Define root CLI behavior and explicit domain command mappings.

## Requirements

### Requirement: Root command shows brief help

The system MUST make `exito` without arguments show brief English help instead of launching the TUI or emitting JSON output.

#### Scenario: Bare command shows help

- GIVEN the executable is available
- WHEN a user runs `exito`
- THEN brief English help text is shown
- AND the command does not enter the TUI

#### Scenario: Root help stays human-readable

- GIVEN no explicit machine-readable command was selected
- WHEN root help is rendered
- THEN the output is text help rather than JSON Envelope output

### Requirement: Root command advertises implemented commands

The root command SHOULD expose CLI guidance for implemented commands without launching interactive behavior.

#### Scenario: Implemented commands are visible

- GIVEN capabilities, generic run, Orders, and Geo commands are implemented
- WHEN root help is shown
- THEN help may advertise `capabilities`, `run`, `orders`, and `geo`
- AND the command still does not emit a JSON Envelope

### Requirement: Geo geocode-address explicit command

The CLI SHALL expose `exito geo geocode-address --city <city> --address <address>` as the explicit command for the `geo.geocode-address` Capability.

#### Scenario: Command routes through shared execution

- **WHEN** a user runs `exito geo geocode-address --city Bogota --address "CL 57 H SUR # 68 D - 75"`
- **THEN** the CLI executes `geo.geocode-address` through the shared Pipeline
- **AND** emits a standard JSON envelope containing `meta.capabilityId` set to `geo.geocode-address`

#### Scenario: Required flags are enforced before execution

- **WHEN** either `--city` or `--address` is omitted
- **THEN** Cobra rejects the command before emitting a JSON envelope

### Requirement: Orders get command is exposed explicitly

The CLI Surface MUST expose `exito orders get --id <order-id>` as an explicit domain command for the `orders.get` Capability.

#### Scenario: Orders get runs through the shared pipeline

- GIVEN the Application has registered `orders.get`
- WHEN a user runs `exito orders get --id A123`
- THEN the command executes the `orders.get` Capability through the shared execution Pipeline
- AND stdout contains a standard JSON Envelope
- AND `meta.capabilityId` is `orders.get`

#### Scenario: Orders get requires an ID flag

- GIVEN the Orders get command is available
- WHEN a user runs `exito orders get` without `--id`
- THEN the command fails before execution
- AND it does not emit a JSON success envelope

### Requirement: TUI starts only through an explicit command

The CLI Surface MUST expose an explicit `exito tui` command for the interactive TUI and MUST NOT launch the TUI from bare `exito`.

#### Scenario: Root help advertises the TUI command

- GIVEN the CLI root is executed without arguments
- WHEN help is rendered
- THEN the help output includes `tui`
- AND the output remains textual help rather than a JSON envelope or interactive TUI session

#### Scenario: TUI command uses shared boot flags

- GIVEN a user supplies `--config` or `--profile`
- WHEN the user runs `exito tui`
- THEN the TUI command bootstraps the Application with those boot flags
