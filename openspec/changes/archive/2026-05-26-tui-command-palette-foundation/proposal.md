# TUI command palette foundation

## Why

The TUI Surface needs a discovery path beyond curated primary Actions. The PRD and glossary distinguish the Command Palette from result filters: it helps people find available Actions across domains without exposing every capability in primary navigation.

## What Changes

- Add a minimal Command Palette mode to the TUI model.
- Curate palette Actions from people-facing Capabilities with command-palette visibility.
- Allow opening/closing the palette and filtering Actions by typed query.

## Out of Scope

- Executing selected Actions from the palette.
- Keyboard selection/cursor movement.
- Result filtering within task outputs.
- Persisted user preferences for palette ordering.
