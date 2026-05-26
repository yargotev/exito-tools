# Archive Report: Orders HTTP getter foundation

## Change

`orders-http-getter-foundation`

## Archived

2026-05-26

## Summary

Synced the Orders HTTP getter foundation delta specs into the source-of-truth specs and archived the completed change.

## Specs Synced

| Domain | Action | Details |
| --- | --- | --- |
| `capability-contract-foundation` | Updated | Added `orders.get domain execution` requirement with configured provider success and not-found scenarios. |
| `application-bootstrap` | Updated | Added `Explicit domain dependency wiring` requirement for configured Orders provider wiring and missing-config behavior. |
| `http-infrastructure` | Updated | Added `Provider domain clients use shared HTTP infrastructure` requirement for Orders HTTP getter authenticated metadata-bearing requests. |

## Archive Contents

- proposal.md ✅
- design.md ✅
- tasks.md ✅ (6/6 complete)
- verify-report.md ✅
- specs/ ✅

## Verification

- ✅ `go test ./...`
- ✅ `make test`
- ✅ `make lint`
- ✅ `go build ./cmd/exito`

## Source of Truth Updated

- `openspec/specs/capability-contract-foundation/spec.md`
- `openspec/specs/application-bootstrap/spec.md`
- `openspec/specs/http-infrastructure/spec.md`

## SDD Cycle Complete

The change has been fully planned, implemented, verified, and archived.
