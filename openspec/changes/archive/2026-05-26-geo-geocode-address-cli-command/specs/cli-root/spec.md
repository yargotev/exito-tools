# cli-root Delta

## ADDED Requirements

### Requirement: Geo geocode-address explicit command

The CLI SHALL expose `exito geo geocode-address --city <city> --address <address>` as the explicit command for the `geo.geocode-address` Capability.

#### Scenario: Command routes through shared execution

- **WHEN** a user runs `exito geo geocode-address --city Bogota --address "CL 57 H SUR # 68 D - 75"`
- **THEN** the CLI executes `geo.geocode-address` through the shared Pipeline
- **AND** emits a standard JSON envelope containing `meta.capabilityId` set to `geo.geocode-address`

#### Scenario: Required flags are enforced before execution

- **WHEN** either `--city` or `--address` is omitted
- **THEN** Cobra rejects the command before emitting a JSON envelope
