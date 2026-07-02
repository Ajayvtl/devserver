# DevServer

Module: `github.com/Ajayvtl/devserver`

## Layout

The repository is now organized around:

- `cmd/devserver`: CLI entrypoint
- `internal/`: application packages such as `app`, `bootstrap`, `config`, `doctor`, `installer`, `logger`, `runner`, `state`, `system`, and `modules`
- `assets/`: runtime assets and scripts
- `pkg/`: shared external packages
- `templates/`: reusable templates
- `configs/`, `scripts/`, `docs/`, `tests/`, and `build/`: supporting project folders

The `internal/modules/` area is reserved for module-specific implementations such as `nginx`, `node`, `python`, `php`, `redis`, `postgres`, `mysql`, `ssl`, `firewall`, `github`, `backup`, `monitor`, and `utils`.
