# VTEX Intelligent Search API local reference

_Last verified: 2026-05-28_

This is a local, summarized reference for Exito Tools work. It is intentionally not a verbatim mirror of VTEX documentation. When changing provider contracts, verify against the official sources listed at the end.

## Server selection

VTEX Intelligent Search is exposed through VTEX account/environment hosts. Exito Tools should prefer these hosts for provider diagnostics instead of custom storefront domains, because storefront domains may add CDN, routing, or response-code behavior that differs from provider/API hosts.

Primary server pattern for product search, facets, suggestions, banners, and correction:

```text
https://{accountName}.{environment}.com.br/api/io/_v/api/intelligent-search
```

For VTEX production/stable environments, `environment` is usually `vtexcommercestable`.

MyVTEX is also documented:

```text
https://{accountName}.myvtex.com/api/io/_v/api/intelligent-search
```

The official schema also lists a custom store-domain server, but Exito Tools should not use the custom domain as its default Intelligent Search base URL unless a test explicitly requires storefront parity.

Delivery Promise pickup-point availability uses a different base path:

```text
https://{accountName}.vtexcommercestable.com.br/api/intelligent-search/v0
```

### Exito Tools defaults

| Brand/profile | Account/environment base URL |
| --- | --- |
| Exito QA | `https://exito.vtexcommercestable.com.br` |
| Exito prod/exitocol | `https://exitocol.vtexcommercestable.com.br` |
| Carulla QA | `https://carulla.vtexcommercestable.com.br` |
| Carulla prod | `https://carulla.vtexcommercestable.com.br` |

Environment override variables:

```env
EXITO_VTEX_INTELLIGENT_SEARCH_BASE_URL_QA=https://exito.vtexcommercestable.com.br
EXITO_VTEX_INTELLIGENT_SEARCH_BASE_URL_PROD=https://exitocol.vtexcommercestable.com.br
CARULLA_VTEX_INTELLIGENT_SEARCH_BASE_URL_QA=https://carulla.vtexcommercestable.com.br
CARULLA_VTEX_INTELLIGENT_SEARCH_BASE_URL_PROD=https://carulla.vtexcommercestable.com.br
```

These endpoints are public storefront/search-engine endpoints and do not require VTEX app key/app token for the current read-only product-search use case.

## Endpoint inventory

| Endpoint | Method | Purpose | Exito Tools status |
| --- | --- | --- | --- |
| `/product_search/{facets}` | `GET` | Product listings for a query and facet context. | Implemented by `catalog.intelligent-search-products`. |
| `/facets/{facets}` | `GET` | Available facets for a query/facet context. | Future capability. |
| `/banners/{facets}` | `GET` | Banners configured for a query/facet context. | Future capability. |
| `/correction_search` | `GET` | Misspelling correction attempt. | Future capability. |
| `/autocomplete_suggestions` | `GET` | Autocomplete terms and attribute suggestions. | Future capability. |
| `/search_suggestions` | `GET` | Similar search-term suggestions. | Future capability. |
| `/top_searches` | `GET` | Frequently searched terms. | Future capability. |
| `/pickup-point-availability/{facets}` | `GET` | Delivery Promise pickup-point availability. | Future capability; different base path. |

## Product search contract

Exito Tools currently calls:

```text
GET /api/io/_v/api/intelligent-search/product_search/{facets}
```

### Facets path

`{facets}` is encoded as alternating path segments:

```text
/{facetKey1}/{facetValue1}/{facetKey2}/{facetValue2}/.../{facetKeyN}/{facetValueN}
```

Rules to preserve in implementation:

- `trade-policy/{id}` is required for Exito Tools calls and is inserted first from `--trade-policy`.
- Repeating the same facet key represents multiple accepted values for that facet.
- Different facet keys narrow the result set together.
- Negative filters use `not:` as the value prefix, for example `color/not:red`.
- Escape each key and value as an individual URL path segment.

Common facet keys:

| Facet key | Meaning | Example |
| --- | --- | --- |
| `trade-policy` | Sales channel/trade policy. | `trade-policy/1` |
| `category-1`, `category-2`, ... | Category hierarchy levels. | `category-1/supermercado/category-2/despensa` |
| `brand` | Brand facet. | `brand/alpina` |
| `productClusterIds` | Collection/cluster ID. | `productClusterIds/123` |
| `price` | Price range. | `price/10000:50000` |
| Store-specific specification facets | Filterable product/SKU specs. | `color/azul` |

### Query modes

The query parameter is `query`; VTEX also documents `q` as an alias. Exito Tools should continue using `query` for clarity.

Supported command modes:

| CLI mode | API query value |
| --- | --- |
| `--text arroz` | `query=arroz` |
| `--query <raw>` | Raw caller-provided value. |
| `--by product-id --value 123` | `query=product.id:123` |
| `--by sku-id --value 123` | `query=sku.id:123` |
| `--by ean --value 770...` | `query=sku.ean:770...` |
| `--by sku-reference --value REF` | `query=sku.reference:REF` |
| `--by slug --value product-slug` | `query=product.link:product-slug` |
| `--by id --value 123` | `query=id:123` |

For multiple values of the same typed lookup, Exito Tools joins them with semicolons, e.g. `sku.id:123;456`.

### Product search query parameters

| Parameter | Type/default | Notes |
| --- | --- | --- |
| `query` | string | Text or typed lookup expression. |
| `count` | number, default 24 in VTEX docs | Page size. Exito Tools sends an explicit `--count`. |
| `page` | number, default 1 | Intelligent Search page number, not Legacy Catalog `_from/_to`. |
| `sort` | string or empty | Empty/omitted means relevance. Documented values include price, orders, name, release, and discount directions. |
| `locale` | BCP 47 string | Only works if the account is indexed for that language. |
| `hideUnavailableItems` | boolean, default false in schema | `true` returns only products with stock according to indexed availability. |
| `simulationBehavior` | `default`, `skip`, `only1P` | Controls simulation behavior. Use only when intentionally testing freshness/performance behavior. |
| `showSponsored` | boolean | VTEX Ads only; not first-class in Exito Tools yet. |
| `sponsoredCount` | string | VTEX Ads only. |
| `advertisementPlacement` | enum | VTEX Ads placement such as search or shelf positions. |
| `repeatSponsoredProducts` | boolean | VTEX Ads only. |

## Other endpoint parameter notes

- `/facets/{facets}`: accepts facets, `query`, `locale`, and `hideUnavailableItems`.
- `/banners/{facets}`: accepts facets, `query`, and `locale`.
- `/correction_search`, `/autocomplete_suggestions`, `/search_suggestions`: accept `query` and optional `locale`.
- `/top_searches`: accepts optional `locale`.
- `/pickup-point-availability/{facets}`: requires account name (`an`) and uses either coordinate/ZIP/country inputs or precomputed delivery/pickup hashes depending on strategy.

## Regionalization notes

For Exito Tools regionalized product search, the current safe-write workflow is:

1. Resolve a VTEX Checkout Region from country and coordinates.
2. Create a transient VTEX segment for the resolved region and trade policy.
3. Call Intelligent Search with the generated `vtex_segment` cookie internally.

The segment token must stay redacted in logs and JSON diagnostics unless a user explicitly requests a copy/paste cookie from the segment command.

## Live verification examples

QA Exito product search:

```bash
curl 'https://exito.vtexcommercestable.com.br/api/io/_v/api/intelligent-search/product_search/trade-policy/1?query=arroz&page=1&count=1'
```

Prod/exitocol product search:

```bash
curl 'https://exitocol.vtexcommercestable.com.br/api/io/_v/api/intelligent-search/product_search/trade-policy/1?query=arroz&page=1&count=1'
```

Exito Tools QA command:

```bash
./exito --profile staging catalog intelligent-search products \
  --brand exito \
  --trade-policy 1 \
  --text arroz \
  --count 1
```

## Sources

- VTEX Developers: Intelligent Search API reference — https://developers.vtex.com/docs/api-reference/intelligent-search-api
- VTEX OpenAPI schema: `VTEX - Intelligent Search API.json` — https://raw.githubusercontent.com/vtex/openapi-schemas/master/VTEX%20-%20Intelligent%20Search%20API.json
- VTEX Search onboarding guide — https://developers.vtex.com/docs/guides/search-overview
