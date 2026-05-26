# Capability Execution Pipeline Specification

## Requirements

### Requirement: Pipeline executes registered capabilities through shared contracts

The system MUST allow explicit domain CLI commands to execute registered Capabilities through the same shared Pipeline used by generic execution.

#### Scenario: Generic run returns invalid input envelope

- GIVEN `exito run <capability-id>` is invoked for a Capability with a required input schema field
- WHEN the supplied complete input object omits that field
- THEN stdout contains a JSON Envelope with `ok: false`
- AND `error.code` is `INVALID_INPUT`
- AND `meta.capabilityId` is the requested Capability ID

#### Scenario: Explicit Orders command preserves envelope metadata

- GIVEN a correlation ID and profile are supplied to the CLI
- WHEN `exito orders get --id A123` executes `orders.get`
- THEN the JSON Envelope metadata includes the selected profile, correlation ID, and `orders.get` Capability ID
