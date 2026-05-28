# Proposal: Add regionalized Intelligent Search flow

## Intent

Add Phase 4 of Intelligent Search: a confirmation-gated convenience flow that resolves VTEX regional coverage, creates a transient VTEX segment, and runs Intelligent Search with that segment without exposing secrets.

## Scope

### In Scope
- Add workflow capability `catalog.regionalized-intelligent-search-products`.
- Add CLI command `exito catalog intelligent-search regionalized-products`.
- Reuse existing Geo region resolution, VTEX segment creation, and Intelligent Search product search ports.
- Keep segment cookies internal and redacted from stdout/diagnostics.

### Out of Scope
- GraphQL storefront parity, `fsGetRegion`, PDP GraphQL, or private storefront contracts.
- OrderForm shippingData writes, Master Data patches, browser/session/localStorage mutation.
- Persisting `vtex_segment` tokens.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `catalog`: Add the regionalized Intelligent Search workflow capability.
- `geo`: Expose resolved VTEX region IDs as stable diagnostic output.
- `cli-root`: Add CLI routing for regionalized Intelligent Search products.

## Approach

Create a workflow package that composes domain ports explicitly during app wiring. The workflow validates search inputs, resolves Checkout Regions from coordinates, chooses the first resolved region ID, creates a VTEX segment with explicit confirmation at the capability level, then searches products with an internal `vtex_segment` cookie.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/workflow/` | New | Regionalized Intelligent Search orchestration. |
| `internal/domain/geo/` | Modified | Stable region ID extraction in result. |
| `internal/app/app.go` | Modified | Explicit workflow capability registration. |
| `internal/surface/cli/root.go` | Modified | New CLI subcommand. |
| `docs/capabilities/` | New | Capability contract documentation. |
| `openspec/specs/` | Modified | Source-of-truth requirements after archive. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Hidden side effect from segment creation | Medium | Capability is safe-write and requires `--confirm`. |
| Token leakage | Medium | Keep cookie internal; output only token metadata/redacted diagnostics. |
| Wrong region choice | Low | Report selected and candidate regions; fail when none returned. |

## Rollback Plan

Remove the new workflow package, app registration, CLI subcommand, docs, and spec deltas. Existing Phase 1–3 capabilities remain unchanged.

## Dependencies

- Existing `geo.resolve-vtex-region`, `catalog.create-vtex-segment`, and `catalog.intelligent-search-products` domain ports.

## Success Criteria

- [ ] CLI requires `--confirm` before creating a segment.
- [ ] Workflow resolves region, creates segment, and searches with internal cookie.
- [ ] Output does not include the raw `vtex_segment` token.
- [ ] `make test` passes.
