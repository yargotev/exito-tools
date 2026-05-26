# Cli Run Specification

## Requirements

### Requirement: Generic run command executes registered capabilities

The CLI Surface MUST expose `exito run <capability-id>` as a machine-readable command that executes a registered Capability by stable Capability ID through the shared execution Pipeline.

#### Scenario: Registered capability runs through the Pipeline

- GIVEN an executable Capability is registered during application boot
- WHEN a user runs `exito run <capability-id>`
- THEN the Capability handler receives the parsed input and execution metadata
- AND stdout contains a standard JSON Envelope
- AND the command exits successfully when the envelope is successful

#### Scenario: Unknown capability returns structured failure envelope

- GIVEN no executable Capability exists for a requested ID
- WHEN a user runs `exito run missing.example`
- THEN stdout contains a failed JSON Envelope
- AND the structured error code is `CAPABILITY_NOT_FOUND`
- AND the command returns a generic failure exit status

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

### Requirement: Generic run command accepts explicit confirmation

The CLI Surface MUST let agents pass explicit confirmation to `exito run <capability-id>` for confirmation-required Capabilities, and MUST fail non-interactively when confirmation is missing.

#### Scenario: Missing run confirmation returns confirmation-required envelope

- GIVEN a registered Capability requires confirmation
- WHEN a user runs `exito run <capability-id>` without `--confirm`
- THEN stdout contains a failed JSON Envelope
- AND the structured error code is `CONFIRMATION_REQUIRED`
- AND the command returns a generic failure exit status

#### Scenario: Run confirmation is passed to the Pipeline

- GIVEN a registered Capability requires confirmation
- WHEN a user runs `exito run <capability-id> --confirm`
- THEN the Capability executes through the shared Pipeline
- AND stdout contains a successful JSON Envelope when the handler succeeds
