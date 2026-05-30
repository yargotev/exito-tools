# VTEX Master Data tools research

## Purpose

Capture the research needed to add a Master Data operational domain to Exito Tools. The intended toolset should let agents and operators inspect Master Data entities, documents, schemas, indices, and later perform controlled document/schema mutations.

## Primary VTEX references

- [Master Data API overview](https://developers.vtex.com/docs/guides/master-data-api)
- [Master Data API v2 reference](https://developers.vtex.com/docs/api-reference/master-data-api-v2)
- [Working with JSON Schemas in Master Data v2](https://developers.vtex.com/docs/guides/working-with-json-schemas-in-master-data-v2)
- [Master Data schema lifecycle](https://developers.vtex.com/docs/guides/master-data-schema-lifecycle)
- [Extracting data with Search and Scroll](https://developers.vtex.com/docs/guides/extracting-data-from-master-data-with-search-and-scroll)
- [Pagination in the Master Data API](https://developers.vtex.com/docs/guides/pagination-in-the-master-data-api)

## Key findings

### Master Data API generations

VTEX Master Data has two relevant API generations:

- **Master Data v1** exposes legacy data entities and document APIs. Data entity configuration is primarily an Admin concern, but the API exposes data entities and documents.
- **Master Data v2** exposes document APIs plus a Control Plane for schemas and indices.

VTEX documents that v1 and v2 entities are not interchangeable. A v2 implementation should not assume it can operate on v1 entities such as `CL` and `AD` unless those entities are actually exposed through the selected API path/account behavior.

### Data plane

Core document operations use paths under:

```text
/api/dataentities/{entity}/documents
/api/dataentities/{entity}/documents/{id}
/api/dataentities/{entity}/search
/api/dataentities/{entity}/scroll
```

Useful document/query inputs:

- `entity`: Master Data data entity acronym/name.
- `documentId`: provider document identifier.
- `_fields`: comma-separated fields to return.
- `_where`: query predicate for search.
- `_schema`: optional schema selector for v2 validation/query behavior.
- `_sort`: recommended for paginated search consistency.
- `REST-Range`: header for search pagination.
- `_size` and `X-VTEX-MD-TOKEN`: scroll pagination state.

### Search versus scroll

Use **search** for targeted, bounded queries and operator inspection. Use **scroll** for extraction-style reads.

Important constraints from VTEX docs:

- Search has a maximum page size of 100 documents per query.
- Search is intended for up to roughly 10,000 documents in a query window.
- Search pagination should include `_sort` to avoid inconsistent page windows.
- Scroll can request up to 1000 documents per page.
- Scroll returns and then reuses `X-VTEX-MD-TOKEN`.
- Scroll tokens expire after 20 minutes.
- VTEX limits each account to 10 simultaneous scrolls.

Exito Tools should expose explicit pagination metadata and should not auto-drain all results by default.

### Control plane

Master Data v2 exposes schema and index operations under paths such as:

```text
/api/dataentities/{entity}/schemas
/api/dataentities/{entity}/schemas/{schemaName}
/api/dataentities/{entity}/indices
/api/dataentities/{entity}/indices/{indexName}
```

Schema documents are JSON Schema-based and may include VTEX-specific metadata such as:

- `v-indexed`
- `v-default-fields`
- `v-security`
- `v-cache`
- `v-triggers`

Schema lifecycle changes can trigger background processing and temporarily affect validation/query behavior. Schema and index mutations are therefore higher-risk than document reads.

### Authentication and configuration

Master Data operations should use server-side VTEX credentials:

```text
X-VTEX-API-AppKey: <app key>
X-VTEX-API-AppToken: <app token>
```

This matches the existing VTEX OMS pattern in `internal/domain/orders/vtex_oms_getter.go`. Secrets must remain in environment variables or non-committed dotenv files, never YAML/docs.

Recommended non-sensitive YAML shape:

```yaml
profiles:
  staging:
    vtexMasterData:
      exito:
        baseUrl: https://exito.vtexcommercestable.com.br
      carulla:
        baseUrl: https://carulla.vtexcommercestable.com.br
```

Recommended environment variables:

```env
EXITO_VTEX_MASTERDATA_BASE_URL_QA=https://exito.vtexcommercestable.com.br
EXITO_VTEX_MASTERDATA_BASE_URL_PROD=https://exitocol.vtexcommercestable.com.br
CARULLA_VTEX_MASTERDATA_BASE_URL_QA=https://carulla.vtexcommercestable.com.br
CARULLA_VTEX_MASTERDATA_BASE_URL_PROD=https://carulla.vtexcommercestable.com.br

# Reuse the existing VTEX server-side credential convention when the same app key/token has Master Data permissions.
EXITO_APP_KEY_QA=...
EXITO_APP_TOKEN_QA=...
EXITO_APP_KEY_PROD=...
EXITO_APP_TOKEN_PROD=...
CARULLA_APP_KEY_QA=...
CARULLA_APP_TOKEN_QA=...
CARULLA_APP_KEY=...
CARULLA_APP_TOKEN=...
```

If reuse of OMS credential variable names becomes confusing, add Master Data-specific credential aliases later, but keep one resolved provider model inside the Configuration Resolver.

## Domain fit for Exito Tools

Master Data should be its own operational domain:

```text
internal/domain/masterdata/
```

Rationale:

- Master Data has its own external API semantics, schemas, indices, pagination, and risks.
- It should not be hidden inside Orders, Catalog, Checkout, or Workflow.
- Domain clients should map VTEX DTOs into domain-owned result shapes before returning data to surfaces.
- CLI/TUI should execute the same capabilities through the shared execution pipeline.

Suggested package files:

```text
internal/domain/masterdata/masterdata.go
internal/domain/masterdata/http_client.go
internal/domain/masterdata/brand_client.go
internal/domain/masterdata/unavailable.go
internal/domain/masterdata/*_test.go
```

Suggested config additions:

- `config.VTEXMasterDataProvider`
- brand provider with base URL, app key/token, source metadata, and configured flags.
- YAML parser support for `profiles.<profile>.vtexMasterData.<brand>.baseUrl`.
- environment override support for Master Data base URLs by profile/brand.

## Data safety

Master Data can contain customer data and operationally sensitive records. Exito Tools should:

- Treat reads as `read-only`, but still avoid logging raw payloads.
- Let callers explicitly select `_fields`; avoid broad defaults that expose PII.
- Redact or summarize payloads in diagnostics and logs.
- Require confirmation for all document writes.
- Treat deletes and schema/index mutations as `destructive` or at least high-risk confirmation-gated operations.
- Keep stdout JSON envelopes clean and machine-readable.

## Open questions before implementation

- Which exact Master Data accounts/hosts should be used for Exito and Carulla in staging and prod?
- Are existing VTEX AppKey/AppToken credentials authorized for Master Data, or do we need separate app keys?
- Which v1 entities are business-critical for first usage? Examples may include `CL`, `AD`, or custom entities.
- Should default read outputs include raw document `fields`, or only selected/safe fields unless `--fields` is supplied?
- Should write capabilities support full JSON payloads from stdin/files only, avoiding many CLI flags?
