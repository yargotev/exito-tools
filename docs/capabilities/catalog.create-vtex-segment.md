# catalog.create-vtex-segment

Creates a VTEX session segment token from an explicit VTEX region ID and sales channel.

- **Capability ID**: `catalog.create-vtex-segment`
- **Domain**: `catalog`
- **CLI Command**: `exito catalog create-vtex-segment [flags]`
- **Risk**: safe-write
- **Requires Confirmation**: yes
- **Audience**: agents

## Examples

```bash
exito catalog create-vtex-segment \
  --brand exito \
  --region-id REGION_ID \
  --sales-channel 1 \
  --confirm

exito catalog create-vtex-segment \
  --brand exito \
  --region-id REGION_ID \
  --sales-channel 1 \
  --include-cookie \
  --confirm
```

## Contract notes

- The command calls VTEX Sessions, preferring `POST /io/api/sessions` and falling back to `POST /api/sessions` when storefront routing requires it, with `public.regionId.value` and `public.sc.value`.
- This capability is confirmation-required because it creates VTEX session/segment state.
- The raw segment token is not included in diagnostics and is not persisted by Exito Tools.
- `--include-cookie` explicitly returns `vtex_segment=<token>` for copy/paste into later Intelligent Search diagnostics.
