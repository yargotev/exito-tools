# Structured warnings in envelope metadata

Exito Tools represents non-fatal issues as structured warnings in JSON Envelope metadata. Warnings have stable codes, messages, and optional details; they do not make the command fail, but they allow agents and the TUI Surface to react to partial data, fallbacks, deprecations, or other noteworthy conditions.
