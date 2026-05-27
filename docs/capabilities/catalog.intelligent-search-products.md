# catalog.intelligent-search-products

Searches VTEX Intelligent Search products using the storefront/search-engine REST API.

- **Capability ID**: `catalog.intelligent-search-products`
- **Domain**: `catalog`
- **CLI Command**: `exito catalog intelligent-search products [flags]`
- **Risk**: read-only
- **Audience**: agents, people

## Examples

```bash
exito catalog intelligent-search products --brand exito --trade-policy 1 --text "leche deslactosada"
exito catalog intelligent-search products --brand exito --trade-policy 1 --by sku-id --value 912350
exito catalog intelligent-search products --brand exito --trade-policy 1 --by sku-id --value 123 --value 456
exito catalog intelligent-search products --brand exito --trade-policy 1 --facet category-1=supermercado --facet category-2=lacteos
```

## Contract notes

- `tradePolicy` / `--trade-policy` is required and encoded as the first path facet: `trade-policy/<id>`.
- Additional facets use repeated `key=value` path segments.
- Query modes are mutually exclusive: `text`, raw `query`, or typed `by` plus repeated `value`.
- Pagination uses Intelligent Search `page` and `count`, not Legacy Catalog `_from` and `_to`.
- Optional cookie inputs are sent to the provider but only cookie names may appear in diagnostics; values must stay redacted.
