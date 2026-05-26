# CLI Run Delta Specification

## ADDED Requirements

### Requirement: Generic run command executes registered capabilities

The CLI Surface MUST expose `exito run <capability-id>` as a machine-readable command that executes a registered Capability by stable Capability ID through the shared execution Pipeline.

#### Scenario: Registered capability runs through the Pipeline

- GIVEN an executable Capability is registered during application boot
- WHEN a user runs `exito run <capability-id>`
- THEN the Capability handler receives the parsed input and execution metadata
- AND stdout contains a standard JSON Envelope

#### Scenario: Unknown capability returns structured failure envelope

- GIVEN no executable Capability exists for a requested ID
- WHEN a user runs `exito run missing.example`
- THEN stdout contains a failed JSON Envelope
- AND the structured error code is `CAPABILITY_NOT_FOUND`

### Requirement: Generic run command accepts complete JSON input objects

The CLI Surface MUST adapt complete JSON input objects into neutral `capability.Input` for the generic run path.

#### Scenario: Inline JSON input is accepted

- GIVEN `--input-json` contains a JSON object
- WHEN a user runs `exito run <capability-id> --input-json '{"id":"123"}'`
- THEN the object is passed to the Capability handler

#### Scenario: File JSON input is accepted

- GIVEN `--input-file` points to a JSON object file
- WHEN a user runs `exito run <capability-id> --input-file input.json`
- THEN the object is passed to the Capability handler

#### Scenario: Piped stdin JSON input is accepted

- GIVEN stdin contains a JSON object
- WHEN a user pipes it into `exito run <capability-id>`
- THEN the object is passed to the Capability handler
