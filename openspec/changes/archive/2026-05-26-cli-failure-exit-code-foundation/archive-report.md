# Archive Report: CLI Failure Exit Code Foundation

## Change

`2026-05-26-cli-failure-exit-code-foundation`

## Archived

2026-05-26

## Summary

Synced the CLI run delta spec into the source-of-truth CLI run spec and archived the completed change.

## Specs Synced

| Domain | Action | Details |
|--------|--------|---------|
| `cli-run` | Updated | Modified the generic run command requirement so successful envelopes exit successfully and failed unknown-capability envelopes return a generic failure exit status. |

## Verification

- ✅ `go test ./...`
- ✅ `go build ./cmd/exito`
- ✅ No critical issues in `verify-report.md`

## Archive Contents

- proposal.md ✅
- specs/ ✅
- design.md ✅
- tasks.md ✅ (5/5 tasks complete)
- verify-report.md ✅

## Source of Truth Updated

- `openspec/specs/cli-run/spec.md`

## SDD Cycle Complete

The change has been fully planned, implemented, verified, and archived.
