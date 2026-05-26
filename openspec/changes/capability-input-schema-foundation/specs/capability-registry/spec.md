# Capability Registry Input Schema Delta Specification

## MODIFIED Requirements

### Requirement: Registry becomes immutable after boot

The system MUST treat the Capability Registry as immutable after finalization, including slice-backed metadata inside registered Capability definitions and their input schemas.

#### Scenario: Returned input schema fields cannot mutate registry state

- GIVEN a finalized registry contains a Capability with input schema fields
- WHEN a caller mutates the fields returned from `All`
- THEN a later `All` call returns the original input schema fields
