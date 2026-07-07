# Product Information Architecture — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Top-Level Product Structure (Global IA)

The DevServer platform is organized into seven core conceptual spaces that dictate the layout hierarchy, navigation groupings, and permission walls.

```
                  ┌─────────────────────────────────┐
                  │        DEVSERVER PLATFORM       │
                  └────────────────┬────────────────┘
                                   │
         ┌───────────┬─────────────┼────────────┬─────────────┐
         ▼           ▼             ▼            ▼             ▼
   [Dashboard]  [Workspace]  [Observe]   [Configure]   [Administration]
     (Overview)   (Develop)    (Monitor)   (Environment)  (Governance & Operations)
```

### The 7 Core Operational Phases:

1. **Dashboard (Overview)**: High-level entry point summarizing system health, task execution metrics, recent personal/organization activity timeline, and smart recommendations.
2. **Workspace (Sandbox)**: The developer's primary IDE space. Contains file tree explorer, code editor, services registry, task configurations, and the AI workspace chat assistant.
3. **Development (DevCenter)**: Documentation, database schemas, module dependency graphs, API reference specs, and the automated "Project Doctor" diagnostic tool.
4. **Deploy (Operations)**: Release pipeline tracking, environment promotions, rollback controls, release notes, and configuration diffs.
5. **Observe (Diagnostics)**: System telemetry (CPU, RAM, Disk, Network), active process monitoring, alerts dashboard, incident logs, and background job queues.
6. **Configure (Settings)**: Environment variables, AES-256 encrypted secrets management, and AI/LLM Provider settings (Ollama, OpenAI, connection states).
7. **Administration (Operations)**: Multi-tenant organization definitions, members listing, roles configuration (RBAC), security configurations (MFA, Active Sessions, API Keys), and data governance (Backup, Restore, Export, Import).

---

## 2. Navigation Hierarchy (Detailed Sitemap)

```
DevServer Global Shell
├── Topbar (Always Visible)
│   ├── Breadcrumbs (Home ➔ Workspace ➔ Section)
│   ├── Search (Ctrl+K Command Palette)
│   │   ├── Commands / Navigation Shortcuts
│   │   ├── Global search (Files, Members, Workspaces, Settings)
│   ├── Notification Popover (Inbox, Mentions, Failed tasks, Alerts)
│   ├── Help Center Popover (Keyboard shortcuts, Tutorials, Documentation, Troubleshooting)
│   └── User Menu Dropdown
│       ├── Profile Settings (MFA, Active Sessions, API keys, Preferences)
│       └── System Status Overview (Health status, Queue status, Active Executors)
│
├── L1 Sidebar (Main Navigation)
│   ├── OVERVIEW
│   │   ├── Dashboard                      /dashboard
│   │   ├── Projects                       /projects
│   │   └── Workspaces                     /workspace (Redirects to active workspace)
│   ├── OPERATIONS & OBSERVE
│   │   ├── Telemetry Monitor              /monitor (CPU/RAM timelines, Process lists, Incidents)
│   │   ├── Deployments                    /deploy (Releases, Rollbacks, Promotions, Env Diffs)
│   │   └── Backup & Data                  /backup (Backup, Restore, Import, Export logs)
│   ├── CONFIGURATION
│   │   ├── Environments                   /config/environments (Variables & Secrets)
│   │   └── AI Providers                   /config/providers (LLM setups, connection tests)
│   ├── DEVCENTER
│   │   ├── Workspace Doctor               /devcenter/project-doctor (Auto-healing tools)
│   │   ├── System Docs                    /devcenter/architecture (Schemas, Graph, APIs)
│   │   └── AI Cost & Usage                /devcenter/ai-usage (Token costs, model usage stats)
│   └── GOVERNANCE
│       └── Admin Settings                 /settings (Organizations, Members, RBAC Roles)
│
└── Active Workspace Context (S-014 Tabs)
    ├── Overview                           (Git status overview, project metadata)
    ├── Files                              (Split-pane tree explorer + Code editor)
    ├── Services                           (Database, Web server, Cache controller)
    ├── Tasks                              (Task Engine queues, run logs, console output)
    └── AI Workspace                       (Active conversation sessions, template library)
```

---

## 3. Global Systems & Features mapping

### 3.1 Global Search & Command Palette (Ctrl+K)
- **Scope**: Search works across Workspaces, Services, Files, Members, Providers, Logs, Tasks.
- **Actions**: Keyboard commands enable quick jump (e.g., `> Go to Settings`), quick workspace cloning (`> Clone Workspace`), or checking provider status.
- **Help Center Integration**: Quick list of keyboard shortcuts (`?`), link to troubleshooting guides, and online product documentation lookup.

### 3.2 Notification System
- **Inbox**: Holds critical alerts (e.g., "Deployment Failed", "Disk Space Low").
- **Mentions**: Notification if another developer tags you or assigns a task.
- **Background Jobs**: Live monitoring of indexing runs, git clone progress, database restore, or exports.

### 3.3 Activity Timeline System
- **Personal Activity**: Files changed, tasks run, or commits made by the active user.
- **Organization Activity**: Audit logs of member invites, billing, or provider changes.
- **Workspace Activity**: Task engine history, runtime service crash alerts, and active debug sessions.

### 3.4 Profile, MFA & Preferences Space
- **Security**: Profile image, password update, MFA setup, Active Sessions manager (with session revocation), and Personal API Keys generator.
- **Preferences**:
  - **Theme**: Dark mode, light mode, or system default.
  - **Shortcuts**: Customizable command palette bindings.
  - **Language & Notifications**: Email alerts or browser push toggles.
  - **Editor Config**: Font size, tab spacing, vim mode toggle.

### 3.5 System Status Space
- **Services Monitor**: Real-time health status of local adapters (Docker, Local, SSH, WSL, Kubernetes).
- **Queue status**: Task execution priority queue length and active worker count.
- **AI Providers status**: Network latencies, token usage counters, and cost estimates.

---

## 4. Operational Flows & Lifecycle States

### 4.1 Data Governance Lifecycle
```
[Export Configuration] ➔ [Config Zip File] ➔ [Import Configuration]
[Trigger Backup] ➔ [Backup Task running] ➔ [Archive created] ➔ [Trigger Restore] (Requires confirmation)
```
- **Safety checks**: Restores and Workspace Deletions block the main viewport and require typing the resource name to prevent accidental data loss.

### 4.2 Workspace Lifecycle
```
[Clone Workspace] ➔ [Cloning State] ➔ [Active Workspace]
                                           │
                        ┌──────────────────┴──────────────────┐
                        ▼                                     ▼
               [Archive Workspace]                    [Delete Workspace]
             (Read-only, preserves data)            (Destructive file deletion)
```

### 4.3 Environment & Provider Lifecycle
- **Duplicate Environment**: Creates copy of variables list without decrypting secrets (creates fresh secret refs).
- **Duplicate Provider**: Quick copy-paste of LLM API configurations for staging/production tests.

### 4.4 AI Development Context
- **Conversations history**: Stores past user/agent sessions.
- **Prompt Templates**: Shared library of developer prompts.
- **Token Costs**: Cost tracker mapping prompt tokens, completion tokens, and dollar expenditures.

### 4.5 Monitoring & Observing
- **Resource graphs**: Live charts updating every 5s.
- **Incident lifecycle**: `Warning` ➔ `Degraded` ➔ `Alerting` ➔ `Resolved`. Includes a timeline detailing process memory leaks or CPU spikes.

### 4.6 Deployment Pipeline
- **Release notes generator**: Auto-creates release draft by comparing Git commits between target environments.
- **Promotion flow**: `Staging` ➔ `Production` promotions. Compares variables/secrets and warns of discrepancies.
- **Rollback execution**: Single-click reversion to the previous stable release.

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
