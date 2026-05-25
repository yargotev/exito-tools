## Exploration: initial-go-application-scaffold

### Current State
The repository is still docs-first: `CONTEXT.md`, PRD, ADRs, capability notes, OpenSpec bootstrap, and env templates exist, but there is no `go.mod`, no Go packages, no test runner, and no executable entrypoint yet. The documented architecture expects one Go Application with explicit Application Wiring, a neutral Capability Core, an immutable Capability Registry, CLI/TUI Surface adapters, and English-only visible contracts.

### Affected Areas
- `CONTEXT.md` — source of truth for product language and architecture terms the scaffold must preserve.
- `docs/prd.md` — defines the recommended implementation order: deep modules first, CLI contracts before TUI.
- `docs/package-layout.md` — defines the target package boundaries for the scaffold.
- `openspec/config.yaml` — states the repo currently has no Go module/test tooling and prefers deep-module foundations first.
- `cmd/exito/main.go` — new executable entrypoint for the Application.
- `internal/app/` — new Application Wiring seam for explicit bootstrapping.
- `internal/capability/` — new neutral Capability Core types needed before real Capabilities exist.
- `internal/registry/` — new immutable Capability Registry scaffold for future discovery and execution.
- `internal/surface/cli/` — new Cobra adapter for root help and future command registration.

### Approaches
1. **Surface-first bootstrap** — Add only `go.mod`, Cobra, and a minimal `exito` root command.
   - Pros: Smallest slice, fastest path to a runnable binary, low review risk.
   - Cons: Delays the documented architecture, risks rework when Capability Core and Registry are added, gives little structure for future domain work.
   - Effort: Low

2. **Architecture-first foundation slice** — Add `go.mod`, package skeleton, minimal Capability/Registry/App Wiring contracts, and a root CLI that shows help without implementing domains yet.
   - Pros: Matches PRD sequence, preserves domain/surface boundaries from day one, creates safe seams for `orders.get`, `geo.geocode-address`, `capabilities`, and `run` next.
   - Cons: Slightly larger first diff, requires discipline to keep placeholders narrow and avoid premature TUI/config complexity.
   - Effort: Medium

### Recommendation
Choose **Architecture-first foundation slice**, but keep it narrow: initialize the Go module, add the documented package layout, implement only the minimal shared contracts needed to boot the Application and root CLI, and add tests for root help plus Registry basics. Defer real Configuration Resolver behavior, `orders.get`, `geo.geocode-address`, `exito capabilities`, `exito run`, HTTP infrastructure, and all TUI work to later slices so the first implementation stays under the 400-line review budget.

### Risks
- The scaffold can sprawl if it tries to include configuration resolution, domain clients, or Bubble Tea in the same slice.
- Placeholder types may calcify into the wrong contract if they are more detailed than the ADRs currently justify.

### Ready for Proposal
Yes — propose a first slice that creates the runnable Go foundation only: `go.mod`, `cmd/exito`, `internal/app`, `internal/capability`, `internal/registry`, `internal/surface/cli`, and focused tests that prove root help behavior and registry immutability/registration rules.
