# Package Layout

```text
cmd/exito/                 # main
internal/app/              # application wiring
internal/config/           # Configuration Resolver
internal/registry/         # Capability registry
internal/capability/       # Capability, schemas, results, structured errors
internal/platform/httpclient/ # shared HTTP infrastructure

internal/domain/orders/    # Orders operational domain
internal/domain/geo/       # Geo operational domain
internal/workflow/         # cross-domain Workflow Capabilities

internal/surface/cli/      # Cobra adapter
internal/surface/tui/      # Bubble Tea adapter
internal/presenter/        # JSON/Markdown/etc presenters
```

## Dependency Rule

`internal/domain/*` must not import `internal/surface/*`, Cobra, or Bubble Tea. Surfaces depend on domains and capability contracts, not the other way around.
