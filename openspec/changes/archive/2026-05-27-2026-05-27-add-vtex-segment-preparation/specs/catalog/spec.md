# Catalog Specification Delta

## ADDED Requirements

### Requirement: VTEX segment preparation capability

Exito Tools MUST expose a confirmation-required `catalog.create-vtex-segment` capability for explicitly preparing a VTEX segment from a caller-provided region ID and sales channel.

#### Scenario: Create segment from region and sales channel

- **Given** VTEX public sessions is configured for the requested brand
- **When** a caller executes `catalog.create-vtex-segment` with `regionId` set to `REGION_ID` and `salesChannel` set to `1`
- **Then** the capability MUST call VTEX Sessions, preferring `POST /io/api/sessions` and falling back to `POST /api/sessions` when storefront routing requires it
- **And** the request body MUST include `public.regionId.value` equal to `REGION_ID`
- **And** the request body MUST include `public.sc.value` equal to `1`
- **And** the result MUST include safe token metadata that indicates whether a segment token was returned.

#### Scenario: Segment creation requires confirmation

- **Given** `catalog.create-vtex-segment` mutates VTEX session state by creating a segment token
- **When** a caller executes it without explicit confirmation
- **Then** the shared Pipeline MUST return a structured `CONFIRMATION_REQUIRED` failure
- **And** the provider MUST NOT be called.

#### Scenario: Token diagnostics are redacted

- **Given** VTEX Sessions returns a segment token
- **When** `catalog.create-vtex-segment` succeeds or fails after receiving provider data
- **Then** diagnostics MUST NOT include the raw token value
- **And** provider payload token fields MUST be redacted.

#### Scenario: Optional cookie output is explicit

- **Given** VTEX Sessions returns a segment token
- **When** the caller sets `includeCookie` to `true`
- **Then** the result MAY include a `cookie` field formatted as `vtex_segment=<token>`
- **And** the capability MUST otherwise omit unredacted cookie output by default.

#### Scenario: CLI command emits JSON envelope

- **Given** the CLI is available
- **When** a caller runs `exito catalog create-vtex-segment --brand exito --region-id REGION_ID --sales-channel 1 --confirm`
- **Then** the command MUST execute `catalog.create-vtex-segment`
- **And** the command MUST emit the standard JSON envelope on stdout.
