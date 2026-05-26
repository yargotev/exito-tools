# Design: CLI Failure Exit Code Foundation

## Approach

The execution Pipeline already classifies capability failures into `ok:false` JSON Envelopes with stable `error.code` values. This slice keeps that structured contract unchanged and adds only the process-status bridge at the CLI Surface.

A new surface-owned `ExitError` carries a generic exit code. Commands write their JSON Envelope first; if the envelope is failed, they return `ExitError{Code: 1}`. The root command still uses `SilenceErrors`/`SilenceUsage`, and `cmd/exito` handles this sentinel by exiting with the code without printing an additional error line.

## Boundaries

- The Pipeline remains surface-independent and does not know about process exit codes.
- Domain packages continue returning structured errors only.
- Cobra validation errors remain ordinary errors so invalid flags/missing required flags do not emit JSON envelopes.
