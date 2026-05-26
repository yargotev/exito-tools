# Verify Report: TUI entrypoint foundation

## Change

`2026-05-26-tui-entrypoint-foundation`

## Result

✅ Passed

## Scope Verified

- `exito tui` is an explicit CLI subcommand and bare `exito` remains textual help.
- The TUI Surface is isolated under `internal/surface/tui` and Bubble Tea does not leak into Operational Domains.
- The initial TUI model renders the effective profile and primary Actions from TUI-visible, people-facing Capabilities.
- Agent-only Capabilities are not promoted as primary TUI Actions.

## Commands

- ✅ `go test ./...`
- ✅ `go build ./cmd/exito`
- ✅ `make lint`

## Issues

No critical issues found.
