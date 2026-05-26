# Design: YAML Default Profile Foundation

## Technical Approach

Extend `internal/config.Resolve` so it selects the Configuration File path first, then loads a saved Default Profile from that file when it exists and `Options.SavedDefaultProfile` is blank. The existing profile resolver then applies the established precedence rules.

The parser intentionally supports only the narrow non-sensitive contract introduced by this slice:

```yaml
defaultProfile: prod
```

It ignores blank lines and comments, accepts simple quoted scalar values, and ignores nested or unsupported YAML content. This keeps the resolver dependency-free while leaving room for a later Viper-backed or full YAML parser behind the same resolver contract.

## Key Decisions

- `Options.SavedDefaultProfile` remains an explicit test/application seam and wins over file-loaded defaults when provided.
- Missing selected config files are ignored for now to preserve existing explicit path selection behavior until full config parsing validates file existence.
- Empty `defaultProfile` values are ignored and fall back to the next precedence source.
- Secrets remain environment/dotenv-only; YAML parsing is limited to the profile name.

## Testing Strategy

- Add resolver tests showing YAML `defaultProfile` is used as the saved default when no explicit/env profile exists.
- Add tests proving `--profile` and `EXITO_PROFILE` still win over YAML defaults.
- Add tests for comments/quoted values and blank ignored values.
- Run targeted config tests first, then `make test` and `go build ./cmd/exito` before reporting completion.
