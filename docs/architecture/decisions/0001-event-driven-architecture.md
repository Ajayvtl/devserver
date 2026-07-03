# ADR 0001: Event-Driven Platform Architecture

## Status
Accepted

## Context
DevServer is transitioning from Phase 1 (Foundation) to Phase 2 (Platform Core). Previously, subsystems like the Workspace Indexer utilized polling loops (e.g., `time.NewTicker`) to detect file changes. Additionally, the Workspace Provider and HTTP API handlers had rigid, coupled dependencies directly on the Indexer and its internals. 

As the platform scales to support hundreds of projects, a task engine, plugin registry, and an AI Context Engine, point-to-point hardcoded relationships will lead to spaghetti code, scalability bottlenecks, and performance issues (e.g., CPU churn from polling). We need a decoupled, event-driven foundation.

## Decision
We are adopting a strict **Event-Driven Architecture** for the backend runtime. 

1. **Central Event Bus (`internal/events`)**: 
   - A thread-safe, publish/subscribe Event Bus acts as the central nervous system.
   - All significant state changes must emit strongly-typed events (e.g., `WorkspaceChanged`, `TaskStarted`, `ContextUpdated`).
   - No service should directly call another service to trigger side effects. Everything subscribes to the Event Bus.

2. **Filesystem Watcher (`internal/core/watcher.go`)**:
   - Polling is completely removed.
   - The platform now relies on `fsnotify` for real-time filesystem events.
   - The watcher debounces bursts of file changes and emits a `WorkspaceChanged` event to the bus.

3. **Decoupled Workflows**:
   - **Old**: Polling Ticker -> Scan -> Write Context
   - **New**: `fsnotify` -> `WorkspaceWatcher` -> Emit `WorkspaceChanged` -> `Indexer` Subscribes -> Scans/Writes -> `Provider` Invalidates Cache.

## Consequences

### Positive
- **Performance**: Eliminates idle CPU usage from polling.
- **Maintainability**: Subsystems (Indexer, AI Context, Task Queue, Plugins) can be built independently and react to events without modifying existing code.
- **Scalability**: The event bus lays the foundation for WebSockets to stream real-time updates directly to the frontend.

### Negative
- **Complexity**: Debugging state transitions requires tracing events rather than following direct function calls.
- **Concurrency**: Event handlers must be written to be thread-safe and non-blocking.
