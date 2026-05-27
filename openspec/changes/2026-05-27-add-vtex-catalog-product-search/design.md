# Design: VTEX catalog product search

`catalog.search-products` lives in a new Catalog operational domain because storefront product lookup is catalog behavior, separate from Orders and Geo.

The domain accepts two lookup modes:

1. Simple mode: `by` + `value`, mapping stable friendly names to VTEX Search API query parameters.
2. Advanced mode: raw repeated `fq` filters plus optional `ft`, `order`, `from`, and `to`.

Simple mode supported keys: `sku-id`, `product-id`, `ref-id`, `ean`, `seller-id`, `category`, `brand-id`, `collection-id`, `text`, and `slug`.

The HTTP adapter calls the public VTEX Legacy Search API endpoint `/api/catalog_system/pub/products/search`. It treats HTTP 200 and 206 as successful responses, captures the VTEX `resources` header for range/total metadata, maps selected product/SKU fields into domain-owned structures, and preserves each product payload in `details`.

Configuration adds public non-sensitive `vtexCatalog.<brand>.baseUrl` values. The provider requires only base URLs, no VTEX app key/token, because this search endpoint is public storefront search.
