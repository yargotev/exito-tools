# Design

`catalog.intelligent-search-products` lives in the existing Catalog domain because it is product-search behavior. It remains independent from Cobra/Bubble Tea and is exposed through the shared capability execution pipeline.

## Domain model

Add an Intelligent Search-specific input/result next to the existing Legacy Catalog search types. The new capability keeps these provider semantics separate from `catalog.search-products`:

- `tradePolicy` is required and becomes the leading path facet `trade-policy/<id>`.
- `facets` are repeated `key=value` path segments and are URL-escaped per segment.
- `text`, typed `by` + repeated `value`, and raw `query` are mutually exclusive query builders.
- Typed lookup modes map to Intelligent Search query expressions: `product.id`, `sku.id`, `sku.ean`, `sku.reference`, `product.link`, and broad `id`.
- `page` and `count` use Intelligent Search pagination rather than Legacy `_from`/`_to`.

## HTTP client

Add a Catalog-domain HTTP searcher for Intelligent Search. It uses shared HTTP request metadata and the configured brand base URL. It calls:

```http
GET /api/io/_v/api/intelligent-search/product_search/{facets}
```

The client accepts optional cookie strings for caller-supplied `vtex_segment` / `vtex_session` diagnostics but must never expose cookie values in output or logs. If diagnostics include cookie presence, they should be boolean/redacted only.

## Configuration

Extend the Configuration Resolver with `vtexIntelligentSearch.<brand>.baseUrl` per profile, plus env/dotenv override keys for Exito and Carulla QA/prod profiles. Base URLs are non-sensitive. Cookies/session/segment tokens are secrets and are not stored in YAML.

## CLI

Add a nested command group:

```bash
exito catalog intelligent-search products [flags]
```

The CLI maps flags into `catalog.intelligent-search-products` input and writes the standard JSON envelope. User-facing help/errors remain English.
