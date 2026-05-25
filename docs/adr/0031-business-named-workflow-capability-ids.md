# Business-named Workflow Capability IDs

Workflow Capabilities should use business-oriented Capability IDs rather than exposing `workflow.*` as a visible technical prefix. If the workflow has its own business identity, it may use a pseudo-domain such as `diagnostico.entrega`; if it clearly belongs to one domain, it should use the dominant domain such as `pedidos.diagnosticar-entrega`.
