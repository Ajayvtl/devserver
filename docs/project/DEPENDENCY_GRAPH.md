# Dependency Graph

This document details the macro-level dependencies, component interactions, and startup order of the DevServer Runtime.

## Runtime Dependency Graph
*See `docs/architecture/runtime_baseline_report.md` for the graphical Mermaid chart.*

1. **`events.Bus`**: The core dependency of almost all asynchronous components.
2. **`tasks.EventStream`**: Depends on `events.Bus`.
3. **`providers.Manager`**: Independent.
4. **`core.Indexer`**: Depends on `events.Bus` and `logger`.
5. **`core.WorkspaceProvider`**: Depends on `core.Indexer`.
6. **`tasks.Engine`**: Depends on `events.Bus` and internal `tasks.Registry`.
7. **`commands.Engine`**: Depends on `events.Bus`, `tasks.Engine`, `state.StoreDB`, and `capabilities.Registry`.
8. **`core.APIServer`**: The final consumer, depending on all of the above to handle REST/HTTP.

## Startup / Shutdown Order
The `runtime.Coordinator` guarantees the following sequential boot order, and shuts down in exact reverse order:
1. `tasks.EventStream`
2. `providers.Manager`
3. `core.WorkspaceProvider`
4. `core.Indexer`
5. `tasks.Engine`
6. `commands.Engine`
7. `core.APIServer`

## Command Flow Architecture (Environment Pattern)
1. **User Intent**: The UI selects an Environment (e.g., "Staging") and issues `POST /api/commands { capability: "service.start", target: "redis", environment: "staging" }`.
2. **API Server**: Validates payload and passes to `commands.Engine`.
3. **Command Engine**: Resolves the capability and retrieves the active `runtime.EnvironmentContext`.
4. **Translation**: Converts the Command into a `tasks.Task` injecting the `executor.ExecutionContext` bound to that environment (e.g., `SSHRuntime`).
5. **Task Engine**: Accepts the Task, queues it, and a worker executes the Provider logic using the specific runtime.
6. **Event Bus**: Throughout execution, events (`TaskStarted`, `TaskProgress`, `TaskCompleted`) are fired.
7. **Event Stream**: Websocket pushes state updates back to the UI.

## Event Flow Architecture
* `Event` structs follow the standard envelope: `ID, Type, Source, Time, Payload`.
* Subsystems push to the synchronous `events.Bus`.
* Background consumers (like `commands.Engine.runInternal()`) listen for `TaskCompleted` to mutate database state.
* `tasks.EventStream` translates these backend events to `wireEvent` JSON payloads and pushes them over the active Websocket to the React UI.
