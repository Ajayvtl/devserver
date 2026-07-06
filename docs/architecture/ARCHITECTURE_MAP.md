# System Architecture Map

This document defines the structural layers, boundaries, and ownership within DevServer.

## Architectural Layers

| Layer | Packages | Ownership | Dependencies Allowed |
| --- | --- | --- | --- |
| **App / Presentation** | `apps/web`, `cmd/devserver` | UI, CLI, routing, bootstrapping | Application |
| **Application** | `internal/application`, `internal/application/ports`, `internal/application/workflow` | Use Cases, Orchestration, Workflow Engine, Ports | Domain, Repository Interfaces, Executor Interfaces |
| **Domain** | `internal/domain` | Pure Business Rules, Entities, Value Objects, Identifiers | None |
| **Repository (Interfaces)** | `internal/repository` | Persistence Contracts | Domain |
| **Executor (Interfaces)** | `internal/executor/contracts` | Execution Contracts | Domain |
| **Infrastructure** | `internal/providers`, `internal/executor/local`, `internal/executor/ssh`, `internal/executor/legacy` | SQL implementations, APIs, OS commands | Repository, Executor Interfaces, Domain |

## Governance Rules
1. **Dependency Direction**: Outer layers depend on inner layers.
2. **Ports and Adapters**: Infrastructure implements interfaces (Ports) defined by the Application or Domain layers.
3. **No Cross-Pollination**: Infrastructure packages should never import each other directly.
