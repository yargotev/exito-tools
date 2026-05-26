## ADDED Requirements

### Requirement: Successful capability execution propagates pagination metadata

The Capability execution Pipeline MUST include pagination metadata returned by a successful handler in Envelope metadata without changing the success data shape.

#### Scenario: Successful handler returns pagination metadata

- GIVEN a registered Capability handler returns successful list data and pagination metadata
- WHEN the execution Pipeline runs the Capability
- THEN the returned envelope has `ok: true`
- AND `meta.pagination` contains the pagination metadata
- AND the successful data result remains under `data`

#### Scenario: Handler failure does not expose result pagination

- GIVEN a registered Capability handler returns an error
- WHEN the execution Pipeline runs the Capability
- THEN the returned envelope has `ok: false`
- AND no successful result pagination metadata is exposed in metadata
