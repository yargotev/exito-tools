# Design

## Approach

Add a small writer function in `internal/config` so persistence rules remain behind the Configuration Resolver boundary. The CLI surface will expose `exito config set-default-profile <profile>` and call the writer with shared boot flags.

## File selection

The writer uses existing `--config`, `EXITO_CONFIG`, local, and user discovery. If no existing selected file exists and no explicit path was supplied, it creates local `./exito.yaml` because that is the closest current project context.

## YAML handling

This slice intentionally supports only the top-level `defaultProfile` scalar. Existing files are updated line-by-line when the key exists, otherwise the key is appended. Secrets remain outside YAML.

## Result contract

The command writes a standard JSON Envelope containing the persisted profile, target path, and target source. The envelope metadata includes request ID, optional correlation ID, selected profile context, and duration.
