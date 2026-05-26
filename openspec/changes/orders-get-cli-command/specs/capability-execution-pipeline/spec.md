# Capability Execution Pipeline Orders CLI Delta Specification

## MODIFIED Requirements

### Requirement: Pipeline executes registered capabilities through shared contracts

The system MUST allow explicit domain CLI commands to execute registered Capabilities through the same shared Pipeline used by generic execution.

#### Scenario: Explicit Orders command preserves envelope metadata

- GIVEN a correlation ID and profile are supplied to the CLI
- WHEN `exito orders get --id A123` executes `orders.get`
- THEN the JSON Envelope metadata includes the selected profile, correlation ID, and `orders.get` Capability ID
