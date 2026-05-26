# Execution Metadata Specification

## Requirements

### Requirement: Request metadata is generated for JSON command output

The system MUST create opaque request metadata for machine-readable command results without depending on any Interaction Surface framework.

#### Scenario: Request ID is generated

- GIVEN a JSON command is being rendered
- WHEN execution metadata is created
- THEN a non-empty opaque `requestId` is available for the JSON Envelope metadata

#### Scenario: Duration is measured in milliseconds

- GIVEN a command has a start and finish time
- WHEN execution metadata is converted to Envelope metadata
- THEN `durationMs` is present as a non-negative integer
