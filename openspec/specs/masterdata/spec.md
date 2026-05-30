# Master Data Specification

## Purpose

Define read-only VTEX Master Data capabilities for document inspection, bounded extraction, and v2 schema/index discovery.

## Requirements

### Requirement: Master Data domain owns read-only operations

Exito Tools MUST expose Master Data behavior through `internal/domain/masterdata` and read-only `masterdata.*` capabilities backed by domain-owned models.

#### Scenario: Capabilities are discoverable

- GIVEN the application boots
- WHEN capability inventory is inspected
- THEN `masterdata.get-document`, `masterdata.search-documents`, `masterdata.scroll-documents`, `masterdata.list-schemas`, `masterdata.get-schema`, and `masterdata.list-indices` MUST be registered
- AND each MUST have read-only risk and no confirmation requirement

#### Scenario: Provider unavailable fails structurally

- GIVEN the requested brand is not configured
- WHEN a Master Data capability executes
- THEN it MUST fail with stable code `MASTERDATA_NOT_CONFIGURED`

### Requirement: Document reads expose safe domain results

The system MUST get, search, and scroll documents without leaking credentials or relying on provider DTOs as public contracts.

#### Scenario: Get document by ID

- GIVEN Master Data is configured for brand `exito`
- WHEN `masterdata.get-document` executes with `entity`, `documentId`, optional `schema`, and optional `fields`
- THEN VTEX `/api/dataentities/{entity}/documents/{id}` MUST be called
- AND the result MUST include brand, entity, document ID, selected fields, document data, and safe diagnostics

#### Scenario: Search enforces VTEX bounded pagination

- GIVEN `masterdata.search-documents` receives `rangeFrom` and `rangeTo`
- WHEN the requested range exceeds 100 documents
- THEN execution MUST fail with `MASTERDATA_INVALID_INPUT`
- AND VTEX MUST NOT be called

#### Scenario: Search warns without sort

- GIVEN a paginated search has no `sort`
- WHEN `masterdata.search-documents` succeeds
- THEN metadata MUST include a structured warning about unstable pagination

#### Scenario: Scroll exposes token metadata

- GIVEN `masterdata.scroll-documents` receives optional `token` and `size`
- WHEN VTEX returns `X-VTEX-MD-TOKEN`
- THEN pagination metadata MUST expose that token as the next cursor
- AND size greater than 1000 MUST fail with `MASTERDATA_INVALID_INPUT`

### Requirement: Schema and index reads expose v2 control-plane metadata

The system MUST list/get schemas and list indices without performing schema or index mutations.

#### Scenario: List and get schemas

- GIVEN Master Data is configured
- WHEN schema list or get capabilities execute
- THEN VTEX v2 schema endpoints MUST be called
- AND returned schema JSON MUST be exposed as read-only domain data

#### Scenario: List indices

- GIVEN Master Data is configured
- WHEN `masterdata.list-indices` executes with `brand` and `entity`
- THEN VTEX `/api/dataentities/{entity}/indices` MUST be called
- AND no write, delete, or schema lifecycle endpoint MUST be called

### Requirement: Master Data outputs protect sensitive data

The system MUST keep stdout JSON valid while avoiding raw credential/header logging and unbounded PII exposure.

#### Scenario: Diagnostics are safe

- GIVEN any Master Data provider call completes
- WHEN the result or error is rendered
- THEN diagnostics MUST omit app keys, app tokens, request headers, and raw cookies
