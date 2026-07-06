# Component Registry

This registry tracks every active component managed by the DevServer Runtime.

| Component Name | Owner Package | Interfaces Implemented | Dependencies | Status | Test Coverage |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `tasks.EventStream` | `internal/tasks` | `runtime.Component` | `events.Bus` | `Active` | Integration |
| `providers.Manager` | `internal/providers` | `runtime.Component` | None | `Active` | Integration |
| `core.WorkspaceProvider` | `internal/core` | `runtime.Component` | `core.Indexer` | `Active` | Integration |
| `core.Indexer` | `internal/core` | `runtime.Component` | `events.Bus`, `logger` | `Active` | Unit, Integration |
| `tasks.Engine` | `internal/tasks` | `runtime.Component` | `events.Bus`, `tasks.Registry` | `Active` | Unit, Integration |
| `commands.Engine` | `internal/commands` | `runtime.Component` | `events.Bus`, `tasks.Engine`, `DB`, `caps` | `Active` | Integration |
| `core.APIServer` | `internal/core` | `runtime.Component` | `WorkspaceProvider`, `TaskEngine`, etc. | `Active` | Integration |
| Component Name | Path | Status | Dependencies | Health | Active Issues |
|---|---|---|---|---|---|
| `editor.Orchestrator` | `internal/orchestrator` | *Pending* | `TaskEngine`, `ProcessManager` | `Pending` | None |
| `terminal.Manager` | `internal/terminal` | *Pending* | *Pending* | `Pending` | None |
| `ai.Runtime` | `internal/ai` | *Pending* | `WorkspaceProvider`, `Indexer` | `Pending` | None |
| `runtime.EnvironmentContext` | `internal/runtime` | *Pending* | `ExecutionContext`, `WorkspaceProvider` | `Pending` | None |
| `executor.Local` | `internal/executor/local` | *Pending* | OS specific drivers | `Pending` | None |
| `executor.Docker` | `internal/executor/docker` | *Pending* | Docker SDK | `Pending` | None |
| `apps/web` | `apps/web` | *Pending* | Next.js dashboard | `Pending` | None |
