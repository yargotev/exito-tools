# Design: TUI command palette foundation

## Approach

Extend the existing side-effect-free TUI model with a small palette state: open/closed, query text, and palette-eligible Actions. The palette remains derived from Capability metadata, keeping Application Wiring and the Registry as the source of truth.

Primary navigation continues to use TUI-visible, people-facing Capabilities. Command Palette discovery uses command-palette-visible, people-facing Capabilities so agent-only or CLI-only capabilities are not promoted to everyday users.

Filtering is in-memory and case-insensitive over Action title, Capability ID, and domain. Selecting/executing Actions is deferred to a later task-runner slice.

## Decisions

- `/` opens the Command Palette from the shell.
- `esc` closes the Command Palette.
- Typing updates the palette query while open.
- `q` quits only when the palette is closed; inside the palette it is query text.

## Risks

- This slice intentionally does not execute selected Actions, so the palette is discovery-only until a later task execution slice.
