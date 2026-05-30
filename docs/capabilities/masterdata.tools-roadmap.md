# Master Data tools roadmap

## Purpose

Define the proposed capability slices for adding VTEX Master Data tooling to Exito Tools. This roadmap is documentation only; implementation should start with an OpenSpec change in a future session.

## Domain placement

Master Data should be implemented as a new **Master Data Domain** rather than as part of Catalog, Orders, Checkout, or Workflow.

Proposed domain name and capability ID prefix:

```text
Domain: masterdata
Capability IDs: masterdata.<action>
```

The domain should follow existing Exito Tools rules:

- Domain package does not import Cobra, Bubble Tea, or `internal/surface/*`.
- Provider DTOs stay inside the domain API client.
- Use cases return domain-owned result models.
- Capabilities are wired explicitly in `internal/app`.
- Registry remains immutable after boot.
- CLI output stays in the standard JSON envelope.

## Slice 1: Read-only discovery and document inspection

Goal: provide useful Master Data exploration without mutating external state.

### Proposed capabilities

#### `masterdata.list-entities`

Lists known Master Data entities where the provider/API generation supports it.

Inputs:

- `brand`: `exito` or `carulla`; default `exito`.
- `apiVersion`: optional `v1` or `v2`; default should be explicit in the spec.

Risk: `read-only`.

Notes:

- v1 exposes data entity metadata more directly.
- v2 may not provide a reliable global list of entities; if unavailable, return a structured unsupported error or omit this from the first implementation.

#### `masterdata.get-entity`

Gets metadata for a single data entity/acronym when supported.

Inputs:

- `brand`
- `entity`
- `apiVersion`

Risk: `read-only`.

#### `masterdata.get-document`

Gets one Master Data document by ID.

Inputs:

- `brand`
- `entity`
- `documentId`
- `schema`: optional v2 schema selector.
- `fields`: optional field list.

Risk: `read-only`.

Output guidance:

- Include brand, entity, document ID, selected fields, and provider diagnostics.
- Avoid exposing secrets or request headers.
- Do not log raw document payloads.

#### `masterdata.search-documents`

Searches Master Data documents with bounded pagination.

Inputs:

- `brand`
- `entity`
- `fields`: field list; required or default to a safe minimal set.
- `where`: optional query predicate.
- `schema`: optional v2 schema selector.
- `sort`: recommended for stable pagination.
- `rangeFrom`: default 0.
- `rangeTo`: default no higher than 99.

Risk: `read-only`.

Contract requirements:

- Enforce VTEX search max page size of 100.
- Include pagination metadata derived from response headers when available.
- Warn when no sort is supplied for paginated requests.
- Do not auto-fetch multiple pages by default.

#### `masterdata.scroll-documents`

Starts or continues a Master Data scroll for large extraction-style reads.

Inputs:

- `brand`
- `entity`
- `fields`
- `where`: optional.
- `schema`: optional.
- `size`: max 1000.
- `token`: optional previous `X-VTEX-MD-TOKEN`.

Risk: `read-only`.

Contract requirements:

- Return the next scroll token in pagination metadata, not hidden state.
- Document token expiry behavior.
- Enforce size max 1000.
- Do not run concurrent scroll drains automatically.

#### `masterdata.list-schemas`

Lists v2 schemas for a data entity.

Inputs:

- `brand`
- `entity`

Risk: `read-only`.

#### `masterdata.get-schema`

Gets a v2 schema definition for a data entity.

Inputs:

- `brand`
- `entity`
- `schema`

Risk: `read-only`.

#### `masterdata.list-indices`

Lists v2 indices for a data entity.

Inputs:

- `brand`
- `entity`

Risk: `read-only`.

### Suggested first OpenSpec change

```text
openspec/changes/<date>-add-masterdata-read-tools/
```

Suggested requirement groups:

- Master Data Domain owns Master Data operations.
- Master Data provider configuration resolves by profile and brand.
- Read-only document inspection uses JSON envelopes and domain-owned results.
- Search and scroll expose explicit pagination and enforce VTEX limits.
- Schema and index reads expose v2 control-plane metadata without mutation.
- Raw payload/log safety prevents credential and PII leakage.

### Minimal implementation order

1. Add config resolver support for `vtexMasterData` base URLs and credentials.
2. Add `internal/domain/masterdata` with brand-aware client and unavailable client.
3. Implement `get-document` and `search-documents` first because they validate auth, URL, mapping, pagination, and output shape.
4. Add `scroll-documents` after search pagination behavior is stable.
5. Add schema/index reads.
6. Wire capabilities in `internal/app`.
7. Add CLI commands only after the generic `exito run <capability-id>` path is covered by tests, unless explicit commands are needed immediately.

## Slice 2: Document writes

Goal: safely create and update documents once read-only behavior is verified.

### Proposed capabilities

#### `masterdata.create-document`

Creates a document.

Inputs:

- `brand`
- `entity`
- `schema`: optional.
- `payload`: JSON object, preferably from `--input-json`, file, or stdin.

Risk: `safe-write`.
Requires confirmation: yes.

#### `masterdata.upsert-document`

Creates or replaces a known document.

Inputs:

- `brand`
- `entity`
- `documentId`
- `schema`: optional.
- `payload`: JSON object.

Risk: `safe-write`, potentially destructive if full replacement semantics are used.
Requires confirmation: yes.

#### `masterdata.patch-document`

Partially updates a document.

Inputs:

- `brand`
- `entity`
- `documentId`
- `schema`: optional.
- `payload`: JSON object.

Risk: `safe-write`.
Requires confirmation: yes.

#### `masterdata.delete-document`

Deletes a document.

Inputs:

- `brand`
- `entity`
- `documentId`

Risk: `destructive`.
Requires confirmation: yes.

### Guardrails

- Missing `--confirm` must return `CONFIRMATION_REQUIRED` and not call VTEX.
- CLI should not prompt interactively in automation paths.
- Output should summarize document IDs/status rather than echo full sensitive payloads.
- Tests should verify request body shape without real provider credentials.

## Slice 3: Schema and index mutations

Goal: expose v2 Control Plane changes only after document reads/writes are stable.

### Proposed capabilities

- `masterdata.save-schema`
- `masterdata.delete-schema`
- `masterdata.put-index`
- `masterdata.delete-index`

Risk:

- `safe-write` for save/update operations only if impact is bounded.
- `destructive` for delete operations.

Requires confirmation: yes.

### Guardrails

- Use JSON input files/stdin for schemas and indices.
- Include warnings about schema lifecycle/background processing.
- Preserve schema names and index names exactly.
- Avoid default schema mutation from TUI primary navigation; keep agent-oriented unless a guided review flow exists.
- Consider requiring profile/brand/entity/schema confirmation fields for deletes.

## Slice 4: Operator conveniences

Potential follow-up capabilities once core contracts are reliable:

- `masterdata.validate-document` — local JSON/schema validation when schema is known.
- `masterdata.diff-document` — compare local payload to provider document.
- `masterdata.export-documents` — controlled extraction with explicit max pages/records.
- `masterdata.inspect-entity` — aggregate entity, schemas, indices, and sample documents into one read-only diagnostic.

These should remain bounded and explicit to avoid accidental large exports.

## Recommended next-session prompt

Use this prompt to begin implementation planning:

```text
Inicia el SDD proposal/spec/design/tasks para el primer slice `add-masterdata-read-tools` en exito-tools. Usa la documentación en `docs/research/vtex-master-data-tools.md` y `docs/capabilities/masterdata.tools-roadmap.md`. El alcance debe ser solo read-only: configuración vtexMasterData por profile/brand, dominio `internal/domain/masterdata`, y capabilities `masterdata.get-document`, `masterdata.search-documents`, `masterdata.scroll-documents`, `masterdata.list-schemas`, `masterdata.get-schema`, y `masterdata.list-indices`. No implementes writes todavía.
```
