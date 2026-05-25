# Run command accepts complete input objects

`exito run <capability-id>` accepts a complete input object through `--input-json`, `--input-file`, or piped stdin, while explicit domain commands expose typed flags and arguments. This keeps the generic run path aligned with neutral Input Schemas and avoids duplicating dynamic flag behavior already provided by explicit commands.
