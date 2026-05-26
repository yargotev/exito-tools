# Capability Contract Input Validation Delta Specification

## ADDED Requirements

### Requirement: Capability input schemas are enforceable at execution time

The system MUST be able to validate complete Capability input objects against neutral input schema metadata before a Capability handler runs.

#### Scenario: Required input field is missing

- GIVEN a Capability has an input schema with a required field
- WHEN execution receives an input object without that field
- THEN execution returns a structured `INVALID_INPUT` failure
- AND the Capability handler is not invoked

#### Scenario: Input field has the wrong type

- GIVEN a Capability has an input schema with a typed field
- WHEN execution receives an input object with an incompatible value type
- THEN execution returns a structured `INVALID_INPUT` failure
- AND the Capability handler is not invoked

#### Scenario: Capability has no input schema

- GIVEN a Capability does not define an input schema
- WHEN execution receives an input object
- THEN execution preserves existing behavior and invokes the Capability handler
