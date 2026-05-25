# Explicit TUI entrypoint

The Exito Tools TUI Surface opens through an explicit command such as `exito tui`; running `exito` without arguments should not silently enter an interactive session. This protects agents and scripts from accidentally blocking on an interactive Bubble Tea UI while still giving people a clear entrypoint.
