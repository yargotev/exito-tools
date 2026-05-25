# Shared HTTP infrastructure with domain API clients

Exito Tools will use shared low-level HTTP infrastructure for cross-cutting concerns such as base URLs, authentication headers, timeouts, retries, and request metadata, while domain-specific API clients live near the domains that consume those external APIs. This keeps infrastructure reusable without centralizing domain semantics in a generic HTTP layer.
