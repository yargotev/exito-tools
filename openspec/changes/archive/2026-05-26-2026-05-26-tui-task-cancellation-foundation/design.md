# Design: TUI task cancellation foundation

## Context

The TUI model currently starts Capability execution by returning a Bubble Tea command that calls the shared execution Pipeline with the model context. Task state supports idle, loading, success, and failure. The PRD and TUI test guidance call out cancelled task states.

## Approach

- Extend `taskStatus` with `cancelled`.
- Store a task-scoped `context.CancelFunc` on the model when starting execution.
- Wrap the model context with `context.WithCancel` in `startExecution` and pass the cancellable context into the execution command closure.
- Handle `esc` while `task.Status == loading` before opening/handling other modes:
  - call the cancel function when present,
  - clear the cancel function,
  - set the task status to `cancelled`, preserving the Capability ID,
  - clear forms, palette, and result filter so the cancelled task state is unambiguous.
- When an `actionExecutedMsg` arrives after a task has already been cancelled for the same Capability ID, ignore it so late completion cannot overwrite the cancelled view.

## Rationale

Keeping cancellation in the TUI Task Runner model preserves domain and Pipeline independence from Bubble Tea. Context cancellation reuses the existing context-aware Capability Execution contract, allowing provider clients and handlers to stop when they respect the context.

## Risks

- Bubble Tea commands may still return a late message after cancellation. The model explicitly ignores late messages for the cancelled Capability ID.
- Handlers that do not observe context cancellation may continue work in the background until they return; this slice establishes the surface contract without provider-specific guarantees.
