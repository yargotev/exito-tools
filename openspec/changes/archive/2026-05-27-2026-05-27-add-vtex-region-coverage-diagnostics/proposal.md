# Add VTEX region coverage diagnostics

## Why

Phase 1 Intelligent Search product lookup is complete. The next safe phase is read-only regional coverage diagnostics so operators can validate whether known coordinates resolve to VTEX seller coverage before later deciding on segment/session creation.

## What changes

- Add a read-only Geo Domain capability `geo.resolve-vtex-region`.
- Add a CLI command `exito geo resolve-vtex-region`.
- Resolve coverage from known longitude/latitude by calling VTEX Checkout Regions:
  `GET /api/checkout/pub/regions?country={country}&sc={salesChannel}&geoCoordinates={longitude};{latitude}`.
- Return sellers, region diagnostics, request diagnostics, and `hasCoverage` using the existing Exito rule: any returned seller ID different from the account/brand has coverage.

## Non-goals

- No orderForm `shippingData` writes.
- No Master Data address patches.
- No VTEX session or segment creation.
- No automatic city/address geocoding in this slice.
