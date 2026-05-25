# Capability Registry Metadata Delta Specification

## MODIFIED Requirements

### Requirement: Registry becomes immutable after boot

The system MUST treat the Capability Registry as immutable after finalization, including slice-backed metadata inside registered Capability definitions.

#### Scenario: Post-boot mutation is rejected

- GIVEN the registry has been finalized
- WHEN code attempts another registration
- THEN the attempt is rejected with a stable failure outcome

#### Scenario: Returned definition metadata cannot mutate registry state

- GIVEN a finalized registry contains a Capability with audience or visibility metadata
- WHEN a caller mutates the slices returned from `All`
- THEN a later `All` call returns the original metadata
