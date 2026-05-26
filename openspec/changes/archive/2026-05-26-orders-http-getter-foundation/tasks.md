# Tasks

- [x] Add Orders HTTP getter using shared HTTP client infrastructure.
- [x] Map provider DTOs to domain-owned `orders.Order`.
- [x] Translate provider failures into structured Orders errors.
- [x] Wire `app.New` to use the HTTP getter when Orders provider config is present.
- [x] Add domain and application wiring tests for request path, auth, metadata, mapping, and failures.
- [x] Verify with `go test ./...`.
