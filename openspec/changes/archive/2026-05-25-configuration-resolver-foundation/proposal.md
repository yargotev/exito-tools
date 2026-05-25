# Proposal: Configuration Resolver Foundation

## Intent

Establish the first deep-module foundation after the Go scaffold: a deterministic Configuration Resolver for profile and source precedence. This unlocks later CLI/TUI execution without letting Viper or Cobra define product behavior.

## Scope

### In Scope
- Add `internal/config` resolver types and pure precedence logic.
- Resolve Effective Profile from explicit option, `EXITO_PROFILE`, saved default, then `staging`.
- Discover Configuration File path from explicit path, `EXITO_CONFIG`, local `./exito.yaml`, user config, then defaults.
- Expose dotenv credential layer order without loading secrets.
- Add focused table-driven Go tests.

### Out of Scope
- Parsing YAML values with Viper.
- Loading dotenv files into process environment.
- Domain-specific Geo/Orders config.
- CLI flags wired into command execution.

## Capabilities

### New Capabilities
- `configuration-resolver`: shared configuration precedence and profile resolution behavior.

### Modified Capabilities
- None.

## Approach

Implement a small dependency-free resolver with injectable working directory, home directory, environment map, explicit options, and saved default profile. Keep it usable by all surfaces and test without real user files.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/config` | New | Configuration resolver package and tests. |
| `openspec/changes/configuration-resolver-foundation` | New | SDD artifacts for this slice. |
| `openspec/specs/configuration-resolver` | New | New permanent spec after archive. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Resolver behavior diverges from docs | Low | Encode documented precedence in tests and specs. |
| Premature YAML/dotenv implementation leaks secrets | Low | Only discover sources; do not read secret values. |

## Rollback Plan

Revert this change's code and artifacts; no runtime user data or committed secrets are changed.

## Dependencies

- Existing Go module and test runner.

## Success Criteria

- [ ] `go test ./...` passes.
- [ ] Resolver precedence is covered by table-driven tests.
- [ ] No secret files are read or committed.
