# Proposal: TUI task cancellation foundation

## Summary

Add the narrow foundation for cancelling an in-flight TUI Capability Execution from the interactive surface. The slice keeps cancellation local to the TUI Task Runner model, uses Go context cancellation for the shared Pipeline execution path, and renders a deterministic cancelled state.

## Problem

The TUI already renders loading, success, and failure task states, but the PRD requires long-running TUI actions to be cancellable. Without a local cancellation state, pressing escape during execution has no explicit effect and users cannot distinguish a deliberately cancelled task from a completed or failed task.

## Scope

- Add a cancelled TUI task state.
- Allow `esc` during task loading to cancel the execution context.
- Ignore late completion messages for tasks already marked cancelled.
- Render the cancelled state in English.
- Cover the behavior with focused TUI model tests.

## Out of Scope

- Provider-specific cancellation UX beyond context propagation.
- Background job management, retries, or task history.
- CLI cancellation behavior changes.
- Profile switching or default profile persistence.

## Rollback

Revert this change to restore the previous TUI loading/success/failure-only task model. No persistent data or public Capability IDs are changed.
