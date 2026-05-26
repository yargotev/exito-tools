# CLI Root Delta

## ADDED Requirements

### Requirement: CLI persists Default Profile explicitly

The CLI Surface MUST expose an explicit command for setting the saved Default Profile and MUST return a machine-readable JSON Envelope.

#### Scenario: Default profile command writes selected configuration

- **GIVEN** a user runs `exito config set-default-profile prod`
- **WHEN** the command succeeds
- **THEN** stdout contains a JSON Envelope with `ok: true`
- **AND** `data.profile` is `prod`
- **AND** `data.configPath` identifies the YAML Configuration File updated

#### Scenario: Blank profile is rejected before writing

- **GIVEN** a user runs `exito config set-default-profile "   "`
- **WHEN** the command validates input
- **THEN** no Configuration File is written
- **AND** the command fails with a user-facing validation error
