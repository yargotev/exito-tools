# Cursor-based pagination without default auto-pagination

List-style Capabilities in Exito Tools use explicit cursor-based pagination when the backend supports it. Commands accept pagination inputs such as `limit` and `cursor`, return pagination metadata such as `nextCursor` and `hasMore`, and do not auto-paginate by default to avoid surprising agents with unbounded work.
