# Change: Add explicit VTEX segment preparation

## Motivation

Regionalized VTEX Intelligent Search diagnostics need a safe way to prepare a reusable `vtex_segment` from an explicit `regionId` and sales channel. Phase 1 already supports caller-supplied segment cookies, and Phase 2 resolves VTEX region coverage. Phase 3 bridges those capabilities while making the session side effect visible and confirmation-gated.

## Scope

- Add a new Catalog capability `catalog.create-vtex-segment`.
- Add CLI command `exito catalog create-vtex-segment` with explicit `--confirm`.
- POST to VTEX Sessions `/io/api/sessions` with `public.regionId` and `public.sc`.
- Return safe token metadata and optionally a command-safe cookie string only when requested.
- Redact segment token values from diagnostics by default and never persist them.

## Out of Scope

- Automatic region resolution inside product search.
- OrderForm shipping data writes or Master Data patches.
- Local storage, browser session mutation, or persisted cookie files.
- GraphQL storefront parity.
