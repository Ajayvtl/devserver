# User Journeys — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Actor Role Definitions

* **Owner**: Ultimate data and org administrator. Authority for billing, deletions, backups, and ownership handover.
* **Admin**: Governance controller. Configures environments, secrets, providers, roles, and members.
* **Developer**: Day-to-day workspace code editor. Configures tasks, branches, coding, and debugging.
* **Operator**: Release manager. Observes system health dashboards, deploys releases, triggers rollbacks, and manages database backups.
* **Viewer**: Read-only stakeholder. Inspects project configurations, check run metrics, and observes logs.

---

## 2. Master End-to-End Workflow Journey (Gap #2 / Refinement)

### Journey 2.0: Project Creation ➔ Deployment ➔ Monitor & Rollback (Developer & Operator)
This master workflow maps the full lifecycle of a development cycle.

```
[Create Project] ➔ [Provision Workspace] ➔ [Clone Repo] ➔ [AI Indexing] ➔ [Run Tasks] ➔ [Deploy Release] ➔ [Telemetry Check] ➔ [Trigger Rollback]
```

* **Step 1: Project Initialization**:
  - Developer navigates to Projects S-007 and clicks "New Project".
  - Fills out project name and repository origin SSH URL, selecting default environment configurations.
* **Step 2: Workspace Provisioning**:
  - The system creates a new workspace database record, maps target environment variables, and allocates a default infrastructure Executor (e.g., Docker container).
* **Step 3: Background Git Clone**:
  - Progress indicator runs. The system triggers `git clone` on the target Executor.
  - Console logs stream file extraction progress.
* **Step 4: AI Indexing Loop**:
  - Once cloned, the Workspace file indexer runs in the background.
  - Code block index markers are mapped, parsing symbols, files, and directory dependencies. The AI chat panel (S-024e) displays "Indexing complete; model is ready to support coding".
* **Step 5: Code Editing & Task Execution**:
  - Developer opens `app.go` in the workspace editor (S-026), writes code edits, and saves files.
  - Launches the Task console (S-024d) and runs `npm run test` or `go test`. The system executes the task on the target executor and streams output.
* **Step 6: Promotion & Deployment**:
  - Tests pass. Developer clicks "Release Version" from the Deploy tab.
  - Inputs release tag version `v1.1.0`. The build engine creates a package and pushes it to the Staging target environment.
* **Step 7: Observability & Telemetry Checks**:
  - The Operator logs in and opens Telemetry Monitor S-009.
  - Observes CPU/RAM graphs. Suddenly, memory usage spikes up to 96% and the CPU chart starts flashing red "Degraded Performance Alert".
* **Step 8: Rollback Mitigation**:
  - Operator goes to Deployments board S-010, selects the previous release `v1.0.9`, and clicks "Execute Rollback".
  - Confirms the warning modal. The deploy engine reverts configuration states, rebuilds, and swaps the staging container back to the previous stable release. Health score returns to 98% (green).

---

## 3. Core Operational Journeys

### Journey 3.1: First Run & Platform Boot (Owner)
* **Goal**: Install and initialize the platform from a bare system setup.
* **Entry Point**: Root page `/` redirects automatically to `/setup` if database is uninitialized.
* **Flow**:
  1. System checks screen (W-001 step 1): Validates network, disk space, and runtime adapters.
  2. Database Provisioning (W-001 step 2): Setup admin schema connections.
  3. Admin Account Setup (W-001 step 3): Input email, password, and basic organization details.
  4. Complete (W-001 step 4): Launches initial login page S-004.
* **Recovery Actions**:
  - Check failure: Present diagnostic page detailing which adapter failed, with diagnostic codes.
  - Setup timeout: Clear cache and auto-trigger setup test again.

### Journey 3.2: Empty Organization Setup (Admin)
* **Goal**: Establish working workspaces, environments, and team access in a newly created empty org.
* **Entry Point**: Redirected to S-012 (Projects List) showing E-007 empty project state.
* **Flow**:
  1. Click "Invite Team" link in empty state.
  2. Dialog M-006 opens: Input email and select `Developer` role.
  3. Navigate to Configuration ➔ Environments S-015: Click "Create Environment" to add variables.
  4. Navigate to Configuration ➔ AI Providers S-014: Set up Ollama connection.
* **Recovery Actions**:
  - Missing key variables: Highlight env list with yellow badge warning that no active workspaces can run without configurations.

### Journey 3.3: Workspace Allocation & Provisioning (Developer)
* **Goal**: Launch a development workspace from zero projects.
* **Entry Point**: Projects list page S-007 showing empty projects state E-007.
* **Flow**:
  1. Click "New Project" to define workspace options.
  2. Select infrastructure provider target (Local, Docker, WSL, Kubernetes, SSH).
  3. Select Git repository URL to clone.
  4. Trigger "Create Workspace". System opens overlay showing cloning logs.
  5. Workspace provisions, indexer finishes, opens Workspace IDE page S-024.
* **Recovery Actions**:
  - Git clone failed: Show clone logs modal with red indicator, offering "Retry clone" or "Change credentials" action.
  - Allocation timeout: Cancel task and return user to Projects list with an error notification.

### Journey 3.4: AI Provider Setup & Fallbacks (Developer / Admin)
* **Goal**: Recover when AI operations fail due to misconfiguration or service outage.
* **Entry Point**: Workspace S-024 ➔ AI Section S-024e.
* **Flow**:
  1. Developer inputs prompt into chat. AI box displays E-006: "AI not configured".
  2. Developer clicks "Configure AI" CTA button (only visible to Admins, otherwise shows "Contact Admin").
  3. Redirects to AI Providers panel S-014.
  4. Admin adds OpenAI API credentials and clicks "Test Connection".
  5. Test connection returns green "Connected" status indicator.
  6. Developer returns to Workspace, chat assistant is now active.
* **Recovery Actions**:
  - Connection failed: Provider page parses response logs and suggests checking API keys or network configurations.

### Journey 3.5: Runtime Executor Interruption (Developer / Operator)
* **Goal**: Gracefully handle the loss of connection to a remote SSH or Kubernetes executor.
* **Entry Point**: Active workspace S-024 ➔ Task Runner Console S-024d.
* **Flow**:
  1. Developer triggers task execution.
  2. Executor connection drops. Workspace status indicator switches to blinking red "Executor Unavailable".
  3. Action buttons disable. File tree becomes read-only to prevent unsaved changes.
  4. System attempts to ping executor in background with exponential backoff (1s, 2s, 4s, 8s, max 16s).
  5. Once link returns, "Active" state is restored, file buffer is written, and task console allows retrying.
* **Recovery Actions**:
  - Re-establish fail: If link remains down for 60s, show dialog modal advising to switch workspace target to "Local" execution environment.

### Journey 3.6: Offline Mode & Interruption Recovery (All Users)
* **Goal**: Keep working during local browser network dropouts.
* **Entry Point**: Global Shell.
* **Flow**:
  1. Network drops. Browser fires offline event. Topbar shows "Offline View" warning banner.
  2. UI disables all mutations (invites, edits, triggers, saves). Buttons gray out.
  3. User can continue to browse active workspace files using cached offline tree.
  4. Editor locks saving, stores active work in browser local storage.
  5. Once internet connection is restored, banner disappears, pending changes sync to workspace automatically.
* **Recovery Actions**:
  - Sync conflict: If files changed on server while user was offline, open merge/diff editor to resolve discrepancies before writing.

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
