# Proposal: Add Master Data Read Tools

## Intent

Add a read-only Master Data Domain so agents can inspect VTEX Master Data documents, schemas, and indices through stable Exito Tools capabilities without mutating provider state.

## Scope

### In Scope
- Resolve `vtexMasterData` configuration by Effective Profile and brand (`exito`, `carulla`).
- Add `internal/domain/masterdata` with domain-owned models, VTEX client mapping, and unavailable-provider behavior.
- Register read-only capabilities: `masterdata.get-document`, `masterdata.search-documents`, `masterdata.scroll-documents`, `masterdata.list-schemas`, `masterdata.get-schema`, and `masterdata.list-indices`.
- Preserve pagination metadata, VTEX limits, structured warnings, and payload/log safety.

### Out of Scope
- Document create/update/patch/delete.
- Schema or index mutations.
- Entity list/get capabilities.
- Auto-draining search/scroll results or parallel scroll orchestration.
- New production dependencies.

## Capabilities

### New Capabilities
- `masterdata`: Master Data provider configuration, read-only document inspection, search/scroll pagination, and v2 schema/index discovery.

### Modified Capabilities
- `configuration-resolver`: Add profile/brand `vtexMasterData` base URL and VTEX app key/token resolution.

## Approach

Follow existing VTEX OMS/Checkout brand-client patterns: config resolves brand providers, `internal/app` wires configured or unavailable clients, domain use cases validate inputs and map VTEX DTOs into safe domain results, and capabilities execute through the existing generic `exito run` pipeline.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/config` | Modified | Add `VTEXMasterDataProvider` and YAML/env resolution. |
| `internal/domain/masterdata` | New | Read-only use cases, HTTP client, brand router, tests. |
| `internal/app` | Modified | Explicitly register Master Data capabilities. |
| `docs/configuration.md` | Modified | Document non-sensitive YAML and env vars. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| PII exposure from raw documents | Medium | Require/encourage explicit fields, no raw payload logging, safe diagnostics only. |
| Unbounded reads | Medium | Enforce search max 100, scroll max 1000, no auto-pagination. |
| Credential leakage | Low | Reuse secret-omitting config structs and app key/token headers only. |

## Rollback Plan

Remove `internal/domain/masterdata`, config additions, app registrations, docs, and this OpenSpec change. No external state rollback is needed because the slice is read-only.

## Dependencies

- Existing Configuration Resolver, capability registry, execution pipeline, and HTTP infrastructure.
- VTEX app key/token credentials with Master Data permissions.

## Success Criteria

- [ ] `make test` passes.
- [ ] Capability inventory includes the six read-only Master Data IDs.
- [ ] Fake-server tests verify VTEX paths, auth headers, limits, pagination, and no write requests.
