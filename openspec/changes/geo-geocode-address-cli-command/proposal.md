# Proposal: Geo geocode-address CLI command

## Summary

Expose the registered `geo.geocode-address` Capability through an explicit domain CLI command: `exito geo geocode-address --city <city> --address <address>`.

## Scope

- Add `geo` CLI command group.
- Add `geo geocode-address` explicit command with required `--city` and `--address` flags.
- Route command execution through the shared Pipeline and standard JSON envelope.
- Propagate profile and correlation metadata.
- Update root help expectations and command tests.

## Out of Scope

- Real Geo provider HTTP client.
- Configuration/dotenv parsing for Geo credentials.
- TUI action implementation.
