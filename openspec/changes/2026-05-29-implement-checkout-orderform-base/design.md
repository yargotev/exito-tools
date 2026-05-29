# Design: Checkout orderForm base

## Architecture

Checkout is implemented as a dedicated Operational Domain under `internal/domain/checkout`, matching ADR-0056. The domain package owns provider DTO mapping, use cases, validation, capability definitions, and domain result models. It does not import Cobra, Bubble Tea, or surface packages.

Application wiring in `internal/app` explicitly creates the Checkout brand client and registers capabilities at boot. The registry remains immutable after boot.

## Configuration

The Configuration Resolver exposes `VTEXCheckoutProvider` with `exito` and `carulla` brand entries. Base URLs are non-sensitive and may come from:

1. Process environment.
2. Profile-specific dotenv.
3. General dotenv.
4. YAML `profiles.<profile>.vtexCheckout.<brand>.baseUrl`.

Environment variable names follow the existing VTEX public endpoint pattern:

- `EXITO_VTEX_CHECKOUT_BASE_URL_QA`
- `EXITO_VTEX_CHECKOUT_BASE_URL_PROD`
- `CARULLA_VTEX_CHECKOUT_BASE_URL_QA`
- `CARULLA_VTEX_CHECKOUT_BASE_URL_PROD`

Cookies, ownership values, customer profile values, and credentials remain execution-time sensitive data and are not read from YAML or serialized in effective configuration.

## Capabilities

### checkout.get-order-form

- Domain: `checkout`
- Risk: `read-only`
- Confirmation: not required
- Visibility: CLI and command palette
- Audience: agents and people
- Inputs: `brand`, `orderFormId`

The use case loads a known VTEX orderForm and returns a domain-owned summary.

### checkout.create-order-form

- Domain: `checkout`
- Risk: `safe-write`
- Confirmation: required
- Visibility: CLI and command palette
- Audience: agents and people
- Inputs: `brand`, `salesChannel`

The shared Pipeline enforces confirmation before handler execution. The HTTP client calls the VTEX current cart endpoint with `forceNewCart=true` and `sc=<salesChannel>`.

## Provider mapping and redaction

The VTEX Checkout HTTP adapter maps provider orderForm DTOs into `OrderFormSummary` with:

- `brand`
- `id`
- `salesChannel`
- `value`
- totalizer summaries
- item summaries
- `itemCount`
- boolean presence flags for client profile and shipping data
- safe diagnostics containing request path and provider status

The summary intentionally does not expose raw client profile fields, shipping address values, cookies, ownership cookies, or provider credential values.

## CLI

The CLI remains a surface adapter. Commands build capability input objects and execute through the shared Pipeline:

```bash
exito checkout get-order-form --brand exito --order-form-id ORDER_FORM_ID
exito checkout create-order-form --brand exito --sales-channel 1 --confirm
```

Without `--confirm`, `create-order-form` emits a standard JSON failure envelope with `CONFIRMATION_REQUIRED` and the provider is not called.
