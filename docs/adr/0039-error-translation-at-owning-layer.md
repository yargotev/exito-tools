# Error translation at the owning layer

Exito Tools translates errors into Structured Errors at the layer that understands their meaning. Shared HTTP infrastructure handles technical failures such as timeouts or network errors, Domain API Clients translate external API semantics, and Use Cases may enrich errors with domain context before surfaces present them.
