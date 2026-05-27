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

### Requirement: Orders get command

The CLI Surface MUST expose `exito orders get --id <order-id>` as an explicit domain command for the `orders.get` Capability and MAY accept an optional GEOMS order-type filter.

#### Scenario: Orders get runs through the shared pipeline

- GIVEN the Application has registered `orders.get`
- WHEN a user runs `exito orders get --id A123`
- THEN the command executes the `orders.get` Capability through the shared execution Pipeline
- AND stdout contains a standard JSON Envelope
- AND `meta.capabilityId` is `orders.get`

#### Scenario: Orders get supports Carulla order type

- **GIVEN** the Application has registered `orders.get`
- **WHEN** a user runs `exito orders get --id A123 --order-type CarullaEcomm`
- **THEN** the command executes the `orders.get` Capability through the shared execution Pipeline
- **AND** the Capability input contains `id` equal to `A123`
- **AND** the Capability input contains `orderType` equal to `CarullaEcomm`

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

### Requirement: CLI persists Default Profile explicitly

The CLI Surface MUST expose an explicit command for setting the saved Default Profile and MUST return a machine-readable JSON Envelope.

#### Scenario: Default profile command writes selected configuration

- GIVEN a user runs `exito config set-default-profile prod`
- WHEN the command succeeds
- THEN stdout contains a JSON Envelope with `ok: true`
- AND `data.profile` is `prod`
- AND `data.configPath` identifies the YAML Configuration File updated

#### Scenario: Blank profile is rejected before writing

- GIVEN a user runs `exito config set-default-profile "   "`
- WHEN the command validates input
- THEN no Configuration File is written
- AND the command fails with a user-facing validation error

### Requirement: Intelligent Search products command

The CLI Surface MUST expose `catalog.intelligent-search-products` through `exito catalog intelligent-search products`.

#### Scenario: CLI text search

- **Given** the CLI is available
- **When** a caller runs `exito catalog intelligent-search products --brand exito --trade-policy 1 --text leche`
- **Then** the command MUST execute `catalog.intelligent-search-products`
- **And** stdout MUST contain the standard JSON envelope.

#### Scenario: CLI typed SKU lookup

- **Given** the CLI is available
- **When** a caller runs `exito catalog intelligent-search products --brand exito --trade-policy 1 --by sku-id --value 912350`
- **Then** the command MUST execute `catalog.intelligent-search-products` with typed lookup input.

#### Scenario: CLI rejects missing trade policy

- **Given** the CLI is available
- **When** a caller runs `exito catalog intelligent-search products --text leche`
- **Then** the command MUST fail before provider execution because `--trade-policy` is required.
