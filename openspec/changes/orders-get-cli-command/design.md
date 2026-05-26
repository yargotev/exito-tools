# Design: Orders Get CLI Command

## Approach

The CLI Surface adds an `orders` command group with a single `get` subcommand. The explicit command builds the same complete input object used by generic `exito run`, then executes `orders.get` via the shared execution Pipeline.

This keeps typed CLI convenience at the surface boundary while preserving the Capability as the source of truth for validation, metadata, and handler execution.

## Error behavior

The `--id` flag is required by Cobra before execution. When the flag is present, the command emits the standard JSON Envelope from the Pipeline. Until a real Orders client exists, valid calls return `ORDERS_NOT_CONFIGURED` in the envelope.
