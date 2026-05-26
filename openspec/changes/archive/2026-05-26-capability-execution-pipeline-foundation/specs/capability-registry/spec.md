# Capability Registry Delta Specification

## ADDED Requirements

### Requirement: Registry supports executable capability lookup

The system MUST allow surfaces and execution code to look up immutable registered Capability entries by stable Capability ID.

#### Scenario: Registered executable capability is found

- GIVEN application wiring registered an executable Capability before finalization
- WHEN the finalized registry is queried by Capability ID
- THEN the matching immutable entry is returned.

#### Scenario: Duplicate capability IDs are rejected

- GIVEN a Capability ID is already registered during boot
- WHEN wiring attempts to register another Capability with the same ID
- THEN the registry rejects the duplicate with a stable failure outcome.
