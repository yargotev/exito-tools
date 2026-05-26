# Proposal: TUI result filter foundation

## Summary

Add a narrow TUI result filter mode that lets a person refine the currently displayed successful task result without reopening the Command Palette or re-running the Capability.

## Motivation

The TUI can now discover Actions, collect required string input, and execute selected Capabilities through the shared Pipeline. The next people-facing foundation is to distinguish Command Palette discovery from Result Filters by allowing the user to filter data already loaded inside a task.

## Scope

- Store successful Pipeline result data in TUI task state.
- Render successful result data as stable text rows.
- Add a result filter mode that opens only when successful result rows exist.
- Filter displayed result rows by typed query.
- Keep filtering local to the TUI; do not change Capability contracts or domain use cases.

## Out of Scope

- Domain-specific result tables or rich widgets.
- Filtering before task execution.
- Server-side filtering or re-executing Capabilities.
- Persisted filters across tasks.
