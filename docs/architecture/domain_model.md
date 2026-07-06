# DevServer Canonical Domain Model

This document is the single source of truth for every core entity in DevServer. It strictly defines the domain objects, ownership relationships, lifecycles, and persistence boundaries before any interface code is written.

---

## 1. Core Principles
1. **Strict Hierarchy**: Objects belong to a strict ownership tree (`Organization -> Project -> Workspace -> Environment -> Resource`).
2. **Separation of Concerns**: *What* something is (Resource) is strictly decoupled from *How* it executes (ExecutionContext).
3. **Capability-Driven**: Runtimes and Actions rely on capability flags (e.g., `SupportsDocker`) rather than brittle conditionals (`if docker { ... }`).
4. **Data-Driven Orchestration**: Complex actions are represented as Directed Acyclic Graphs (DAGs) of Actions inside a Workflow.

---

## 2. Domain Objects

### 2.1 Project
* **Definition**: The highest-level logical grouping of code and intent (e.g., "Aurora Platform").
* **Ownership**: Owned by Organization (future). Owns Workspaces.
* **Responsibilities**: Stable identity, global team settings, and project-wide metadata.
* **Lifecycle**: Created once per codebase; persists indefinitely until manually deleted.
* **Serialization & Persistence**: Stored in DevServer's central SQLite database.
* **Dependencies**: Independent.

### 2.2 Workspace
* **Definition**: The physical instantiation of a Project on a host machine (e.g., `/Users/dev/aurora`).
* **Ownership**: Owned by Project. Owns Environments.
* **Responsibilities**: Mapping a logical project to physical disk paths, file indexing, and local machine state.
* **Lifecycle**: Created when a project is cloned/opened; destroyed if the directory is removed or unregistered.
* **Serialization & Persistence**: Stored in DevServer's central SQLite database (`workspaces` table).

### 2.3 Environment
* **Definition**: A specific boundary of execution context (e.g., `Local`, `Development`, `QA`, `Staging`, `Production`).
* **Ownership**: Owned by Workspace. Owns Resources, Variables, Secrets, Policies, Connections, Workflows, Templates, and Schedules.
* **Responsibilities**: Total isolation of configuration and execution state. A Redis resource in `Local` is entirely distinct from Redis in `Production`.
* **Lifecycle**: Created implicitly (Local) or explicitly (Staging/Prod); persists as long as the environment is actively managed.
* **Serialization & Persistence**: Definitions stored in SQLite or a shared `.devserver/environments.json` file to allow team syncing.

### 2.4 Executor (formerly Execution Context)
* **Definition**: The capability-driven mechanism for executing commands (e.g., `Local`, `SSH`, `Docker`, `Kubernetes`).
* **Ownership**: Owned by Environment.
* **Responsibilities**: Abstracting OS/Network boundaries.
* **Capabilities**: Defines data-driven capabilities (e.g., Name, Version, Constraints).
* **Lifecycle**: Instantiated dynamically at runtime when an Environment is selected.
* **Serialization & Persistence**: State lives only in runtime memory. Configuration (e.g., SSH Host/Port) is stored as a Connection.

### 2.5 Template
* **Definition**: A blueprint for creating Resources or Workflows (e.g., `Node App Template`, `Deployment Workflow Template`).
* **Ownership**: Owned by Project or global registry.
* **Responsibilities**: Providing baseline defaults and schemas.
* **Serialization & Persistence**: Stored in config/SQLite.

### 2.6 ResourceSpec (Desired State)
* **Definition**: An application or infrastructure dependency instantiated from a Template (e.g., `Redis Instance`, `Node App`).
* **Ownership**: Owned by Environment.
* **Properties**: `ID`, `Type`, `Name`, `Owner`, `DesiredConfiguration`.
* **Responsibilities**: Representing the declared desired state. 
* **Serialization & Persistence**: Metadata stored in SQLite/config.

### 2.7 RuntimeState (Observed State)
* **Definition**: The observed live state of a Resource.
* **Ownership**: Ephemeral, tied to the Executor mapping of a Resource.
* **Properties**: `Status`, `Metrics`, `ObservedAt`.
* **Serialization & Persistence**: Never persisted to long-term storage; lives in memory or ephemeral caches.

### 2.8 Action
* **Definition**: A formalized intent executed against a target (e.g., `Start`, `Stop`, `Build`).
* **Ownership**: Standalone or belonging to a Workflow/Resource.
* **Properties**: `ID`, `Type`, `Target`, `Arguments`. (Note: Retry/Rollback/Timeout belong to WorkflowNode).
* **Responsibilities**: Encapsulating exactly *what* needs to be done.
* **Serialization & Persistence**: Execution history stored in SQLite `tasks` table.

### 2.9 Event
* **Definition**: A first-class record of an occurrence within the platform (e.g., `Node Started`, `Health Check Failed`).
* **Ownership**: Produced by Executors, Actions, or Resources.
* **Properties**: `ID`, `Type`, `Source`, `Payload`, `OccurredAt`.
* **Responsibilities**: Enabling asynchronous, data-driven orchestration and tracking.
* **Serialization & Persistence**: Stored in SQLite for audit and replay.

### 2.10 Workflow & WorkflowNode
* **Definition**: A Directed Acyclic Graph (DAG) orchestration of Actions.
* **Ownership**: Owned by Environment.
* **Properties**: Nodes (Actions wrapped with `Retry`, `Timeout`, `Rollback`, `Condition`), Edges (Dependencies).
* **Triggers**: Reserved for `Manual`, `Schedule`, `Webhook`, `Event`, `API`.
* **Responsibilities**: Handling complex multi-step processes including parallelism and conditions.
* **Serialization & Persistence**: Templates stored in config. Execution state stored in SQLite.

### 2.10 Connection
* **Definition**: Authorized integrations with external systems, classified into `Transport` (SSH, HTTP) and `Integration` (GitHub, AWS).
* **Ownership**: Owned by Environment (or globally by Project).
* **Responsibilities**: Providing authenticated clients.
* **Serialization & Persistence**: Metadata stored in SQLite.

### 2.11 Secret
* **Definition**: Secure reference to sensitive data.
* **Ownership**: Owned by Environment.
* **Properties**: `Reference`, `Provider` (Vault, AWS Secrets, SQLite), `RotationPolicy`. (Values are NEVER stored in the domain).
* **Responsibilities**: Safely retrieving credentials at runtime.
* **Serialization & Persistence**: Metadata stored in SQLite.

### 2.12 Variable
* **Definition**: Environment-specific plaintext configuration (e.g., `NODE_ENV=development`, `PORT=3000`).
* **Ownership**: Owned by Environment.
* **Serialization & Persistence**: Stored in SQLite or `.env` files.

### 2.13 Policy
* **Definition**: Rules governing allowed Actions (e.g., "Cannot Stop Production Database without confirmation").
* **Ownership**: Owned by Environment.
* **Responsibilities**: Preventing destructive actions and enforcing guardrails.
* **Serialization & Persistence**: Stored in config/SQLite.

### 2.14 Dependency
* **Definition**: Generic DAG edges between entities (e.g., `From`, `To`, `Type`, `Condition`). Types include `Soft`, `Hard`, `Runtime`, `Build`.
* **Ownership**: First-class entity defining relationships.
* **Serialization & Persistence**: Stored as relational mappings in SQLite.

---

## 3. State Transitions

### Action Lifecycle
1. **Pending**: Action is queued in the DAG.
2. **Executing**: ExecutionContext is actively processing the Action.
3. **Evaluating**: Policy or Health Checks are running post-execution.
4. **Success**: Action completed.
5. **Failed**: Action threw an error or timed out.
6. **Rolling Back**: A rollback Action was triggered due to failure in the DAG.

## 4. Dependency Rules (Who can reference whom)
- **UI** depends strictly on **Environment** and its **Workflows/Resources**.
- **Resources** depend on **Variables**, **Secrets**, and other **Resources**.
- **Workflows** depend on **Actions** and **Policies**.
- **Actions** depend on **ExecutionContexts** to perform work.
- **ExecutionContexts** depend on **Connections** (e.g., SSH keys) to establish transport.
- **Providers (Implementations)** depend on **Actions** and **ExecutionContexts**.
