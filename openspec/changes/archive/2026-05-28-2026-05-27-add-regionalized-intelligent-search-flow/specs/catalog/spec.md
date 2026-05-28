## ADDED Requirements

### Requirement: Regionalized Intelligent Search product workflow

Exito Tools MUST expose a confirmation-required `catalog.regionalized-intelligent-search-products` workflow capability that resolves a VTEX region from coordinates, creates a VTEX segment, and runs VTEX Intelligent Search with that segment.

#### Scenario: Resolve, segment, and search

- **Given** VTEX Checkout Regions, VTEX Sessions, and VTEX Intelligent Search are configured for the requested brand
- **When** a caller executes `catalog.regionalized-intelligent-search-products` with `country`, `longitude`, `latitude`, `tradePolicy`, and a valid Intelligent Search query mode
- **Then** Exito Tools MUST resolve VTEX regions using Checkout Regions and `geoCoordinates={longitude};{latitude}`
- **And** it MUST create a VTEX segment from the selected region ID and trade policy
- **And** it MUST execute Intelligent Search with an internal `vtex_segment` cookie.

#### Scenario: Workflow requires confirmation

- **Given** the workflow creates a VTEX session segment
- **When** a caller executes it without explicit confirmation
- **Then** the shared Pipeline MUST return `CONFIRMATION_REQUIRED`
- **And** no provider MUST be called.

#### Scenario: Segment token is not exposed

- **Given** VTEX Sessions returns a segment token
- **When** the workflow succeeds or fails after segment creation
- **Then** stdout JSON and diagnostics MUST NOT include the raw `vtex_segment` token
- **And** the result MAY include only token metadata such as whether a token was set and token length.

#### Scenario: No resolved region fails safely

- **Given** Checkout Regions returns no region ID for the requested coordinates
- **When** the workflow executes
- **Then** it MUST fail with stable structured error code `REGIONALIZED_SEARCH_NO_REGION`
- **And** it MUST NOT create a VTEX segment.
