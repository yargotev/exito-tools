## ADDED Requirements

### Requirement: Envelope metadata supports pagination

The system MUST model cursor-based pagination as optional structured Envelope metadata for list-style Capability results.

#### Scenario: Pagination metadata serializes in metadata

- GIVEN a JSON Envelope includes pagination metadata
- WHEN the envelope is serialized
- THEN `meta.pagination.nextCursor` contains the opaque cursor when one exists
- AND `meta.pagination.hasMore` declares whether another page is available
