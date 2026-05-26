## ADDED Requirements

### Requirement: Successful capability execution propagates warnings

The Capability execution Pipeline MUST include non-fatal warnings returned by a successful handler in Envelope metadata without changing the success state.

#### Scenario: Successful handler returns warning metadata

- GIVEN a registered Capability handler returns successful data and one Structured Warning
- WHEN the execution Pipeline runs the Capability
- THEN the returned envelope has `ok: true`
- AND `meta.warnings` contains the warning
- AND the warning does not change the successful data result

#### Scenario: Handler failure does not expose result warnings

- GIVEN a registered Capability handler returns an error
- WHEN the execution Pipeline runs the Capability
- THEN the returned envelope has `ok: false`
- AND no successful result warnings are exposed in metadata
