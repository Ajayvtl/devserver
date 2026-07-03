# ADR 0003: Task Engine

## Status
Accepted

## Context
The current task execution is simulated: `simulateTask()` in `server.go` hardcodes four progress phases with `time.Sleep(450ms)` and broadcasts via a task-specific `TaskHub` WebSocket. The existing `core.Runner` provides sequential execution with rollback support, but has no queue, no priority, no concurrency control, and no event integration.

As DevServer scales to support deployments, backups, database migrations, plugin installation, and AI-generated actions, every operation must be modeled as a task that is resumable, observable, and cancellable.

## Decision
We introduce `internal/tasks/` as a production-grade execution engine.

1. **Task Definition**: Tasks carry ID, name, workspace scope, priority, retry policy, timeout, dependencies, metadata, and progress reporting. They implement the existing `core.Task` interface so the `Runner`'s rollback support is preserved.

2. **Priority Queue**: FIFO by default, with priority override. Supports cancellation, pause/resume, and bounded capacity.

3. **Worker Pool**: Configurable concurrency. Workers execute tasks through the `core.Runner`, gaining rollback for free. Each worker emits typed events on every state transition.

4. **Event Integration**: Every lifecycle transition (Queued → Started → Progress → Completed/Failed/Cancelled/RolledBack) publishes a typed event to the EventBus. No polling.

5. **Event Stream**: A generalized WebSocket bridge subscribes to the EventBus and multiplexes all event types (task, workspace, plugin, AI) to connected frontends. Replaces the task-specific `TaskHub`.

## Consequences

### Positive
- Every operation in the system becomes a first-class task with observability.
- The frontend receives real-time updates through a single WebSocket channel.
- Rollback support from `core.Runner` is preserved without reimplementation.
- Future capabilities (deployments, AI actions, plugin install) slot in without new execution paths.

### Negative
- The task package introduces additional complexity in the runtime initialization.
- Queue ordering with priorities requires careful concurrency handling.
