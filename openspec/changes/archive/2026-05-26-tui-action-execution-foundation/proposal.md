# Proposal: TUI action execution foundation

## Summary

Add the narrow first TUI task-runner foundation so a people-facing Command Palette Action can be selected and executed through the shared Capability execution Pipeline, with visible loading, success, and structured failure state.

## Motivation

The previous TUI Command Palette slice made Actions discoverable but intentionally deferred execution. The next smallest useful TUI slice is to prove that TUI Actions invoke the same registered Capabilities and shared Pipeline as the CLI rather than duplicating domain logic.

## Scope

- Add Command Palette selection movement.
- Execute the selected palette Action with an empty input object through the shared Pipeline.
- Render loading, success, and structured failure state in the TUI model.
- Keep forms/input collection, cancellation controls, result filters, and rich result views out of scope.

## Non-goals

- Building per-capability TUI forms.
- Executing primary navigation Actions directly.
- Provider-specific rendering.
- Full cancellation UX.
