# Design: TUI entrypoint foundation

## Approach

Introduce `internal/surface/tui` as the Bubble Tea adapter. The initial model is intentionally read-only: it receives the already bootstrapped Application, selects primary actions from Capability metadata, and renders a small shell with the current profile and actions. This keeps Operational Domains independent from Bubble Tea and preserves Application Wiring as the source of registered Capabilities.

The CLI Surface adds a `tui` subcommand that bootstraps the Application using the same `--config` and `--profile` flags as the other commands, then starts the TUI program. Bare `exito` remains textual help to avoid blocking agents.

## Decisions

- Use Bubble Tea only inside `internal/surface/tui`.
- Treat primary actions as Capabilities with TUI visibility and people audience.
- Keep the initial model side-effect-free; future slices can add task execution and navigation.

## Risks

- Running Bubble Tea in non-interactive environments may fail; this is acceptable because the TUI is explicitly opt-in and should not be used by automation unless a terminal is available.
