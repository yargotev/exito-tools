# CLI Run Delta: confirmation required foundation

## Requirements

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
