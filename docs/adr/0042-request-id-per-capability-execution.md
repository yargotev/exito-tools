# Request ID per Capability execution

Each Exito Tools Capability Execution will generate or propagate a `requestId` used in JSON envelope metadata, structured logs, and outbound HTTP request metadata when applicable. `requestId` identifies a single CLI/TUI execution path; distributed tracing terminology such as `traceId` is deferred until needed.
