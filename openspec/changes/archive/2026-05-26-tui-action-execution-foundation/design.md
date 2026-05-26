# Design: TUI action execution foundation

## Approach

Extend the existing side-effect-free TUI model with enough execution state to run a selected Command Palette Action. The model keeps the finalized Registry and effective Profile from the shared Application. Pressing enter while the Command Palette is open resolves the currently selected filtered Action and returns a Bubble Tea command that executes the Capability through `execution.Pipeline` using an empty input object.

The command returns a model message containing the standard Envelope. The model renders a simple task status line for loading, success, or failure, including the Capability ID and structured error code/message when available.

## Decisions

- Selection is only introduced inside the Command Palette for this slice.
- `up`/`down` move the selected palette row.
- Filtering clamps selection back to a valid row.
- Execution uses the shared Pipeline and Registry; no TUI-specific use-case path is introduced.
- Empty input is intentional for the foundation; schema-aware forms are deferred.

## Risks

- Capabilities with required input will fail with `INVALID_INPUT` until TUI forms are added.
- There is still no cancellation key or long-running progress view; this slice only establishes the execution seam and terminal states.
