# Onboarding, Installation & Mobile Specifications

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Mobile Navigation & Sitemap Hierarchy (Gap #10)

On viewports below `720px` width, the standard L1 sidebar is replaced by a bottom tab navigation menu and a topbar hamburger drawer.

### 1.1 Bottom Navigation Bar (Primary Actions)
Contains 4 high-frequency icons:
1. **Home**: `/dashboard` (S-005 dashboard metrics)
2. **Projects**: `/projects` (Workspace listings)
3. **Deployments**: `/deploy` (Release board S-010)
4. **Profile**: `/profile` (Session settings S-013)

### 1.2 Hamburger Drawer Menu (Secondary Pages)
Tapping the topbar menu button slides in a left-aligned drawer containing:
* **Governance & Settings**:
  * Member Management (`/settings#members`)
  * Organizations Switcher (`/settings#orgs`)
  * Backup & Restore Logs (`/backup`)
* **DevCenter Spec Docs**:
  * Architecture Graph (`/devcenter/architecture`)
  * Database Schema (`/devcenter/database`)
  * System Doctor (`/devcenter/project-doctor`)
* **Configuration Options**:
  * Environments (`/config/environments`)
  * AI Providers (`/config/providers`)

---

## 2. Production Installation & Onboarding Journey (Gap #11)

This flow maps out the lifecycle from bare server setup to workspace ready-state.

```
[System Check] ➔ [DB Initialization] ➔ [Create Root Admin] ➔ [Create Org] ➔ [AI Config] ➔ [Var Setup] ➔ [Clone Repo] ➔ [IDE Ready]
```

1. **Bare Install & System Check**: Verify server environment compatibility (WSL, Docker, SSH ports).
2. **Database Provisioning**: Install admin system schemas, schemas check pass.
3. **Register Root Account**: Create first administrative user profile (S-003 wizard).
4. **Organization Allocation**: Spawn default organizational tenant.
5. **AI Provider Connection**: Input endpoint address for Ollama, OpenAI, or other model engines.
6. **Environment Variable Configuration**: Define dev-stage variables and inject encrypted secrets.
7. **Allocate Workspace & Clone Repository**: Input Git URL, clone repositories to runner storage, and run index scan.
8. **IDE Ready**: Open Workspace IDE panel (S-024).

---

## 3. Empty Organization Onboarding Wizard (Gap #12)

If a user authenticates and enters an organization with:
* `0 projects`
* `0 members`
* `0 providers`
* `0 workspaces`

The normal dashboard is hidden, and the system runs the **Organization Onboarding Wizard (W-004)**.

### Onboarding Steps

```
┌────────────────────────────────────────────────────────────────────────┐
│  Onboarding Wizard: Setup Organization                                  │
├────────────────────────────────────────────────────────────────────────┤
│  Step 1: Invite Team Members (Email input + Role dropdown)             │
│  Step 2: Configure AI Provider (Ollama or OpenAI API connection)        │
│  Step 3: Define Development Environment (Variables & Secrets setup)    │
│  Step 4: Create First Workspace (Input git repository URL to clone)    │
└────────────────────────────────────────────────────────────────────────┘
```

* **Step 1: Invite Team Members**: Text area to invite colleagues via email and assign access roles.
* **Step 2: Configure AI Provider**: Select Ollama (default local endpoint) or OpenAI, input API keys, and test connection.
* **Step 3: Create Development Environment**: Pre-configure variables (e.g., `PORT=8080`, `DB_USER=root`) and secrets reference keys.
* **Step 4: Create Workspace**: Git repository link input. Tapping "Finish & Clone" starts cloning in the background and opens the Workspace Explorer IDE view.

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
