# catalog.regionalized-intelligent-search-products

Resolves VTEX Checkout Regions from coordinates, creates a transient VTEX segment, and runs VTEX Intelligent Search products with that segment.

- **Capability ID**: `catalog.regionalized-intelligent-search-products`
- **Risk**: `safe-write`
- **Requires Confirmation**: yes
- **CLI Command**: `exito catalog intelligent-search regionalized-products [flags]`

## Examples

```bash
exito catalog intelligent-search regionalized-products \
  --brand exito \
  --country COL \
  --trade-policy 1 \
  --longitude -74.160580822 \
  --latitude 4.598090587 \
  --text arroz \
  --confirm
```

```bash
exito catalog intelligent-search regionalized-products \
  --brand exito \
  --trade-policy 1 \
  --longitude -74.160580822 \
  --latitude 4.598090587 \
  --by sku-id \
  --value 912350 \
  --confirm
```

## Behavior

1. Calls VTEX Checkout Regions with `geoCoordinates={longitude};{latitude}`.
2. Selects the first returned region ID.
3. Creates a VTEX segment for that region and trade policy.
4. Runs Intelligent Search with the generated `vtex_segment` cookie internally.

The raw segment token is not printed, logged, persisted, or exposed in diagnostics. Output includes only token metadata such as `tokenSet` and `tokenLength`.

## Non-goals

- No GraphQL storefront parity.
- No orderForm `shippingData` writes.
- No Master Data address patches.
- No browser, localStorage, or cookie jar mutation.
