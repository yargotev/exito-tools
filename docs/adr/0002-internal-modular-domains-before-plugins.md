# Internal modular domains before plugins

Exito Tools will start with operational domains implemented as internal modules that can register their Use Cases, CLI Commands, and TUI Actions. We are not building dynamic external plugins yet because the current need is extensibility inside the repo without the complexity of runtime loading, versioning, or a plugin marketplace.
