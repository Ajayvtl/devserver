# Runtime Migration: Dependency Graph

## 1. Current State (Fragmented Hard-Wiring)

Currently, subsystems are instantiated and wired manually in `cmd/devserver/run.go`. Subsystems cannot dynamically discover one another, and initialization order is hardcoded.

```mermaid
graph TD
    Main[cmd/devserver/run.go]
    
    Main -->|instantiates| EvtStream[core.EventStream]
    Main -->|instantiates| CapReg[capabilities.Registry]
    Main -->|instantiates| TaskReg[tasks.Registry]
    Main -->|instantiates| ProvMgr[providers.Manager]
    Main -->|instantiates| Indexer[indexer.WorkspaceIndexer]
    Main -->|instantiates| WSProvider[core.WorkspaceProvider]
    Main -->|instantiates| Srv[core.Server]
    
    %% Tangled dependencies
    ProvMgr -.->|needs| EvtStream
    TaskReg -.->|needs| ProvMgr
    WSProvider -.->|needs| Indexer
    Srv -.->|needs| WSProvider
    Srv -.->|needs| Indexer
```

## 2. Target State (Runtime Coordinator)

By introducing the `internal/runtime` package, the `Runtime` becomes a thin coordinator. It does not wrap existing systems; instead, existing systems (like `tasks.Engine`) adapt to implement the `Component` interface.

```mermaid
graph TD
    Main[cmd/devserver/run.go] -->|bootstraps| Runtime[runtime.Coordinator]
    
    %% Components register themselves
    Runtime -->|manages lifecycle| EvtStream[core.EventStream : Component]
    Runtime -->|manages lifecycle| CapReg[capabilities.Registry : Component]
    Runtime -->|manages lifecycle| TaskEng[tasks.Engine : Component]
    Runtime -->|manages lifecycle| ProvMgr[providers.Manager : Component]
    Runtime -->|manages lifecycle| WSProvider[core.WorkspaceProvider : Component]
    
    %% Dynamic discovery
    TaskEng -.->|queries| Runtime
    Runtime -.->|returns| ProvMgr
```

### Key Differences
1. **No Wrapper Layers**: `TaskEngine` implements `Component`. We do not create a redundant `TaskRuntime` layer.
2. **Standardized Lifecycle**: The `Runtime` orchestrates `Initialize()`, `Start()`, `Health()`, and `Stop()` across all registered Components safely.
3. **Event Envelope**: `EventBus` utilizes standard envelopes (`ID`, `Type`, `Source`, `Time`, `Payload`) instead of generic maps.
