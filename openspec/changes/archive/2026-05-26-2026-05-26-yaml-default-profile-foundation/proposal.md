# Proposal: YAML Default Profile Foundation

## Summary

Teach the Configuration Resolver to load a saved Default Profile from the selected non-sensitive YAML Configuration File using a narrow top-level `defaultProfile` scalar.

## Motivation

The resolver already models Effective Profile precedence, including a saved Default Profile, but the saved value can only be injected programmatically. The product roadmap requires users to set a Default Profile explicitly, and the TUI Session Profile slice deliberately kept temporary session changes separate from saved defaults. Reading the saved default from YAML is the smallest next step toward persistent profile behavior without adding write flows or new dependencies.

## Scope

- Read a top-level `defaultProfile` value from the selected configuration file when the caller did not already provide `SavedDefaultProfile`.
- Preserve Effective Profile precedence: explicit `--profile`, then `EXITO_PROFILE`, then saved YAML Default Profile, then `staging`.
- Keep secret/provider values outside YAML; this slice only reads the non-sensitive profile name.
- Use a dependency-free narrow parser for the supported scalar so no production dependency is added.

## Out of Scope

- Writing or mutating configuration files.
- Full YAML document parsing, nested profile definitions, or domain overrides.
- Rebootstrapping provider clients when a running TUI Session Profile changes.
- Adding a CLI/TUI command to set the Default Profile.
