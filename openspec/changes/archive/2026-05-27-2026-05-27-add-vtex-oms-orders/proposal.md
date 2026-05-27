# Add independent VTEX OMS Orders lookup

## Why

GEOMS remains the current `orders.get` source, but VTEX OMS exposes additional order detail that is useful for operations. The lookup should be available independently so users can compare VTEX OMS with GEOMS without changing the stable GEOMS-backed contract.

## Scope

- Add an independent `orders.get-vtex` capability and `exito orders get-vtex` CLI command.
- Keep `orders.get` backed by GEOMS unchanged.
- Resolve VTEX OMS base URLs as non-sensitive configuration and app key/token credentials from environment/dotenv only.
- Support Exito and Carulla brand credential names across QA and production profiles.

## Out of Scope

- Browser/client exposure of VTEX credentials.
- Replacing or enriching `orders.get` with VTEX OMS data.
- VTEX Catalog, Masterdata, or Intelligent Search capabilities.
