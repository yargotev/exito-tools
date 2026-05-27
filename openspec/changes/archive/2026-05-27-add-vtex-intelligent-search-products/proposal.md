# Add VTEX Intelligent Search product search

## Why

Exito Tools already exposes Legacy VTEX Catalog search through `catalog.search-products`, but VTEX Intelligent Search is a separate storefront/search-engine API with different facet, query, pagination, sorting, and regionalization semantics. Operators need a read-only CLI slice to inspect Intelligent Search product results directly without changing the stable Legacy Catalog capability.

## Scope

- Add a new read-only Catalog capability `catalog.intelligent-search-products`.
- Add an explicit CLI command `exito catalog intelligent-search products`.
- Resolve non-sensitive Intelligent Search base URLs per profile and brand from YAML/env/dotenv.
- Query VTEX Intelligent Search REST `GET /api/io/_v/api/intelligent-search/product_search/{facets}`.
- Require `tradePolicy` / `--trade-policy` and encode it as the first path facet.
- Support text search, typed lookup modes, raw query diagnostics, repeated path facets, page/count/sort, availability, simulation behavior, and optional caller-provided cookie diagnostics.
- Map provider responses into domain-owned product/SKU summaries while preserving raw provider payload only under diagnostics/details.

## Out of Scope

- Changing or replacing `catalog.search-products`.
- Automatic VTEX region resolution, Checkout orderForm mutation, Master Data writes, or hidden VTEX session/segment creation.
- GraphQL SearchQuery/ProductPageQuery parity.
- VTEX Ads first-class flags unless a later change requests them.
- Persisting, logging, or printing cookie values.
