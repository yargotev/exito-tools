# Change: CLI Default Profile Command

## Summary

Add a narrow CLI command that lets users explicitly persist the saved Default Profile in the non-sensitive YAML Configuration File.

## Motivation

The resolver can already read `defaultProfile` from YAML. The next slice should provide an explicit user path to update that value so future CLI and TUI sessions use the preferred profile without relying on manual file edits.

## Scope

- Add `exito config set-default-profile <profile>`.
- Write only the top-level `defaultProfile` YAML scalar.
- Select the target file using existing configuration discovery; when no file exists, create local `./exito.yaml`.
- Return a standard JSON Envelope on success.

## Non-goals

- Full YAML parsing or arbitrary config editing.
- Writing secrets or provider tokens to YAML.
- TUI default-profile persistence.
