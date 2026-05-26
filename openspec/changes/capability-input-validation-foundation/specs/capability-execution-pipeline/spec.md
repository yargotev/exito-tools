# Capability Execution Pipeline Input Validation Delta Specification

## MODIFIED Requirements

### Requirement: Pipeline executes registered capabilities through shared contracts

The system MUST validate schema-shaped Capability inputs before invoking registered handlers so all Interaction Surfaces receive consistent input failure behavior.

#### Scenario: Generic run returns invalid input envelope

- GIVEN `exito run <capability-id>` is invoked for a Capability with a required input schema field
- WHEN the supplied complete input object omits that field
- THEN stdout contains a JSON Envelope with `ok: false`
- AND `error.code` is `INVALID_INPUT`
- AND `meta.capabilityId` is the requested Capability ID
