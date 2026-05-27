# CLI Root Specification Delta

## ADDED Requirements

### Requirement: Intelligent Search products command

The CLI Surface MUST expose `catalog.intelligent-search-products` through `exito catalog intelligent-search products`.

#### Scenario: CLI text search

- **Given** the CLI is available
- **When** a caller runs `exito catalog intelligent-search products --brand exito --trade-policy 1 --text leche`
- **Then** the command MUST execute `catalog.intelligent-search-products`
- **And** stdout MUST contain the standard JSON envelope.

#### Scenario: CLI typed SKU lookup

- **Given** the CLI is available
- **When** a caller runs `exito catalog intelligent-search products --brand exito --trade-policy 1 --by sku-id --value 912350`
- **Then** the command MUST execute `catalog.intelligent-search-products` with typed lookup input.

#### Scenario: CLI rejects missing trade policy

- **Given** the CLI is available
- **When** a caller runs `exito catalog intelligent-search products --text leche`
- **Then** the command MUST fail before provider execution because `--trade-policy` is required.
