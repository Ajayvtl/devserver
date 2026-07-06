# Environment Runtime Architecture

This document establishes the fundamental abstraction for the DevServer platform: The Environment Model.

## The Core Problem
Historically, DevServer mapped Runtimes and Services directly to a `Workspace`, mixing Execution, Technology, Application, and Infrastructure into a single overloaded `Provider` concept (e.g., `RedisProvider`, `DockerRedisProvider`, `SystemdRedisProvider`). This causes a combinatorial explosion of providers.

## The Solution: The Resource & Context Model
The architecture strictly separates WHAT something is (Resource) from HOW it executes (Context).

The hierarchy is strictly:

`Project -> Workspace -> Environment -> Execution Context -> Resource -> Operation`

### Definitions

* **Project**: The logical grouping of code and intent (e.g., "Aurora Platform").
* **Workspace**: The physical instantiation of a Project on a host (e.g., `/Users/dev/aurora`).
* **Environment**: A specific boundary of execution (e.g., `Local`, `Development`, `QA`, `Staging`, `Production`). Each Environment owns its own Resources.
* **Execution Context**: The mechanism for executing commands within an Environment (e.g., `Local`, `SSH`, `Docker`, `WSL`, `Kubernetes`).
* **Resource**: The definition of an application or infrastructure dependency (e.g., `Redis`, `Node`, `MySQL`, `Git`). The Resource only knows *what* it is.
* **Operation**: An intent or capability supported by a Resource (e.g., `Start`, `Stop`, `Restart`, `Install`, `Detect`, `Backup`).
* **Workflow**: An orchestration layer that sequences Operations across multiple Resources and Contexts (e.g., "Deploy Production" -> Git Push -> SSH Pull -> Build -> Restart -> Health Check -> Rollback).

## Architectural Constraints

1. **Execution and Resource Must Never Mix**. A Resource (`Redis`) must never contain OS-specific logic (e.g., `systemctl`). It defines Operations, and the `ExecutionContext` performs them.
2. **No Overloaded Providers**. We do not build `SSHRedisProvider`. We build a `Redis` Resource, which requests a `Start` Operation, dispatched through an `SSH` Execution Context.
3. **The UI Binds to Environments**. Operations pages (Services, Databases, Domains) derive their state from the currently selected Environment, not the Workspace.

## Workflow Example: Starting Redis

**Current Anti-Pattern:**
`Workspace -> RedisProvider.Start()`

**New Abstraction Pattern:**
1. User selects `[Environment: Staging]` in the Global Context Switcher.
2. User triggers `Start` on the `Redis` Resource.
3. The UI dispatches `Operation: Start(Resource: Redis)`.
4. Command Bus resolves the Staging Environment, retrieving its bound `ExecutionContext` (e.g., `SSHRuntime`).
5. `SSHRuntime` translates the generic `Start` intent into `systemctl start redis` via SSH.
