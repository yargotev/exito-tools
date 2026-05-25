# Tasks: Initial Go Application Scaffold

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 260-340 |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | single PR |
| Delivery strategy | auto-forecast |
| Chain strategy | pending |

Decision needed before apply: No
Chained PRs recommended: No
Chain strategy: pending
400-line budget risk: Low

### Suggested Work Units

| Unit | Goal | Likely PR | Notes |
|------|------|-----------|-------|
| 1 | Ship runnable Go scaffold with tests | PR 1 | Single work unit; keep code, tests, and docs together |

## Phase 1: Foundation

- [x] 1.1 Create `go.mod` with module `github.com/yargotev/exito-tools` and add Cobra as the only startup dependency.
- [x] 1.2 Create `internal/capability/types.go` with minimal `Definition`, `StructuredError`, and generic `Envelope` skeletons only.
- [x] 1.3 Create `internal/registry/registry.go` with boot-time registration, `Finalize()`, immutable snapshot reads, and a stable post-finalize error.

## Phase 2: Application Wiring

- [x] 2.1 Create `internal/app/app.go` to build the registry, finalize it during boot, and expose the minimal Application seam to surfaces.
- [x] 2.2 Create `internal/surface/cli/root.go` with the English-only Cobra root command, no business subcommands, and help text that avoids claiming deferred commands exist.
- [x] 2.3 Create `cmd/exito/main.go` to call `app.New()`, build the root command, and execute the CLI surface without launching a TUI.

## Phase 3: Verification

- [x] 3.1 Create `internal/registry/registry_test.go` as table-driven tests for register-before-finalize, empty finalize, and post-finalize rejection from the registry spec.
- [x] 3.2 Create `internal/surface/cli/root_test.go` to run the bare root/help path, capture stdout, and verify English help text with no JSON envelope markers.
- [x] 3.3 Run `go test ./...` and verify the bootstrap, CLI-root, registry, and contract-foundation scenarios pass in the new module.

## Phase 4: Cleanup

- [x] 4.1 Update `openspec/changes/initial-go-application-scaffold/tasks.md` task checkboxes during apply and keep the slice limited to scaffold files only.
- [x] 4.2 Confirm visible strings in `internal/surface/cli/root.go` remain English-only and that no Orders, Geo, `run`, or `capabilities` command is implied.
