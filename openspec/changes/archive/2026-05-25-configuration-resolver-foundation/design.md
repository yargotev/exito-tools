# Design: Configuration Resolver Foundation

## Technical Approach

Create `internal/config` as a dependency-free deep module with deterministic resolver inputs. The resolver does not parse YAML or load secrets yet; it computes selected profile, selected configuration source, candidate file paths, and credential layer order.

## Key Decisions

- Use injectable `WorkDir`, `HomeDir`, and `Env` for deterministic tests and future CLI/TUI reuse.
- Keep source selection independent from Cobra/Viper so product precedence is explicit.
- Treat explicit config path and `EXITO_CONFIG` as selected values even if the file does not currently exist; this lets later parsing return a precise file error.
- Keep dotenv values unread in this slice to avoid secret handling before the credential loader contract exists.

## Components

| Component | Responsibility |
|-----------|----------------|
| `Options` | Resolver inputs from flags, environment, saved defaults, and filesystem roots. |
| `Effective` | Resolved profile, selected config source, candidates, and credential layers. |
| `Resolve` | Applies precedence and returns deterministic resolution output. |

## Testing Strategy

Use table-driven Go tests with temporary directories and injected environment maps. Tests cover profile precedence, configuration file discovery, explicit path behavior, and dotenv layer ordering.

## Deferred Work

- YAML parsing and validation.
- Viper as an implementation helper behind this resolver.
- Dotenv loading into an isolated credential map.
- CLI persistent flags wired through application boot.
