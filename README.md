# DevServer

Module: `github.com/Ajayvtl/devserver`

## Layout

The repository is organized around:

- `apps/web`: Next.js frontend application
- `cmd/devserver`: command entrypoint only
- `internal/core`: application, registry wiring, runner, and lifecycle
- `internal/bootstrap`: server preparation and generic module orchestration
- `internal/config`: YAML, ENV, defaults, and validation
- `internal/logger`: structured logging with `zerolog`
- `internal/executor`: all command execution paths
- `internal/platform`: Linux distribution detection
- `internal/registry`: module registration and lookup
- `internal/state`: persistent installed-module state
- `internal/modules/*`: module implementations
- `assets/`: embedded runtime assets
- `configs/`: configuration files
- `templates/`: reusable templates
- `pkg/`: shared reusable packages

Design rules:

- command packages only expose `func Run(...)`
- business logic lives outside `cmd`
- operations accept `context.Context`
- logging is structured from the beginning
- modules are registered, not hardcoded
- every long-running flow should be expressed as tasks

## Frontend

The Next.js UI lives in [`apps/web/`](apps/web/). It starts at `/`, then routes through setup, login, and dashboard using mock services and reusable UI components so the backend can be connected incrementally later without reshaping the pages.

## Local Startup

Use the root launcher:

```bash
node dev.js
```
Test Indexer

It starts the Go backend and the Next.js frontend together from the project root.
