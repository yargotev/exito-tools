# Capability Execution Delta: confirmation required foundation

## Requirements

### Requirement: Pipeline enforces confirmation-required capabilities

The Capability execution Pipeline MUST reject confirmation-required Capabilities unless the execution request carries explicit confirmation.

#### Scenario: Missing confirmation returns structured failure

- GIVEN a registered Capability has `requiresConfirmation` set to true
- WHEN the Pipeline executes it without explicit confirmation
- THEN the handler is not called
- AND the returned envelope has `ok: false`
- AND the structured error code is `CONFIRMATION_REQUIRED`
- AND standard metadata is still included

#### Scenario: Provided confirmation allows execution

- GIVEN a registered Capability has `requiresConfirmation` set to true
- WHEN the Pipeline executes it with explicit confirmation
- THEN the Capability handler is called
- AND the returned envelope reflects the handler result
