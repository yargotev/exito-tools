# Proposal: Capability Input Validation Foundation

## Summary

Use neutral Capability input schemas in the shared execution pipeline to reject invalid complete input objects before invoking Capability handlers.

## Scope

- Add basic schema validation for required fields.
- Add basic JSON-shaped type validation for string, number, boolean, object, and array inputs.
- Return standard failed envelopes with stable `INVALID_INPUT` structured errors.
- Verify generic `exito run` receives invalid-input envelopes through the shared pipeline.

## Non-goals

- Implement full JSON Schema validation.
- Validate unknown/extra input fields.
- Add nested object field validation, enums, min/max constraints, or custom validators.
- Generate explicit domain command flags from input schemas.
- Add real Orders or Geo domain capabilities.
