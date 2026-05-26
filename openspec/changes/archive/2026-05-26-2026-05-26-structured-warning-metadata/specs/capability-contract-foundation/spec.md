## ADDED Requirements

### Requirement: Envelope metadata supports structured warnings

The system MUST model non-fatal warnings as structured metadata with stable codes, messages, and optional details.

#### Scenario: Structured warning serializes in metadata

- GIVEN a JSON Envelope includes one Structured Warning in metadata
- WHEN the envelope is serialized
- THEN `meta.warnings[0].code` contains the stable warning code
- AND `meta.warnings[0].message` contains the warning message
- AND optional warning details are serialized under `meta.warnings[0].details`
