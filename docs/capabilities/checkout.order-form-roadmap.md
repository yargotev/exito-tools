# checkout orderForm roadmap

Guided VTEX Checkout purchase assembly is the next Exito Tools roadmap area. Checkout capabilities own orderForm state and are separate from Catalog product discovery and Orders lookup.

## Roadmap capabilities

### checkout.get-order-form

- **Risk**: read-only
- **Audience**: agents, people
- **Goal**: Retrieve a known VTEX orderForm by ID and return a redacted summary suitable for diagnostics.
- **CLI shape**: `exito checkout get-order-form --brand <brand> --order-form-id <id>`

### checkout.create-order-form

- **Risk**: safe-write
- **Requires Confirmation**: yes
- **Audience**: agents, people
- **Goal**: Create a fresh VTEX orderForm/cart for a brand and sales channel.
- **CLI shape**: `exito checkout create-order-form --brand <brand> --sales-channel <sc> --confirm`

### checkout.add-items

- **Risk**: safe-write
- **Requires Confirmation**: yes
- **Audience**: agents, people
- **Goal**: Add selected SKU items to an existing orderForm after product discovery.
- **CLI shape**: `exito checkout add-items --brand <brand> --order-form-id <id> --item sku=<sku>,quantity=<qty>[,seller=<seller>] --confirm`

### checkout.update-client-profile

- **Risk**: safe-write
- **Requires Confirmation**: yes
- **Audience**: agents, people
- **Goal**: Attach customer profile data to an existing orderForm.
- **CLI shape**: `exito checkout update-client-profile --brand <brand> --order-form-id <id> --input-json <profile-json> --confirm`

### checkout.update-shipping-data

- **Risk**: safe-write
- **Requires Confirmation**: yes
- **Audience**: agents, people
- **Goal**: Attach delivery address data and selected logistics options to an existing orderForm.
- **CLI shape**: `exito checkout update-shipping-data --brand <brand> --order-form-id <id> --input-json <shipping-json> --confirm`

## Guided flow

1. Search products with `catalog.intelligent-search-products` or `catalog.regionalized-intelligent-search-products`.
2. Create or load an orderForm with Checkout.
3. Add selected SKU IDs and quantities.
4. Attach or update client profile data.
5. Attach or update shipping data and logistics selections.
6. Inspect the prepared orderForm summary.

## Safety rules

- Checkout writes are confirmation-gated safe-write capabilities.
- A guided purchase assembly flow must execute Checkout writes sequentially; no parallel mutations for the same orderForm.
- PII, cookies, and provider tokens must be redacted from stdout JSON, logs, diagnostics, and committed docs by default.
- First slice stops before place-order, process-order, and payment endpoints.
