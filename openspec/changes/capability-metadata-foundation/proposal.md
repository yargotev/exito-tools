# Proposal: Capability Metadata Foundation

## Summary

Extend neutral Capability definitions with the first stable metadata needed by discovery, execution surfaces, and future TUI filtering.

## Scope

- Add version and domain metadata to Capability definitions.
- Add risk and confirmation metadata.
- Add audience and visibility metadata for CLI/TUI/palette exposure decisions.
- Preserve defensive-copy behavior in the immutable registry now that definitions contain slices.
- Verify `exito capabilities` exposes the metadata in JSON inventory output.

## Non-goals

- Add full input/output JSON schema metadata.
- Implement confirmation enforcement for risky commands.
- Add real Orders or Geo domain capabilities.
