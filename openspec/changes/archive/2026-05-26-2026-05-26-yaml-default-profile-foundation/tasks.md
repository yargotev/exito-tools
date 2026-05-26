# Tasks: YAML Default Profile Foundation

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 180-260 |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | single slice |

## Phase 1: Spec and docs

- [x] 1.1 Add a Configuration Resolver delta requirement for YAML-saved Default Profile loading.
- [x] 1.2 Document the supported `defaultProfile` YAML key.

## Phase 2: Resolver implementation

- [x] 2.1 Select the config path before resolving the effective profile.
- [x] 2.2 Load a top-level `defaultProfile` scalar from an existing selected config file when no saved default was injected.
- [x] 2.3 Preserve explicit profile and `EXITO_PROFILE` precedence over the saved YAML default.

## Phase 3: Tests and verification

- [x] 3.1 Add focused resolver tests for YAML default profile behavior.
- [x] 3.2 Run `go test ./internal/config`.
- [x] 3.3 Run `make test` and `go build ./cmd/exito`.
