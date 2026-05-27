# Proposal: Add VTEX catalog product search

## Motivation
Operators and agents need to inspect VTEX storefront product catalog data by SKU ID, product ID, EAN, reference ID, seller, category, brand, collection, text, or raw VTEX filters without leaving Exito Tools.

## Scope
- Add a read-only `catalog.search-products` capability backed by VTEX Legacy Search API product search.
- Add `exito catalog search-products` with simple `--by/--value` lookup and advanced `--fq`, `--ft`, `--order`, `--from`, and `--to` options.
- Configure public VTEX catalog base URLs per brand/profile.
- Preserve provider payload details for diagnostics while exposing a stable summary result.

## Out of Scope
- Catalog mutation, indexing, pricing updates, or inventory changes.
- Authenticated private Catalog API endpoints.
- Archiving the change in this slice.
