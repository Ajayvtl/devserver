# DevServer Product Specification

This canonical specification replaces fragmented UX and Menu structure documentation. It ensures the UI, Backend, and Architecture remain permanently aligned around the Environment Model.

## The Environment Model Paradigm

DevServer shifts the industry standard from Workspace-centric to **Environment-centric** orchestration.

Instead of `Workspace -> Services`, the hierarchy is strictly:
`Project -> Environment -> Resources (Services/Deployments/Databases)`

A user selects a **Global Context Switcher** (Local / Staging / Production). All operational tabs (Services, Databases, Domains) instantly re-bind their data sources and capabilities to that environment's Execution Context.

---

## Canonical Page Definitions

### 1. Overview Page (`/workspace/[id]`)
* **Purpose**: High-level observability for the selected environment.
* **Actors**: Developer, DevOps
* **Data Sources**: `EnvironmentInfo`, `HealthScore`
* **Runtime Dependencies**: Active Execution Context (Local/SSH)
* **Actions**: `Workflow: Boot Environment`
* **Resource Dependencies**: Workspace
* **Future Scope**: Multi-environment split view.

### 2. Services Page
* **Purpose**: Manage infrastructure processes within the active environment.
* **Data Sources**: `Resource[]`
* **Runtime Dependencies**: `ExecutionContext` (Local, Docker, SSH)
* **Actions**: `Operation: Start, Stop, Restart, Tail Logs`
* **Command Bus Events**: `resource.start`, `resource.stop`
* **Resource Dependencies**: Redis, MySQL, Node, Docker
* **Failure Modes**: Missing runtime dependencies, port conflicts
* **Offline Behaviour**: Disables controls, retains last known state.

### 3. Databases Page
* **Purpose**: Provision and manage data stores.
* **Data Sources**: `Resource[]`
* **Actions**: `Operation: Provision, Connect, Dump, Restore`
* **Command Bus Events**: `database.provision`
* **Resource Dependencies**: Postgres, MySQL, Redis
* **Future Scope**: Automatic `.env.local` injection for connection strings.

### 4. Deployments Page
* **Purpose**: Manage multi-environment synchronization.
* **Data Sources**: SQLite `deployments` table, Git branch history.
* **Actions**: `Workflow: Sync to Environment`, `Workflow: Rollback`
* **Command Bus Events**: `environment.sync`
* **Future Scope**: Real-time streaming of SSH output during remote sync.

## The Synchronization Triangle
DevServer orchestrates state across three planes:
1. **Local Space**: Development context (hot reloading, local DBs).
2. **Middle Space**: Source of truth (GitHub/GitLab, registry).
3. **Remote Space**: Deployed execution (Staging/Production via SSH/K8s).

Workflows bridge these spaces (e.g., "Deploy to Staging" = Git Push -> SSH -> Pull -> Build -> Restart -> Health Check -> Rollback if failed).

---

## Strategic Execution Order

To implement this abstraction without causing regressions, the rollout is phased:

* **Phase 1 (Architecture Freeze)**: Finalize the canonical object model (`Project`, `Workspace`, `Environment`, `ExecutionContext`, `Resource`, `Operation`, `Workflow`).
* **Phase 2 (Runtime)**: Implement `Environment`, `ExecutionContext`, and routing layer.
* **Phase 3 (Resources)**: Convert existing providers (Redis, Node, Docker, MySQL, Nginx) strictly to the `Resource` model.
* **Phase 4 (Workflow Engine)**: Introduce orchestration (`Boot Development`, `Deploy Staging`, `Backup Database`).
* **Phase 5 (UI)**: Add the Global Context Switcher and bind all pages to the active Environment.
* **Phase 6 (Future Subsystems)**: Only after Phase 5 is stable, proceed with Logs, Domains, Monitoring, AI, Editor, and Terminal.
