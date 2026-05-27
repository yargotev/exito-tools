# catalog.search-products

Searches VTEX public catalog products using the Legacy Search API.

- **Capability ID**: `catalog.search-products`
- **Domain**: `catalog`
- **CLI Command**: `exito catalog search-products [flags]`
- **Risk**: read-only
- **Audience**: agents, people

## Simple lookup

```bash
exito catalog search-products --by sku-id --value 912350 --brand exito
exito catalog search-products --by product-id --value 534690
exito catalog search-products --by ref-id --value 912350
exito catalog search-products --by ean --value 7706060050094
exito catalog search-products --by seller-id --value VMIABBA
exito catalog search-products --by category --value 34185087/34185141/34185508
exito catalog search-products --by brand-id --value 6
exito catalog search-products --by collection-id --value 172
exito catalog search-products --by text --value minibar
exito catalog search-products --by slug --value nevera-minibar-97-litros-abba-534690
```

## Advanced lookup

```bash
exito catalog search-products \
  --fq skuId:912350 \
  --fq sellerId:VMIABBA \
  --from 0 \
  --to 0 \
  --order OrderByPriceASC
```

Supported raw VTEX filters include category (`C:/.../`), brand (`B:{id}`), specification (`specificationFilter_{id}:{value}`), price range (`P:[0 TO 200000]`), collection (`productClusterIds:{id}`), product ID, SKU ID, reference ID, EAN13, sales-channel availability, and seller ID.

## Output

The result includes product and SKU summaries, parsed `resources` pagination metadata when VTEX returns it, and preserves the full provider payload for each product under `details`.
