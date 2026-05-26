# Apply Progress: YAML Default Profile Foundation

## Status

Implemented.

## Completed

- Added Configuration Resolver support for reading a top-level `defaultProfile` scalar from the selected YAML Configuration File when no saved default profile is supplied directly.
- Preserved precedence: explicit profile, `EXITO_PROFILE`, saved YAML Default Profile, then `staging`.
- Documented the supported `defaultProfile` YAML key.
- Added focused resolver tests for YAML default profile loading, explicit/env override precedence, quoted values with comments, and blank fallback behavior.

## Verification

- `go test ./internal/config`
- `make test`
- `go build ./cmd/exito`
