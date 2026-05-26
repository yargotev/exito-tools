# Change: capabilities-inventory-command

## Why

Agents need a machine-readable way to inspect the currently registered Capability Registry before real domain commands are added. The project already has a boot-time registry; this slice exposes that registry through a narrow `exito capabilities` CLI command without implementing `run`, Orders, Geo, or the execution pipeline.

## What Changes

- Add a JSON presenter foundation for CLI envelopes.
- Add `exito capabilities` as a machine-readable command that boots the application and emits the finalized registry inventory.
- Keep bare `exito` human-readable help and avoid advertising deferred `orders`, `geo`, or `run` commands.

## Out of Scope

- Real Orders or Geo capability implementations.
- Generic `exito run` execution.
- Request ID/correlation/duration execution pipeline.
- TUI capability discovery.
