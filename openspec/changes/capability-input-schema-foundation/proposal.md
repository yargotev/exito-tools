# Proposal: Capability Input Schema Foundation

## Summary

Extend neutral Capability definitions with minimal schema-first input metadata so discovery surfaces can expose complete input object contracts before real Orders or Geo capabilities are implemented.

## Scope

- Add neutral input schema and input field contracts to the Capability core.
- Add stable input type constants for common JSON-shaped values.
- Include input schemas in `exito capabilities` inventory output when definitions provide them.
- Preserve registry defensive-copy behavior for slice-backed input schema fields.

## Non-goals

- Validate `exito run` input objects against schemas.
- Generate explicit domain command flags from schemas.
- Add real Orders or Geo domain capabilities.
- Model full JSON Schema, nested field constraints, enums, or output schemas.
