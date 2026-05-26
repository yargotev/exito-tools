# TUI entrypoint foundation

## Why

The product decisions require interactive behavior to be opt-in through an explicit `exito tui` command. The current CLI root exposes machine-readable commands but has no TUI entrypoint yet.

## What Changes

- Add a narrow Bubble Tea-backed TUI Surface package.
- Add `exito tui` as the explicit interactive entrypoint while keeping bare `exito` as textual help.
- Render an initial task-first shell that shows the effective profile and curated primary actions from TUI-visible, people-facing Capabilities.

## Out of Scope

- Command palette search.
- Capability execution from the TUI.
- Profile switching/default profile persistence.
- Result filters and task runner states.
