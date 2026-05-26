# Design: Pagination metadata foundation

## Approach

Introduce `capability.PaginationMeta` as a small, surface-neutral metadata type for cursor-based list results. `capability.ExecutionResult` gains an optional `Pagination` pointer so handlers can attach pagination metadata only when the result represents a paged list. `capability.EnvelopeMeta` gains optional `pagination` so JSON command output can carry pagination metadata without changing successful data payloads.

The execution Pipeline copies handler pagination metadata into success envelope metadata. Failure envelopes remain focused on structured errors and do not propagate handler result metadata, because handlers that return an error do not have a successful result contract.

## Decisions

- `nextCursor` is a string and remains opaque to Exito Tools consumers.
- `hasMore` is a bool inside `pagination`; omitting the whole pagination object represents non-paginated results.
- Pipeline propagation copies the pagination struct value to avoid exposing mutable handler-owned pointer state.

## Risks

- Future providers may require additional paging fields. Those can be compatible additions inside `pagination` when a concrete list Capability needs them.
