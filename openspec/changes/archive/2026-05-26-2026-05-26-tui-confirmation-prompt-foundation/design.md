# Design: TUI confirmation prompt foundation

## Context

Capability definitions already expose `Risk` and `RequiresConfirmation`, and the shared execution Pipeline enforces confirmation by returning `CONFIRMATION_REQUIRED` when `ExecuteRequest.Confirmed` is false. The TUI model currently starts execution from the command palette or input form with no confirmation step.

## Approach

- Add a `confirmationState` to the TUI model with the selected Capability definition and collected input.
- Route command-palette selection and submitted forms through a small `confirmOrStartExecution` helper:
  - if the Capability does not require confirmation, start execution normally with `Confirmed: false`,
  - if confirmation is required, close palette/form state and render the confirmation prompt instead of executing.
- Add prompt key handling before palette/form handling:
  - `y` or `enter` starts execution with `Confirmed: true`,
  - `n` or `esc` clears the prompt and does not execute,
  - `ctrl+c` quits.
- Extend `startExecution` and `executeAction` to carry the confirmation boolean into `execution.ExecuteRequest`.
- Render confirmation copy from existing metadata: title, Capability ID, risk level, and description.

## Rationale

The TUI should adapt shared Capability metadata into an interactive user experience while keeping domain use cases and the Pipeline independent from Bubble Tea. Passing `Confirmed: true` only after a local user confirmation preserves the Pipeline as the policy enforcement point and avoids duplicating execution semantics.

## Risks

- The initial prompt is generic because no production mutating Capability currently provides domain-specific impact details. The prompt still uses stable metadata and can be enriched later without changing the Pipeline contract.
- `enter` confirms for keyboard convenience; the prompt labels both `y` and `enter` as confirmation and `n`/`esc` as cancellation.
