# Apply Progress: Regionalized Intelligent Search flow

## TDD Cycle Evidence

| Task | Evidence |
| --- | --- |
| Geo region IDs | Added failing/covering assertion in `internal/domain/geo/geo_test.go`, then implemented `regions[]` extraction. |
| Workflow orchestration | Added `internal/workflow/regionalized_intelligent_search_test.go` for region → segment → search, confirmation, no-region stop, and token redaction. |
| App registration | Extended `TestNewWiresBootCapabilities` for the new workflow capability. |
| CLI routing | Added CLI tests for confirmation-required and confirmed regionalized command input mapping. |

## Implementation Notes

- Capability: `catalog.regionalized-intelligent-search-products`.
- CLI: `exito catalog intelligent-search regionalized-products`.
- Risk: `safe-write`, confirmation-required.
- Raw segment cookie is used only as an internal input to Intelligent Search and is not included in the workflow result.
