# Onboarding, Installation & Mobile Specifications

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Mobile Navigation & Sitemap Hierarchy

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

## 2. Platform Installation & Onboarding Experience (Gap #3 / Refinement)

The installation setup flow runs in an isolated step-by-step progress UI.

### 2.1 Installation Progress UI
The bootstrap interface S-003 features a central panel with an active step indicator, a detailed task logs console, and an overall percentage progress bar.

```
┌────────────────────────────────────────────────────────┐
│  DevServer Platform Setup                              │
├────────────────────────────────────────────────────────┤
│  [ Step 2 of 4: Provisioning Database ]                 │
│  Progress: [████████████████░░░░░░░░░] 60%             │
├────────────────────────────────────────────────────────┤
│  Task Log Output:                                      │
│  > Checking port 5432... Open.                         │
│  > Initializing system tables... Success.              │
│  > Running migrations [14 of 22]... Fail.              │
├────────────────────────────────────────────────────────┤
│  [ Retry Step ]   [ View Diagnostic Logs ]             │
└────────────────────────────────────────────────────────┘
```

### 2.2 Partial Failure Recovery & Retry/Resume Pathways
To prevent restarting the installation from step 1 when a network drop or database timeout occurs:
- **Save checkpoints**: Each successful setup phase writes its success state to local disk configurations.
- **Fail states**: If an execution step fails, the progress bar turns red, the task logs display the error output, and a "Retry Step" button is enabled.
- **Diagnostics modal**: Provides access to copy system error trace logs for troubleshooting.
- **Resume logic**: Reloading the browser resumes setup from the last uncompleted checkpoint.

### 2.3 Post-Install Checklist & Health Verification
Before redirecting the user to the login screen, the installer runs a suite of health verification tests:
1. **Database link check**: Runs database latency pings.
2. **Indexer engine test**: Launches a mock files indexer test.
3. **Local Docker socket ping**: Validates executor socket permissions.
4. **Network connectivity test**: Pings standard package index URLs.

#### Health Output Checklist Screen:
```
[✔] Database Connection: Online (Latency 4ms)
[✔] File System Indexer: Verified
[✔] Docker Runtime adapter: Connected (Docker Engine v24.0)
[✔] LLM Provider Connectivity: Connected
```

---

## 3. Empty Organization Onboarding Wizard

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
