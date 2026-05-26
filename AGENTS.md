# AGENTS.md

## Purpose

Exito Tools is one Go application for Exito operational workflows, exposed through multiple interaction surfaces. Keep the product unified: CLI and TUI are surfaces over the same capabilities and use cases, not separate apps.

## Operating mode

- Prefer retrieval-led reasoning: before changing architecture, contracts, or terminology, read the relevant local source of truth instead of relying on memory.
- Start by exploring the project context, then read only the task-relevant docs listed below.
- Keep changes narrow and reviewable. Avoid repo-wide rewrites unless explicitly requested.
- Do not add production dependencies, change public contracts, or introduce destructive behavior without explicit user approval.
- Keep user-facing labels, help, errors, and visible contracts in English until i18n exists.

## Local docs index

Read these files as needed rather than duplicating their content here:

| Need | Read |
| --- | --- |
| Domain language, product terms, architectural vocabulary | `CONTEXT.md` |
| Product roadmap and intended capability sequence | `docs/prd.md` |
| Package boundaries and dependency rule | `docs/package-layout.md`, `docs/adr/0024-go-package-layout.md` |
| Specific architecture decisions | `docs/adr/*.md` |
| Current source-of-truth requirements | `openspec/specs/*/spec.md` |
| Active change proposals/tasks/designs | `openspec/changes/<change-id>/` |
| Configuration behavior | `docs/configuration.md`, `openspec/specs/configuration-resolver/spec.md` |
| Capability contracts | `docs/capabilities/*.md`, `openspec/specs/capability-*/spec.md` |
| Agent skill inventory | `.atl/skill-registry.md` |

## Architecture rules

- Domain packages under `internal/domain/*` must not import Cobra, Bubble Tea, or `internal/surface/*`.
- Surface packages adapt capabilities into CLI/TUI behavior; business behavior belongs in domain use cases or shared execution/application layers.
- Use explicit application wiring in `internal/app`; do not introduce hidden side-effect registration.
- Treat the capability registry as immutable after boot.
- Preserve stable Capability IDs. Use `<domain>.<kebab-case-action>` and introduce a versioned ID for incompatible changes.
- Map external DTOs into domain-owned models/results before exposing them to use cases or surfaces.
- Keep Viper, if used, behind the Configuration Resolver. Product precedence is explicit, not library-defined.
- Secrets come from environment variables or non-committed dotenv files, never committed YAML/docs.

## CLI/TUI contracts

- CLI is machine-first. JSON is the default command output shape for command results.
- Do not contaminate stdout JSON envelopes with logs or debug output. Logs/diagnostics go to stderr or a non-disruptive TUI path.
- CLI failures use stable structured error codes and generic exit-code categories.
- TUI actions should execute the same capabilities/use cases as CLI commands.
- Long-running capability execution must be context-aware and cancellable.
- Risk/destructive actions require explicit confirmation; non-interactive CLI flows must not silently proceed.

## Development commands

| Task | Command |
| --- | --- |
| Format | `make fmt` |
| Tidy modules | `make tidy` |
| Lint | `make lint` |
| Test | `make test` |
| Pre-commit checks | `make precommit` |
| Build CLI | `go build ./cmd/exito` |
| Coverage | `go test ./... -cover` |

Run the smallest relevant verification first. Before declaring implementation complete, run `make test` at minimum; run `make precommit` when changes touch formatting, lint-sensitive code, or before commit/PR handoff.

## Spec-driven workflow

- Align proposals, specs, designs, tasks, and implementation with `CONTEXT.md`, `docs/prd.md`, ADRs, and existing OpenSpec specs.
- For requirement changes, update the appropriate `openspec/changes/<change-id>/` artifacts before or alongside implementation.
- Use Given/When/Then scenarios and RFC 2119 language in specs.
- Preserve documented package layout and domain independence unless an ADR/change explicitly updates them.
- Archive/sync completed changes only after verification passes and the user asks to archive or close the slice.

## Git and collaboration

- Check `git status --short` before editing and before reporting completion.
- Do not overwrite user changes. If unexpected modified files appear, stop and ask.
- Keep commit messages clean and without AI/co-author attribution.
- Do not commit unless the user asks.
