# ADR 0002: Workspace Manager and Dynamic Lifecycle

## Status
Accepted

## Context
In Phase 1, DevServer was hardcoded to manage a single workspace instantiated at startup (e.g., `core.NewIndexer(log, []core.WorkspaceSpec{{ID: "devserver", Root: root}})`). The `Indexer`, `WorkspaceWatcher`, and `WorkspaceProvider` expected static workspaces.

As we scale to support 100+ concurrent projects, repositories, servers, and containers, the platform requires a centralized orchestrator to handle the dynamic lifecycle of workspaces (Open, Close, Delete, Import, Export, Activate, Deactivate) without directly coupling API routes to individual subsystems.

## Decision
We introduce the `WorkspaceManager` (`internal/core/workspace_manager.go`) as the centralized lifecycle orchestrator.

1. **Responsibilities**:
   - Acts as the single entry point for all workspace lifecycle operations.
   - Coordinates the initialization and teardown across the `WorkspaceWatcher`, `Indexer`, and `EventBus`.
   
2. **Subsystem Upgrades**:
   - `Indexer` and `WorkspaceWatcher` will expose dynamic `AddWorkspace(id, root)` and `RemoveWorkspace(id)` methods to support on-the-fly registration.
   - `WorkspaceManager.Open(id, root)` orchestrates these calls and ensures the first synchronous index (`Refresh`) executes before returning.

3. **Event Emitting**:
   - The Manager publishes state transitions (e.g., `WorkspaceActivated`, `WorkspaceDeactivated`) to the `EventBus` so the Task Engine and AI Context Service can react appropriately.

## Consequences

### Positive
- **Scalability**: The backend can natively support multi-tenant workspaces and bulk imports.
- **Encapsulation**: HTTP Handlers and the Task Engine interact strictly with `WorkspaceManager` rather than micro-managing the Watcher and Indexer.
- **Foundation for Remote**: Opens the door to managing remote workspaces (SSH, Docker, Git) via abstraction in the future.

### Negative
- **State Drift Risk**: If `WorkspaceManager` fails halfway through an orchestration (e.g., adds to Indexer but Watcher errors), the workspace could end up in a zombie state. We must ensure robust error handling and transactional rollbacks.
