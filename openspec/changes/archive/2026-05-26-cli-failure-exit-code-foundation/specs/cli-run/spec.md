# CLI Run Delta Specification

## MODIFIED Requirements

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
