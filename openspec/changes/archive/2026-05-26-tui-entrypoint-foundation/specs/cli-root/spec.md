# CLI Root Delta Specification

## ADDED Requirements

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
