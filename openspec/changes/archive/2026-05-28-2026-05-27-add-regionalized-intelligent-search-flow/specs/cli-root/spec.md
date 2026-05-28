## ADDED Requirements

### Requirement: Regionalized Intelligent Search products command

The CLI Surface MUST expose `catalog.regionalized-intelligent-search-products` through `exito catalog intelligent-search regionalized-products`.

#### Scenario: CLI regionalized text search

- **Given** the CLI is available
- **When** a caller runs `exito catalog intelligent-search regionalized-products --brand exito --country COL --trade-policy 1 --longitude -74.160580822 --latitude 4.598090587 --text arroz --confirm`
- **Then** the command MUST execute `catalog.regionalized-intelligent-search-products`
- **And** stdout MUST contain the standard JSON envelope.

#### Scenario: CLI requires confirmation

- **Given** the CLI is available
- **When** a caller runs `exito catalog intelligent-search regionalized-products --trade-policy 1 --longitude -74 --latitude 4 --text arroz` without `--confirm`
- **Then** stdout MUST contain a JSON Envelope with `ok: false` and error code `CONFIRMATION_REQUIRED`.
