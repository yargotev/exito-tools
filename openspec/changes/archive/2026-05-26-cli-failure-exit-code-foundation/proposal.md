# Proposal: CLI Failure Exit Code Foundation

## Summary

Make machine-readable capability commands return a generic non-zero process status when they emit a failed JSON Envelope, while keeping the structured error code in stdout as the detailed automation contract.

## Scope

- Add a small CLI-surface exit error type for generic process status mapping.
- Return the exit error from `run`, `orders get`, and `geo geocode-address` after writing failed envelopes.
- Update `cmd/exito` so structured exit errors terminate without extra stderr text that could confuse agents.
- Preserve successful command behavior and Cobra validation errors.

## Non-goals

- Add domain-specific exit codes for each structured error.
- Change JSON Envelope shapes or pipeline structured error mapping.
- Add logging configuration or TUI behavior.
