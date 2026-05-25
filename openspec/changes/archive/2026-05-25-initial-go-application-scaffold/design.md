# Design: Initial Go Application Scaffold

## Technical Approach

Implement the narrow architecture-first scaffold from the proposal: create the Go module, a tiny `cmd/exito` entrypoint, explicit `internal/app` wiring, an immutable boot/finalized registry, minimal shared capability contracts, and a Cobra-based CLI root that only renders English help. This satisfies the bootstrap, CLI root, registry, and contract foundation specs while deferring real domains, config resolution, TUI, `run`, `capabilities`, and JSON runtime emission.

## Architecture Decisions

| Decision | Options | Tradeoff | Choice / Rationale |
|---|---|---|---|
| Boot seam | Wire Cobra directly in `main` vs `internal/app` bootstrap | `main` is smaller but hides future composition | Use `internal/app` so explicit Application Wiring matches ADR 0025 and future domains register through one visible seam. |
| Registry lifecycle | Mutable map forever vs builder + finalized snapshot | Permanent mutability is simpler but breaks ADR 0026 | Use boot-time registration plus `Finalize()` returning an immutable registry/snapshot; post-finalize registration returns a stable error. |
| CLI root behavior | Custom help renderer vs Cobra defaults with narrow config | Custom text gives full control but adds maintenance | Use Cobra v1.9.1 default help generation with a constructed root command and no business subcommands; this matches ADR 0016 with minimal code. |
| Shared contracts | Full future contract model now vs skeleton types | Full model risks overdesign in first slice | Add only stable skeletons for capability metadata, structured error, and envelope-shaped result types so later slices can extend without moving packages. |

## Data Flow

`cmd/exito/main.go` → `internal/app.New()` → boot `registry.Builder` → finalize registry → `internal/surface/cli.NewRoot(app)` → Cobra `ExecuteContext()` → help text to stdout

No Operational Domain is loaded in this slice. The finalized registry may be empty, but the Application still boots and exposes the CLI Surface.

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `go.mod` | Create | Initialize module and Cobra dependency. |
| `cmd/exito/main.go` | Create | Minimal process entrypoint calling application boot and CLI execution. |
| `internal/app/app.go` | Create | Defines the Application struct and explicit bootstrap/finalization flow. |
| `internal/capability/types.go` | Create | Minimal capability metadata, structured error, and envelope skeletons. |
| `internal/registry/registry.go` | Create | Builder/finalized registry types and stable mutation errors. |
| `internal/surface/cli/root.go` | Create | Builds the Cobra root command with English help and future command seams. |
| `internal/surface/cli/root_test.go` | Create | Verifies bare-root help behavior with stdout capture and `SetArgs(nil)`. |
| `internal/registry/registry_test.go` | Create | Table-driven tests for registration, empty finalize, and post-finalize rejection. |

## Interfaces / Contracts

```go
package capability

type Definition struct {
    ID          string
    Title       string
    Description string
}

type StructuredError struct {
    Code    string
    Message string
}

type Envelope[T any] struct {
    OK    bool
    Data  *T
    Error *StructuredError
}
```

```go
package registry

type Builder interface {
    Register(capability.Definition) error
    Finalize() Registry
}

type Registry interface {
    All() []capability.Definition
}
```

The CLI package depends on `app` and finalized `registry` only; domain packages remain absent, preserving one-way surface dependencies.

## Testing Strategy

| Layer | What to Test | Approach |
|-------|-------------|----------|
| Unit | Registry registration/finalization behavior | Table-driven tests; assert empty finalize succeeds and post-finalize mutation returns the stable error. |
| Integration | Root CLI help path | Construct root command, set args to none/help, capture output, and assert English help text plus no JSON envelope markers. |
| E2E | Not in this slice | Defer compiled-command tests until explicit business commands/output runtime exist. |

## Migration / Rollout

No migration required. Roll out as a single foundation work unit that introduces the Go module and internal package seams without changing user data or existing runtime behavior.

## Open Questions

- [ ] What exact module path should `go.mod` use for this repository name and future imports?
