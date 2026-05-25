# Explicit domain command convention

Explicit CLI commands in Exito Tools follow `exito <domain> <action> [flags]`, where the domain is canonical, the action uses kebab-case when needed, and flags are derived from the Capability Input Schema. Each explicit command maps to a stable Capability ID, such as `exito pedidos consultar` mapping to `pedidos.consultar`.
