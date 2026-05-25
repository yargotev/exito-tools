# TUI Profile session selection is separate from Default Profile changes

The Exito Tools TUI Surface resolves its starting Profile with the same precedence as the CLI Surface and shows the active Profile to the user. Changing the Profile for the current TUI session is separate from explicitly setting a new Default Profile, preventing temporary navigation from silently changing future CLI or TUI behavior.
