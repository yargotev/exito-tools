# Environment-based credentials instead of OS keychain

Exito Tools will read credentials from non-committed `.env` files and environment variables rather than OS-specific secure credential stores. This matches the existing operating style for the project and keeps CLI, TUI, local development, and agent/CI execution aligned without depending on macOS Keychain, Windows Credential Manager, or Linux Secret Service.
