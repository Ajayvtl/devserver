# State Transition Models — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Workspace Lifecycle State Machine

A workspace manages its runtime state based on cloning operations, indexing tasks, and lifecycle actions.

```mermaid
stateDiagram-v2
    [*] --> Creating
    Creating --> Provisioning : DB entry registered
    Provisioning --> Indexing : Repository cloned
    Indexing --> Ready : File index complete
    
    Ready --> Busy : Task triggered
    Busy --> Ready : Task complete
    
    Ready --> Offline : Network / Heartbeat lost
    Offline --> Ready : Link restored
    
    Ready --> Archived : Archive triggered
    Archived --> Ready : Restore active
    
    Archived --> Deleted : Delete confirmed
    Ready --> Deleted : Delete confirmed
    Deleted --> [*]
```

### Transition Validation Rules

| Current State | Event | Next State | Triggering API Action | Permission |
|---|---|---|---|---|
| **Creating** | Clone repo success | **Provisioning** | Git clone finished | System |
| **Provisioning**| Index files start | **Indexing** | Indexer spawned | System |
| **Indexing** | Index finished | **Ready** | Index files complete | System |
| **Ready** | Run Task | **Busy** | `POST /api/tasks` | Developer+ |
| **Busy** | Task finished | **Ready** | Executor task end | System |
| **Ready** | Connection lost | **Offline** | Heartbeat failure (30s timeout) | System |
| **Offline** | Heartbeat received| **Ready** | Network link restored | System |
| **Ready** | Archive workspace | **Archived** | `POST /api/workspaces/archive` | Operator+ |
| **Archived** | Restore workspace | **Ready** | `POST /api/workspaces/unarchive`| Developer+ |
| **Ready** | Delete workspace | **Deleted** | `DELETE /api/workspaces` | Admin+ |

---

## 2. Task Engine State Machine

Tasks are managed by the execution runner queue and report real-time progress via WebSockets.

```mermaid
stateDiagram-v2
    [*] --> Queued
    Queued --> Running : Worker allocated
    Running --> Success : Process exits with 0
    Running --> Failed : Process exits with non-0
    Running --> Cancelled : User terminates task
    
    Success --> [*]
    Failed --> [*]
    Cancelled --> [*]
```

### Transition Validation Rules

| Current State | Event | Next State | Triggering Event |
|---|---|---|---|
| **Queued** | Executor slots clear | **Running** | Runner consumes queue item |
| **Running** | Task process finishes | **Success** | Exit code `0` returned |
| **Running** | Task process crashes | **Failed** | Exit code `> 0` returned |
| **Running** | Terminate clicked | **Cancelled** | SIGINT/SIGKILL sent to executor |

---

## 3. Deployment Pipeline State Machine

Deployments track release states targeted at specific environments.

```mermaid
stateDiagram-v2
    [*] --> Pending
    Pending --> Building : Git fetch complete
    Building --> Testing : Compile success
    Testing --> Deploying : Tests pass
    Deploying --> Active : Health checks pass
    
    Building --> Failed : Compile error
    Testing --> Failed : Test failures
    Deploying --> RolledBack : Health check fail / Manual Rollback
    
    Active --> RolledBack : Manual Rollback triggered
    Active --> Deprecated : Shipped new active release
    
    Failed --> [*]
    RolledBack --> [*]
```

---

## 4. Member Invitation State Machine

Governs team membership lifecycle states within an organization.

```mermaid
stateDiagram-v2
    [*] --> Invited
    Invited --> Active : Accepts invite & logs in
    Active --> Suspended : Admin disables account
    Suspended --> Active : Admin re-enables account
    Active --> Removed : Delete member triggered
    Invited --> Expired : 7-day invite timeout reached
    Expired --> Invited : Admin resends invitation
    Removed --> [*]
```

---

## 5. Infrastructure Executor State Machine

Represents connection states of remote docker, WSL, local, or SSH runner environments.

```mermaid
stateDiagram-v2
    [*] --> Disconnected
    Disconnected --> Connecting : Connection test triggered
    Connecting --> Online : SSH/API handshake success
    Connecting --> Offline : Handshake timeout / Network unreachable
    Online --> Degraded : Memory/Disk > 90% or CPU throttled
    Online --> Offline : Connection lost
    Offline --> Connecting : Auto-retry ping
```

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
