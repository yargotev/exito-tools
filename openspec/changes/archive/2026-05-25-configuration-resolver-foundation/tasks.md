# Tasks: Configuration Resolver Foundation

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 280-380 |
| 400-line budget risk | Low |
| Chained PRs recommended | No |
| Suggested split | single PR |

## Phase 1: Resolver core

- [x] 1.1 Create `internal/config` resolver types for options, selected sources, effective config, and credential layers.
- [x] 1.2 Implement Effective Profile precedence: explicit profile, `EXITO_PROFILE`, saved default, `staging`.
- [x] 1.3 Implement Configuration File discovery: explicit path, `EXITO_CONFIG`, local `./exito.yaml`, user config, defaults.
- [x] 1.4 Implement profile-aware credential layer descriptors for process env, `.env.{profile}`, and `.env`.

## Phase 2: Tests and validation

- [x] 2.1 Add table-driven tests for profile precedence.
- [x] 2.2 Add table-driven tests for config file discovery precedence.
- [x] 2.3 Add tests for dotenv layer ordering.
- [x] 2.4 Run `gofmt`, `go test ./...`, and `go build ./cmd/exito`.
