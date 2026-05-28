# VTEX Intelligent Search research for future CLI implementation

_Last researched: 2026-05-28_

## Purpose

This document captures the current understanding of VTEX Intelligent Search variants so a future session can design and implement an Exito Tools CLI capability without redoing discovery. It intentionally focuses on **read-only search and diagnostics**.

## Executive summary

VTEX has two product-search families that should not be mixed:

- **Legacy Search API**: existing Exito Tools `catalog.search-products` capability uses `/api/catalog_system/pub/products/search` and raw `fq`, `ft`, `_from`, `_to`, `O` parameters.
- **Intelligent Search API**: storefront/search-engine API under `/api/io/_v/api/intelligent-search` with path-based facets, query-string text/ID lookup, page/count pagination, relevance-by-default sorting, autocomplete, correction, banners, facets, and optional VTEX Ads parameters.

For a future CLI, model Intelligent Search as a **new Catalog read-only capability**, not as a breaking change to `catalog.search-products`. Recommended ID: `catalog.intelligent-search-products` or `catalog.search-products-intelligent` (choose via OpenSpec before implementation). Keep Legacy and Intelligent Search separate because their filters, pagination, sorting, response shape, and regionalization behavior differ.

## Primary endpoints

Base servers from the OpenAPI schema:

- `https://{accountName}.{environment}.com.br/api/io/_v/api/intelligent-search`, with `environment=vtexcommercestable` for the stable VTEX environment.
- `https://{accountName}.myvtex.com/api/io/_v/api/intelligent-search`.
- `https://{storeDomain}/api/io/_v/api/intelligent-search` is documented, but Exito Tools should not use custom storefront domains as Intelligent Search defaults because CDN/storefront routing can differ from account/environment API hosts.
- Pickup-point availability uses `https://{accountName}.vtexcommercestable.com.br/api/intelligent-search/v0`.

See `docs/research/vtex-intelligent-search-api-reference.md` for the local summarized API reference.

| Endpoint | Purpose | CLI candidate |
| --- | --- | --- |
| `GET /product_search/{facets}` | Product listing for text, ID, categories, brand, specs, collections, price, trade policy. | `exito catalog intelligent-search products` |
| `GET /facets/{facets}` | Possible facets for a query/facet context. | `exito catalog intelligent-search facets` |
| `GET /autocomplete_suggestions` | Suggested terms and attributes for autocomplete. | Later: `suggest` |
| `GET /search_suggestions` | Suggested similar search terms. | Later: `suggest-terms` |
| `GET /correction_search` | Misspelling correction. | Later: `correct` |
| `GET /top_searches` | Top 10 searches in last 14 days. | Later: `top-searches` |
| `GET /banners/{facets}` | Banners registered for query/facet context. | Later: `banners` |
| `GET /pickup-point-availability/{facets}` | Delivery Promise pickup-point availability. | Later only if needed; requires extra geo/hash semantics. |

## Product search request model

### Facets path

`/product_search/{facets}` requires path-based facets, not `fq` query parameters.

Format:

```text
/{facetKey1}/{facetValue1}/{facetKey2}/{facetValue2}/.../{facetKeyN}/{facetValueN}
```

Important rules:

- `trade-policy/{id}` is mandatory according to the OpenAPI parameter description. Treat it as required in the CLI unless runtime validation proves a store-specific exception.
- Multiple values of the **same** facet behave as OR/union.
- Different facet keys behave as AND/intersection.
- A facet key may repeat, for example `color/blue/color/red`.
- Negative filters use value prefix `not:`, for example `color/blue/size/not:42`.
- The order of terms in the path is documented as not relevant to the search, but preserve user order in diagnostics.
- URL-escape every facet key/value path segment independently. Do not concatenate raw values.

Common facet keys:

| Facet key | Meaning | Example |
| --- | --- | --- |
| `trade-policy` | Sales channel / trade policy. Required. | `trade-policy/1` |
| `category-1`, `category-2`, ... | Department/category/subcategory levels. Full parent path must be declared. | `category-1/supermercado/category-2/lacteos` |
| `brand` | Brand slug/identifier. | `brand/alpina` |
| `{specificationName}` | Filterable product/SKU specification slug/name. | `color/blue` |
| `productClusterIds` | Collection ID. | `productClusterIds/262` |
| `price` | Price range as `min:max`. | `price/100:500` |
| `region-id` | Search GraphQL exposes this selected facet for regionalization. Validate on REST before first-class support. | `region-id/<regionId>` |

### Query text and ID lookup

The query-string parameter is named `query`; `q` is an alias. It can contain natural text or typed ID lookup expressions.

Supported ID forms from current docs/OpenAPI:

| Lookup | Query format |
| --- | --- |
| Product ID | `query=product:<id>` or `query=product.id:<id>` |
| SKU ID | `query=sku:<id>` or `query=sku.id:<id>` |
| SKU reference | `query=sku.reference:<ref>` |
| EAN | `query=sku.ean:<ean>` |
| Product slug/link text | `query=product.link:<slug>` |
| Broad ID | `query=id:<id>` where VTEX may match ProductID, ProductRefID, SKUID, SKURefID, or EAN |
| Multiple IDs | `query=sku.id:1;2;3` or same type with semicolon-separated values |

Constraints:

- All IDs in a multiple-ID query should be the same type.
- Intelligent Search also supports partial matches for leading digits in the storefront search bar, but for API diagnostics the CLI should prefer explicit typed query strings to avoid ambiguity.
- Product/SKU/specification fields are searchable only if the store's Intelligent Search settings/indexing support them.

### Query parameters for product search

| Parameter | Type/default | Meaning / recommendation |
| --- | --- | --- |
| `query` / `q` | string | Text or typed ID lookup. CLI should expose `--text`, `--by/--value`, and raw `--query`. |
| `count` | number, default `24` | Products per page. Use explicit bounds in CLI docs; do not auto-page silently. |
| `page` | number, default `1` | Page number, unlike Legacy `_from/_to`. |
| `sort` | enum or omitted | Omitted/null/empty means relevance. Values: `price:desc`, `price:asc`, `orders:desc`, `name:desc`, `name:asc`, `release:desc`, `discount:desc`. |
| `locale` | BCP 47 string | Target language; account must be indexed for it. |
| `hideUnavailableItems` | boolean, docs recommend `true` though OpenAPI default is `false` | When true, filters to available products by indexed availability. For CLI, default should be explicit in docs/config rather than implicit. |
| `simulationBehavior` | `default`, `skip`, `only1P` | Price/promotion freshness vs speed tradeoff. `skip` is faster but less current. |
| `showSponsored` | boolean, default `false` | VTEX Ads only. Keep optional/raw until store need is confirmed. |
| `sponsoredCount` | string | VTEX Ads only. |
| `advertisementPlacement` | enum | VTEX Ads only: `top_search`, `middle_search`, `search_shelf`, `cart_shelf`, `plp_shelf`, `autocomplete`, `homepage`. |
| `repeatSponsoredProducts` | boolean | VTEX Ads only. |

### Other endpoint parameters

- `/top_searches`: `locale`.
- `/autocomplete_suggestions`, `/search_suggestions`, `/correction_search`: `query`/`q`, `locale`.
- `/facets/{facets}`: `facets`, `query`, `locale`, `hideUnavailableItems`.
- `/banners/{facets}`: `facets`, `query`, `locale`.
- `/pickup-point-availability/{facets}`: `facets`, `query`, `an`, plus either `coordinates` + `zip-code` + `country` or `deliveryZonesHash` + `pickupsHash` depending on strategy.

## Regionalization, trade policy, segment, and geo

### Trade policy / sales channel

- The Intelligent Search REST facets path requires `trade-policy/{tradePolicyId}`.
- Search GraphQL also has `salesChannel` arguments, and some integrations use `salesChannel` query string for the same purpose when supported. Do not assume the REST endpoint accepts every GraphQL argument.
- For CLI v1, make `--trade-policy` required and encode it as `trade-policy/<id>`.

### VTEX segment/session cookies

VTEX Sessions uses two cookies:

- `vtex_session`: identifies the individual browsing session.
- `vtex_segment`: stores commercial-condition cache-key information such as channel, price tables, region, campaigns, and UTMs.

A CLI outside the browser may need to reproduce cookies if the store's Intelligent Search relies on personalized/B2B pricing, price tables, campaigns, or regionalized context. Treat cookie support as **advanced diagnostics**:

- `--cookie 'vtex_segment=...'`
- `--cookie 'vtex_session=...'`
- or `--cookie-file` later if needed.

Never log cookies or put them in stdout diagnostics. They can imply customer/session/commercial context.

### Geo and region ID

Regionalization optimizes search results by seller availability in the shopper region. VTEX sessions can derive `checkout.regionId` from `public.country` plus either `public.postalCode` or `public.geoCoordinates`; Checkout session data declares inputs `public.regionId`, `country`, `postalCode`, `geoCoordinates` and outputs `checkout.regionId`.

Possible future flow:

1. Resolve sellers/region using Checkout public regions API or Session Manager.
2. Pass resulting context to Intelligent Search via session cookies and/or a selected facet like `region-id/<regionId>` if validated for REST.
3. Use `hideUnavailableItems=true` to narrow visible inventory.

Do **not** bake geo-to-region behavior into the first Intelligent Search product command unless tested against real Exito/Carulla stores. Make region support raw/advanced first.

## Search behavior and operational caveats

- Intelligent Search relevance is the default sort; fixed sort values such as price or name can change how multi-term results rank.
- The search engine applies autocorrect/fuzzy behavior according to search behavior settings.
- Searchable fields include product name, brand, ProductID, ProductRefID, SKUID, SKURefID, and EAN by default; specifications/categories require search settings.
- Search results are limited to 50 pages.
- Facet availability depends on store configuration; non-filterable specifications do not show as facets.
- `hideUnavailableItems` depends on indexed availability; `simulationBehavior` influences freshness/cost of pricing/promotions.

## Recommended Exito Tools capability shape

### Capability

```text
Capability ID: catalog.intelligent-search-products
Domain: catalog
Risk: read-only
Audience: agents, people
Visibility: CLI first; TUI later after real workflows are known
```

### CLI flags

```bash
# Text search
exito catalog intelligent-search products \
  --brand exito \
  --trade-policy 1 \
  --text "leche deslactosada" \
  --count 12 \
  --page 1 \
  --hide-unavailable

# SKU lookup
exito catalog intelligent-search products \
  --brand exito \
  --trade-policy 1 \
  --by sku-id \
  --value 12345

# Multiple SKU IDs
exito catalog intelligent-search products \
  --brand exito \
  --trade-policy 1 \
  --by sku-id \
  --value 12345 \
  --value 67890

# Raw facets for diagnostics
exito catalog intelligent-search products \
  --brand exito \
  --trade-policy 1 \
  --facet category-1=supermercado \
  --facet category-2=lacteos \
  --facet color=not:red \
  --sort price:asc

# Advanced cookie/segment diagnostics
exito catalog intelligent-search products \
  --brand exito \
  --trade-policy 1 \
  --text arroz \
  --cookie 'vtex_segment=...'
```

Recommended simple `--by` enum:

- `text` -> `query=<value>`
- `product-id` -> `query=product.id:<value>`
- `sku-id` -> `query=sku.id:<value>`
- `ean` -> `query=sku.ean:<value>`
- `sku-reference` -> `query=sku.reference:<value>`
- `slug` -> `query=product.link:<value>`
- `id` -> `query=id:<value>`

Recommended advanced flags:

- `--trade-policy <id>` required.
- `--facet key=value` repeated. The command should append `trade-policy=<id>` first, then repeated user facets.
- `--negative-facet key=value` could encode `key/not:<value>` if we want safer UX than asking users to type `not:`.
- `--query <raw>` for exact API query; mutually exclusive with `--by/--value` and `--text`.
- `--locale <bcp47>`.
- `--count`, `--page`, `--sort`.
- `--hide-unavailable` / `--include-unavailable` with a clear default.
- `--simulation-behavior default|skip|only1P`.
- VTEX Ads flags only if requested: `--show-sponsored`, `--sponsored-count`, `--advertisement-placement`, `--repeat-sponsored-products`.
- `--cookie` hidden/agent-oriented; redact in logs.

### Output mapping

Do not expose raw VTEX response as the stable top-level contract. Map to domain-owned results and keep raw payload under diagnostics, following existing Catalog pattern.

Suggested output fields:

```json
{
  "query": "sku.id:12345",
  "facets": [
    {"key":"trade-policy","value":"1"},
    {"key":"brand","value":"alpina"}
  ],
  "page": 1,
  "count": 24,
  "sort": "relevance",
  "products": [
    {
      "productId": "...",
      "productName": "...",
      "brand": "...",
      "linkText": "...",
      "categoryId": "...",
      "items": [
        {
          "skuId": "...",
          "name": "...",
          "ean": "...",
          "referenceId": "...",
          "sellerIds": ["1"]
        }
      ],
      "diagnostics": {"providerPayload": {}}
    }
  ],
  "diagnostics": {
    "requestPath": "/product_search/trade-policy/1",
    "requestQuery": {"query":"sku.id:12345"}
  }
}
```

Preserve response headers useful for pagination/cache if VTEX returns them; unlike Legacy Search, Intelligent Search primarily uses `page`/`count` and may expose totals in response metadata depending on provider shape.

## Implementation notes for future session

1. Create an OpenSpec change before implementation. This is a new provider/capability contract.
2. Reuse shared HTTP infrastructure and Catalog domain package boundaries.
3. Add non-sensitive config for Intelligent Search base URL per brand/profile. Keep any app keys/cookies/session tokens out of YAML.
4. Unit-test URL construction heavily: facet escaping, repeated facets, negative facets, typed multi-ID query, and mutual exclusion of text/query/by-value modes.
5. Add HTTP client tests using fixtures from real VTEX response snippets once available.
6. Validate `region-id` REST support against a real store before making it first-class.
7. Keep user-facing CLI help/errors in English.

## Sources

- VTEX Developers: Intelligent Search API overview — https://developers.vtex.com/docs/guides/intelligent-search-api-overview
- VTEX OpenAPI schema: `VTEX - Intelligent Search API.json` — https://raw.githubusercontent.com/vtex/openapi-schemas/master/VTEX%20-%20Intelligent%20Search%20API.json
- VTEX Developers: Consult product search information — https://developers.vtex.com/docs/guides/consult-product-search-information
- VTEX Help: Search behavior — https://help.vtex.com/en/docs/tutorials/search-behavior
- VTEX Developers: Search GraphQL schema — https://developers.vtex.com/docs/apps/vtex.search-graphql
- VTEX Developers: Sessions System Overview — https://developers.vtex.com/docs/guides/sessions-system-overview
- VTEX Developers: Enable the Region for SKUs — https://developers.vtex.com/docs/guides/enable-the-region-for-skus
- VTEX Developers: Session data available from VTEX apps — https://developers.vtex.com/docs/guides/session-data-available-from-vtex-apps
- VTEX Developers: Get sellers by region or address — https://developers.vtex.com/docs/guides/get-sellers-by-region-or-address
- VTEX Developers: FastStore regionalization — https://developers.vtex.com/docs/guides/faststore/storefront-features-implementing-regionalization-features

## Addendum: Exito/FastStore regionalization flow notes

_User-provided implementation notes captured on 2026-05-27 for future validation against the storefront repositories._

### Checkout Regions coverage API

The primary VTEX coverage validation API for known coordinates is Checkout Regions:

```http
GET https://{account}.vtexcommercestable.com.br/api/checkout/pub/regions?country={country}&sc={salesChannel}&geoCoordinates={longitude};{latitude}
```

Important details:

- VTEX expects coordinates as `longitude;latitude` in `geoCoordinates`, not latitude first.
- Current Exito Tools CLI coverage diagnostics consider coverage true when VTEX Checkout Regions returns at least one seller:

```ts
sellers.length > 0
```

- Historical Exito storefront logic considered coverage true only when VTEX returned at least one seller whose `id` was different from the account name:

```ts
sellers.some((seller) => seller.id !== account)
```

That older rule existed for storefront/product-price flows: for example, Intelligent Search may return a price for seller `exitocol`, but that seller does not identify the exact white-label store fulfilling the order. For the read-only CLI diagnostic, `hasCoverage` follows the broader VTEX Regions meaning and the raw returned sellers remain available for future business-specific interpretation.

Coordinate acquisition before calling Regions may come from:

- Sitidata / ServiInformación:
  - `POST https://sitidataws.sitimapa.co/api/multizonificador/geocoder/`
  - body: `{ city, address }`
- SmartQuick:
  - `GET https://exito.smartquick.com.co/restexito/servicio_rest/MantieneReceptor/consulta_direccion_zona_xml/{address}/{city}/EXITOCOM`

When coverage is true, the existing service updates Checkout orderForm shipping data:

```http
POST /api/checkout/pub/orderForm/{orderFormId}/attachments/shippingData
```

with `selectedAddresses[].geoCoordinates = [longitude, latitude]`.

For saved addresses that are poorly georeferenced, the service may patch Master Data AD:

```http
PATCH /api/dataentities/AD/documents/{id}
```

with:

```json
{
  "geoCoordinate": ["longitude", "latitude"],
  "validatedWithSmart": true
}
```

This coverage flow does not directly mutate `vtex_segment`; it updates orderForm shipping data and sometimes Master Data address records.

### FastStore/session regionalization flow

The storefront regionalization flow appears to be:

1. Session starts with default channel, for example `{"salesChannel":"1","regionId":""}`.
2. Address or pickup-point changes call `updateRegion(...)`.
3. `updateRegion` calls storefront GraphQL `fsGetRegion` with address/city/shipping type/sales channel/orderForm/auth cookie/country.
4. If coverage exists, the returned region ID is saved into FastStore session channel:

```json
{"salesChannel":"1","regionId":"REGION_ID"}
```

5. The application creates/updates VTEX segment through:

```http
POST /io/api/sessions
```

payload:

```json
{
  "public": {
    "regionId": { "value": "REGION_ID" },
    "sc": { "value": "1" }
  }
}
```

6. The response `segmentToken` is stored as cookie `vtex_segment`.
7. PLP/search/shelf/PDP product queries include selected facets like:

```json
{ "key": "channel", "value": "{\"salesChannel\":\"1\",\"regionId\":\"REGION_ID\"}" }
{ "key": "locale", "value": "..." }
```

8. PDP has an additional refresh path that reads `vtex_segment` and calls legacy catalog search by SKU:

```http
GET /api/catalog_system/pub/products/search/?fq=skuId:{skuId}
Cookie: vtex_segment=...
```

so `commertialOffer` is refreshed with regionalized seller/pricing data.

Security note: the referenced storefront flow reportedly also writes `vtex_segment` to localStorage. That should be reviewed before touching the flow because project guidance says not to store tokens in localStorage.

### GraphQL/API inventory related to regionalized product search

| Layer | Operation / endpoint | Role |
| --- | --- | --- |
| Storefront GraphQL | `/api/graphql` | Main front GraphQL gateway. |
| Region GraphQL | `fsGetRegion(...)` | Computes `regionId`, `hasCoverage`, selected addresses, postal code, shipping type. |
| VTEX private GraphQL | `https://{workspace}--{storeId}.myvtex.com/_v/private/graphql/v1` | Remote schema exposing `fsGetRegion`, `cities`, `dependenciesByCity`, `availableAddresses`, `fsUpdateDependency`. |
| Cities GraphQL | `cities(channel, shippingType)` | Lists cities by channel and shipping type. |
| Dependencies GraphQL | `dependenciesByCity(channel, shippingType, idCity)` | Lists pickup stores/dependencies with coordinates, postal code, address. |
| Addresses GraphQL | `availableAddresses(email)` | Lists saved addresses. |
| Session validation GraphQL | `validateSession(session, search)` | Validates FastStore session after channel update. |
| VTEX Sessions | `/io/api/sessions` | Creates segment token for `regionId` + sales channel. |
| Product Search GraphQL | `SearchQuery` | PLP/shelf search with `channel` and `locale` facets. |
| PDP GraphQL | `ProductPageQuery` | PDP product locator with `channel` and `locale` facets. |
| Legacy catalog SKU refresh | `/api/catalog_system/pub/products/search/?fq=skuId:{skuId}` | PDP regionalized offer refresh with `vtex_segment` cookie. |
| Session by coordinates | `/api/sessions`, `/api/segments/{segmentToken}` | Creates/reads segment by `country`, `geoCoordinates`, and `sc`. |
| Checkout simulation | `/api/checkout/pub/orderForms/simulation?sc={salesChannel}` | Discounts/cross-selling support, not primary region calculation. |

### Implication for Intelligent Search CLI

Adding GraphQL is viable, but it should be scoped carefully:

- **CLI v1 can call Intelligent Search REST directly** with explicit `--trade-policy`, `--query`/`--by`, and optional raw `--cookie vtex_segment=...` for users who already have a segment token.
- **A regionalized CLI mode likely needs a separate region/session preparation step**, either:
  - call Checkout Regions directly from coordinates and only report coverage/sellers, or
  - call the storefront/private GraphQL `fsGetRegion` flow, then create a VTEX session segment through `/io/api/sessions`, then pass `vtex_segment` to product search.
- Treat GraphQL region resolution as a **separate capability** or subcommand, not hidden behavior inside product search, because it may require orderForm/auth/country/shipping type/city/address and has side effects if it updates orderForm or Master Data in the existing service.
- A safe read-only first slice could be:
  1. `catalog.intelligent-search-products`: product search with explicit trade policy/facets/query and optional caller-provided segment cookie.
  2. Later `geo/region.resolve-vtex-region` or `catalog.resolve-region`: read-only coverage/region diagnostics from coordinates via Checkout Regions.
  3. Later `catalog.create-vtex-segment` only if explicitly approved, because creating/updating VTEX sessions is a write-like side effect even if operationally low risk.

## Roadmap for next implementation sessions

Use this sequence when the user says to start Intelligent Search implementation.

### Phase 1 — read-only Intelligent Search product CLI

Goal: implement the useful search CLI without hidden regionalization side effects.

Recommended OpenSpec/change scope:

- Add a new read-only Catalog capability, preferably `catalog.intelligent-search-products`.
- Add explicit CLI command, preferably:

```bash
exito catalog intelligent-search products [flags]
```

Must include:

- Brand/profile-based Intelligent Search base URL configuration.
- Required `--trade-policy <id>` mapped to `trade-policy/<id>` facet.
- Lookup modes:
  - `--text <term>` -> `query=<term>`.
  - `--by product-id|sku-id|ean|sku-reference|slug|id --value <value>` -> typed `query` expression.
  - repeated `--value` for same-type multi-ID lookup with semicolon join.
  - raw `--query <value>` for diagnostics.
- Repeated `--facet key=value` for additional path facets.
- Pagination/sort flags: `--page`, `--count`, `--sort`.
- Availability/simulation flags: `--hide-unavailable`, `--include-unavailable`, `--simulation-behavior default|skip|only1P`.
- Optional advanced cookie flag for caller-supplied segment context, with redaction in logs/output.
- Domain-owned output mapping with raw provider payload only under diagnostics.

Do not include automatic region/session creation in Phase 1. If a user has a segment, they can pass it explicitly.

### Phase 2 — read-only region and coverage diagnostics

Goal: connect the existing Geo direction with VTEX coverage while remaining read-only.

Recommended capability options:

- `geo.resolve-vtex-region` if the language is primarily geospatial/coverage.
- `catalog.resolve-vtex-region` if it is primarily search/catalog context.

Must support:

- Known coordinates input.
- Optional city/address input only if reusing existing geocoder behavior is approved.
- Checkout Regions call:

```http
GET /api/checkout/pub/regions?country={country}&sc={salesChannel}&geoCoordinates={longitude};{latitude}
```

- Preserve the `longitude;latitude` format exactly.
- Return sellers, region/coverage diagnostics, and `hasCoverage = sellers.length > 0`.
- Document the historical Exito storefront rule (`any seller.id != account`) as a product-price/white-label ambiguity rule, not the CLI diagnostic coverage rule.
- No orderForm shippingData writes and no Master Data AD patches in this read-only phase.

### Phase 3 — explicit VTEX segment/session preparation

Goal: create a reusable segment token for regionalized product queries.

This phase is side-effectful and should require an explicit OpenSpec risk decision.

Possible capability:

```text
catalog.create-vtex-segment
```

Behavior:

- Accept `regionId` and `salesChannel` explicitly, or consume the output of Phase 2.
- Call `/io/api/sessions` with:

```json
{
  "public": {
    "regionId": { "value": "REGION_ID" },
    "sc": { "value": "1" }
  }
}
```

- Return token metadata safely; do not log or persist token by default.
- Do not store `vtex_segment` in localStorage. If cookie export is needed, print a redacted/explicit command-safe form only when requested.

### Phase 4 — regionalized Intelligent Search convenience flow

Goal: combine the previous capabilities without hiding side effects.

Possible UX:

```bash
exito catalog intelligent-search products \
  --trade-policy 1 \
  --text arroz \
  --region-id REGION_ID \
  --create-segment-confirmed
```

or a workflow capability that explicitly orchestrates:

1. Resolve region/coverage.
2. Create segment token.
3. Run Intelligent Search with `vtex_segment` cookie.

This should not be Phase 1 because it combines read-only search with session mutation and possibly GraphQL/session dependencies.

### Phase 5 — optional GraphQL parity

Goal: mirror storefront behavior only if REST + segment is insufficient.

Investigate/implement only after Phase 1/2 validation:

- Storefront `/api/graphql` `SearchQuery`/`ProductPageQuery` with `channel` and `locale` selected facets.
- Private GraphQL `fsGetRegion` if Checkout Regions is insufficient.
- Cities/dependencies/address GraphQL if CLI must guide store/pickup selection.

Keep this as separate capability/workflow work because it introduces storefront-specific contracts beyond VTEX Intelligent Search REST.

## Live regionalization validation notes

The Phase 4 workflow was validated against `prod`/exitocol with the confirmation-required command:

```bash
exito --profile prod catalog intelligent-search regionalized-products \
  --brand exito \
  --country COL \
  --trade-policy 1 \
  --longitude <longitude> \
  --latitude <latitude> \
  --text <query> \
  --confirm
```

Known coordinates for repeatable tests:

| Place | Address | Latitude | Longitude | Prod region observed |
| --- | --- | ---: | ---: | --- |
| Poblado | Cra. 42 #10-58, El Poblado, Medellín, Antioquia | `6.2107605` | `-75.569325` | `U1cjZXhpdG9jb2wwMzM=` |
| Bello | Av. 32 #55-137, Hermosa Provincia, Bello, Antioquia | `6.34160679` | `-75.54018773` | `U1cjZXhpdG9jb2wwMzA=` |

Validation findings:

- `agua` confirmed that Poblado and Bello resolve to different prod regions, but common water SKUs kept the same selling prices in the sampled results. Ranking and top-result composition changed by region.
- `tomate` confirmed regional price differences. SKU `506130` (`Tomate Chonto Insuperable FRESCAMPO 1000 gr`) returned `$3.800` in Poblado and `$4.240` in Bello. SKU `1601910` (`Tomate Milano 1 und`) had no Poblado price in the sampled SKU lookup and returned `$2.580` in Bello.
- Some tomato SKUs kept the same selling price while list price changed by region: `124214809`, `124233565`, and `639051` had higher list prices in Bello.
- SKU-specific lookup with `--by sku-id --value <sku>` is the preferred validation method when confirming regional price changes, because text search ranking can differ by region and obscure like-for-like comparisons.
- Staging/QA is useful for integration behavior, but Bello returned `hasCoverage=false` with fallback region `U1cj` and anomalous QA prices in the sampled run. Use prod/exitocol for this Poblado-vs-Bello price comparison unless QA coverage data is corrected.
