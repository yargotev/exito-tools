# Capability Execution Specification

## Requirements

### Requirement: Capability execution pipeline wraps handler results

The system MUST provide a surface-independent Capability execution pipeline that invokes registered Capability handlers with context and returns the standard JSON Envelope shape.

#### Scenario: Registered capability succeeds

- GIVEN a Capability with a context-aware handler is registered
- WHEN the execution pipeline runs that Capability ID with input
- THEN the handler receives the input and execution context
- AND the returned envelope has `ok: true`
- AND `data` contains the handler result
- AND metadata includes request ID, optional correlation ID, profile, capability ID, and duration.

#### Scenario: Registered capability returns structured error

- GIVEN a Capability handler returns a structured error
- WHEN the execution pipeline runs that Capability ID
- THEN the returned envelope has `ok: false`
- AND the envelope error preserves the structured error code and message
- AND standard metadata is still included.

#### Scenario: Capability ID is unknown

- GIVEN no Capability is registered for an ID
- WHEN the execution pipeline runs that Capability ID
- THEN the returned envelope has `ok: false`
- AND the error code is `CAPABILITY_NOT_FOUND`
- AND standard metadata is still included.

### Requirement: Successful capability execution propagates warnings

The Capability execution Pipeline MUST include non-fatal warnings returned by a successful handler in Envelope metadata without changing the success state.

#### Scenario: Successful handler returns warning metadata

- GIVEN a registered Capability handler returns successful data and one Structured Warning
- WHEN the execution Pipeline runs the Capability
- THEN the returned envelope has `ok: true`
- AND `meta.warnings` contains the warning
- AND the warning does not change the successful data result

#### Scenario: Handler failure does not expose result warnings

- GIVEN a registered Capability handler returns an error
- WHEN the execution Pipeline runs the Capability
- THEN the returned envelope has `ok: false`
- AND no successful result warnings are exposed in metadata

### Requirement: Successful capability execution propagates pagination metadata

The Capability execution Pipeline MUST include pagination metadata returned by a successful handler in Envelope metadata without changing the success data shape.

#### Scenario: Successful handler returns pagination metadata

- GIVEN a registered Capability handler returns successful list data and pagination metadata
- WHEN the execution Pipeline runs the Capability
- THEN the returned envelope has `ok: true`
- AND `meta.pagination` contains the pagination metadata
- AND the successful data result remains under `data`

#### Scenario: Handler failure does not expose result pagination

- GIVEN a registered Capability handler returns an error
- WHEN the execution Pipeline runs the Capability
- THEN the returned envelope has `ok: false`
- AND no successful result pagination metadata is exposed in metadata

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
