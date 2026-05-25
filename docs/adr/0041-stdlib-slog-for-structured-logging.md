# Stdlib slog for structured logging

Exito Tools will use Go's standard `log/slog` package for structured logging. It avoids an extra dependency, supports levels and structured fields, and is sufficient for CLI/TUI observability when logs are kept separate from machine-readable command output.
