# Generic exit codes with specific error codes

The Exito Tools CLI Surface will use a small set of generic process exit codes for coarse automation decisions, while detailed failure semantics live in the structured `error.code` field. Agents should use the exit code to classify the broad failure category and the Error Code to make domain-specific decisions.
