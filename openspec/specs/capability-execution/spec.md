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
