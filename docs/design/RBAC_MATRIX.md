# RBAC & Role Navigation Matrix — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Role Permission Matrix

The table below maps every role to allowed actions across all primary platform screens. 
* **Yes**: Action is allowed and visible.
* **No**: Action is blocked (API returns 403) and UI element is hidden/disabled.

| Screen / Feature Area | Action | Owner | Admin | Developer | Operator | Viewer |
|---|---|---|---|---|---|---|
| **Setup Wizard (`/setup`)** | Initialize Platform | Yes | No | No | No | No |
| **Organizations** | Create Organization | Yes | No | No | No | No |
| | Delete Organization | Yes | No | No | No | No |
| | Rename Organization | Yes | Yes | No | No | No |
| **Members & Roles** | Invite Member | Yes | Yes | No | No | No |
| | Edit Member Role | Yes | Yes | No | No | No |
| | Remove Member | Yes | Yes | No | No | No |
| | Transfer Ownership | Yes | No | No | No | No |
| **AI Providers** | Add AI Provider | Yes | Yes | No | No | No |
| | Edit AI Provider | Yes | Yes | No | No | No |
| | Delete AI Provider | Yes | Yes | No | No | No |
| | Test Connection | Yes | Yes | No | No | No |
| **Environments & Secrets**| Create Environment | Yes | Yes | No | No | No |
| | Delete Environment | Yes | Yes | No | No | No |
| | Set Env Variable | Yes | Yes | Yes | No | No |
| | View Secrets (Refs) | Yes | Yes | Yes | Yes | Yes |
| | Decrypt Secret (API) | Yes | Yes | Yes | No | No |
| **Projects & Workspaces** | Create Project / Workspace| Yes | Yes | Yes | No | No |
| | Clone Workspace | Yes | Yes | Yes | No | No |
| | Archive Workspace | Yes | Yes | Yes | Yes | No |
| | Delete Workspace | Yes | Yes | No | No | No |
| | Edit Code in Workspace | Yes | Yes | Yes | No | No |
| **Tasks & Run Commands**  | Trigger Task / Action | Yes | Yes | Yes | Yes | No |
| | Configure / Edit Task | Yes | Yes | Yes | No | No |
| | Terminate Running Task | Yes | Yes | Yes | Yes | No |
| **Deployments & Pipelines**| Trigger Release | Yes | Yes | Yes | Yes | No |
| | Rollback Release | Yes | Yes | No | Yes | No |
| | Promote Staging ➔ Prod | Yes | Yes | No | Yes | No |
| **Diagnostics & Telemetry**| View CPU/RAM Charts | Yes | Yes | Yes | Yes | Yes |
| | View Process List | Yes | Yes | Yes | Yes | Yes |
| | Restart Degraded Process| Yes | Yes | No | Yes | No |
| **Backup & Data Restore** | Trigger Backup Job | Yes | Yes | No | Yes | No |
| | Trigger Restore Action | Yes | No | No | No | No |
| | Export Configuration | Yes | Yes | No | No | No |
| | Import Configuration | Yes | No | No | No | No |
| **Security & Auditing**   | View Audit Logs | Yes | Yes | No | No | No |
| | Revoke User Sessions | Yes | Yes | No | No | No |
| | Generate System API Keys| Yes | Yes | No | No | No |

---

## 2. Navigation Visibility by Role

Developers should not see menus they cannot act on. Below are the navigation structures displayed for each role.

### 2.1 Owner Navigation Map
```
Topbar: [Breadcrumbs] ➔ [Search Ctrl+K] ➔ [Notifications Inbox] ➔ [Help Center] ➔ [User Profile Dropdown]
Sidebar:
 ├── OVERVIEW
 │    ├── Dashboard
 │    ├── Projects
 │    └── Active Workspace
 ├── OPERATIONS & OBSERVE
 │    ├── Telemetry Monitor
 │    ├── Deployments Board
 │    └── Backup & Data Hub
 ├── CONFIGURATION
 │    ├── Environments
 │    └── AI Providers
 ├── DEVCENTER
 │    ├── Workspace Doctor
 │    ├── System Docs
 │    └── AI Cost & Token Usage
 └── GOVERNANCE
      └── Admin Settings (Full Orgs, Members, Audit Logs, API keys, Session managers)
```

### 2.2 Admin Navigation Map
```
Topbar: [Breadcrumbs] ➔ [Search Ctrl+K] ➔ [Notifications Inbox] ➔ [Help Center] ➔ [User Profile Dropdown]
Sidebar:
 ├── OVERVIEW
 │    ├── Dashboard
 │    ├── Projects
 │    └── Active Workspace
 ├── OPERATIONS & OBSERVE
 │    ├── Telemetry Monitor
 │    ├── Deployments Board
 │    └── Backup & Data Hub
 ├── CONFIGURATION
 │    ├── Environments
 │    └── AI Providers
 ├── DEVCENTER
 │    ├── Workspace Doctor
 │    ├── System Docs
 │    └── AI Cost & Token Usage
 └── GOVERNANCE
      └── Admin Settings (Members, Audit Logs, API keys; *NO Organization deletion/restore*)
```

### 2.3 Developer Navigation Map
```
Topbar: [Breadcrumbs] ➔ [Search Ctrl+K] ➔ [Notifications Inbox] ➔ [Help Center] ➔ [User Profile Dropdown]
Sidebar:
 ├── OVERVIEW
 │    ├── Dashboard (Dev-scoped metrics)
 │    ├── Projects
 │    └── Active Workspace (Full IDE View)
 ├── OPERATIONS & OBSERVE
 │    └── Telemetry Monitor (Read-only System Metrics)
 ├── CONFIGURATION
 │    └── Environments (Read-only variables; *No Provider modifications*)
 └── DEVCENTER
      ├── Workspace Doctor
      └── System Docs (Architecture, Database Schema, Dependency Graph)
```

### 2.4 Operator Navigation Map
```
Topbar: [Breadcrumbs] ➔ [Search Ctrl+K] ➔ [Notifications Inbox] ➔ [Help Center] ➔ [User Profile Dropdown]
Sidebar:
 ├── OVERVIEW
 │    ├── Dashboard (Status metrics)
 │    ├── Projects
 │    └── Active Workspace (Log-only viewer)
 ├── OPERATIONS & OBSERVE
 │    ├── Telemetry Monitor (Full charts & Process list controls)
 │    ├── Deployments Board (Rollbacks & Promotions)
 │    └── Backup & Data Hub (Triggering backups; *NO Restore*)
 └── DEVCENTER
      └── System Docs (Architecture, Database Schema)
```

### 2.5 Viewer Navigation Map
```
Topbar: [Breadcrumbs] ➔ [Search Ctrl+K] ➔ [Notifications Inbox] ➔ [Help Center] ➔ [User Profile Dropdown]
Sidebar:
 ├── OVERVIEW
 │    ├── Dashboard (Read-only status)
 │    ├── Projects (Read-only browse)
 │    └── Active Workspace (Read-only logs / files)
 └── DEVCENTER
      └── System Docs (Architecture, Database Schema)
```

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
