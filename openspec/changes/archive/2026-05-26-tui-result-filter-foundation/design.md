# Design: TUI result filter foundation

## Approach

Extend the TUI task state to retain successful result data from the shared Pipeline. The View renders result data as deterministic text rows after a success state. A lightweight filter mode is toggled with `f` when successful result rows exist. While the filter is active, typed input updates a local query and only matching result rows are rendered.

## Decisions

- Result filtering is local to loaded TUI task data and never re-executes the Capability.
- The filter key is `f` because `/` remains reserved for Command Palette discovery.
- `esc` closes the active result filter and clears its query.
- The first renderer is intentionally generic and deterministic: maps render as sorted `key: value` rows, slices render their elements, and scalar data renders as one row.

## Risks

- Generic text rows are not a substitute for future domain-specific result layouts.
- Large result sets are not virtualized in this foundation slice.
