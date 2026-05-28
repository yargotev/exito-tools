# Design: Regionalized Intelligent Search flow

## Context

Phase 1 provides read-only Intelligent Search with optional caller-provided cookies. Phase 2 resolves VTEX Checkout Regions from coordinates. Phase 3 creates a VTEX segment token from an explicit region ID and sales channel. Phase 4 should combine these steps while keeping the session creation side effect explicit.

## Decisions

- Add a workflow capability ID `catalog.regionalized-intelligent-search-products` because the user-facing outcome is still Catalog product search, even though it composes Geo and Catalog ports.
- Mark the capability `safe-write` and `requiresConfirmation: true` because it creates a VTEX session segment.
- Implement the orchestration under `internal/workflow` so domain packages remain independent of each other and of Cobra.
- Extend `geo.ResolveVTEXRegionResult` with `regions[]` containing IDs and sellers. This is an additive output change required for orchestration.
- Do not expose raw cookies in the workflow result. The segment creator is invoked with `IncludeCookie: true` internally, but the result only records `tokenSet`/`tokenLength` and sanitized nested results.

## Flow

1. Validate brand, country, coordinates, trade policy, and query mode.
2. Resolve VTEX regions using `geo.VTEXRegionResolver` with `salesChannel = tradePolicy`.
3. Fail with `REGIONALIZED_SEARCH_NO_REGION` when no region ID is returned.
4. Create a VTEX segment with the selected region ID and trade policy.
5. Fail with `REGIONALIZED_SEARCH_SEGMENT_UNAVAILABLE` when no token/cookie is returned.
6. Execute Intelligent Search with the generated `vtex_segment` cookie.
7. Return selected region, region diagnostics, segment token metadata, and product search result with cookie names only.

## Non-goals

- No GraphQL parity.
- No orderForm or Master Data mutation.
- No token persistence or user-visible raw token.
