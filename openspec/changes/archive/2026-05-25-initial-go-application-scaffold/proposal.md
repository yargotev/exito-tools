# Proposal: Initial Go Application Scaffold

## Intent

Start implementation after the docs-first phase by creating a runnable Go Application foundation that preserves the documented one-application, CLI/TUI-ready architecture and English-only visible contracts.

## Scope

### In Scope
- Initialize `go.mod` and `cmd/exito/main.go`.
- Add minimal `internal/app`, `internal/capability`, `internal/registry`, and `internal/surface/cli` seams.
- Add focused tests for root help and registry boot/finalization behavior.

### Out of Scope
- `orders.get`, `geo.geocode-address`, full Configuration Resolver, HTTP Infrastructure, and Bubble Tea TUI.
- `exito capabilities`, `exito run`, real JSON envelope emission, and real structured error translation beyond narrow skeleton types.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- None.

## Approach

Use the architecture-first slice from exploration: create only the seams required to boot the Application and show root CLI help. Keep domain packages absent or empty, keep Cobra confined to the CLI Surface, and keep registry/capability contracts minimal so later slices can add real Capabilities without reworking boot structure.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `go.mod` | New | Initialize module and minimal deps |
| `cmd/exito/main.go` | New | Executable entrypoint |
| `internal/app/` | New | Explicit Application Wiring seam |
| `internal/capability/` | New | Minimal contract and error/envelope skeletons |
| `internal/registry/` | New | Registration and finalization scaffold |
| `internal/surface/cli/` | New | Root Cobra command/help adapter |
| `**/*_test.go` | New | Root help and registry tests |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Scaffold expands into real product scope | Med | Reject non-foundation work in this slice |
| Early contracts become overdesigned | Med | Keep types minimal and ADR-aligned |

## Rollback Plan

Revert scaffold files as one work unit: remove `go.mod`, `cmd/exito`, and new `internal/*` foundation packages if boot structure proves misaligned.

## Dependencies

- Go toolchain and Cobra for root command wiring.

## Success Criteria

- [ ] `exito` builds and shows brief text help without launching TUI.
- [ ] Foundation tests pass for root help and immutable registry behavior.
