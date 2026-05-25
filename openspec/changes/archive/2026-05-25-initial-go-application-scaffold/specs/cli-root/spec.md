# CLI Root Specification

## Purpose

Define root CLI behavior.

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

### Requirement: Root command preserves future discovery seams

The root command SHOULD expose CLI guidance without claiming deferred commands exist yet.

#### Scenario: Deferred commands are not misrepresented

- GIVEN this scaffold has no real business capability commands
- WHEN root help is shown
- THEN the help avoids implying that Orders, Geo, `run`, or `capabilities` already work
