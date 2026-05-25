# Configuration file discovery precedence

Exito Tools discovers non-sensitive YAML configuration by explicit path first, then `EXITO_CONFIG`, then local `./exito.yaml`, then user config at `~/.config/exito-tools/config.yaml`, then internal defaults. A local project config wins over user config because agents and commands should respect the repository context in which they run.
