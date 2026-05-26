# Proposal: TUI input form foundation

## Summary

Add the narrow first schema-aware TUI input form so selected Command Palette Actions with required string inputs can collect values before execution.

## Motivation

The TUI action execution foundation executes selected Actions through the shared Pipeline, but it always sends an empty input object. Capabilities such as `orders.get` and `geo.geocode-address` therefore fail with `INVALID_INPUT`. The next smallest useful slice is to collect simple required string fields from the Capability Input Schema and submit them through the existing TUI execution path.

## Scope

- Detect selected palette Actions with required string input fields.
- Render a minimal input form for those fields.
- Support typing, backspace, enter-to-advance, and submit after all fields are populated.
- Execute the Action through the shared Pipeline with collected input.
- Keep optional fields, non-string field editors, validation hints, masking, and rich form layout out of scope.

## Non-goals

- Full schema renderer for every input type.
- Editing optional fields.
- Primary action forms.
- Provider-specific result views.
