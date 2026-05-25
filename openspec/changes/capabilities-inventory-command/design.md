# Design: capabilities-inventory-command

## Approach

Expose the existing immutable registry through a dedicated CLI subcommand. The CLI command remains a surface adapter: it bootstraps the application using parsed root flags, reads `Application.Registry.All()`, wraps the result in the shared envelope shape, and writes JSON to stdout through a small presenter package.

## Boundaries

- `internal/surface/cli` owns Cobra command wiring and command-specific output decisions.
- `internal/presenter` owns JSON encoding details.
- `internal/capability` owns envelope and capability contract shapes.
- Operational domains remain absent and surface-independent.

## Decisions

- `exito capabilities` returns an empty inventory (`[]`) successfully while no real capabilities are registered.
- The command uses the standard envelope top-level shape now (`ok`, `data`, `meta`) but keeps request IDs and duration metadata deferred to the execution pipeline slice.
- Root help may list `capabilities` because it is implemented; it still must not imply `orders`, `geo`, or `run` exist.
