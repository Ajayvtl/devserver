# User Journeys — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## Role Definitions

| Role | Scope | Description |
|---|---|---|
| **Owner** | Organization | Creates org, transfers ownership, full admin escalation |
| **Admin** | Organization | Manages members, providers, environments, settings |
| **Developer** | Workspace | Creates projects, writes code, runs tasks, uses AI |
| **Operator** | Infrastructure | Manages deployments, monitors services, backup/restore |
| **Viewer** | Read-only | Views dashboards, files, logs — cannot mutate |

---

## Journey 1: Owner

### 1.1 Platform Initialization

| Attribute | Value |
|---|---|
| **Goal** | Bootstrap DevServer for the first time and establish governance |
| **Entry Point** | `/` → Bootstrap Router → `/setup` |
| **Navigation** | Setup Wizard (S-003) → Login (S-004) → Settings (S-009) |

**Success Path:**
1. Visit root URL → System detects uninitialized state
2. Setup Wizard auto-launches → System check passes
3. Database provisioned → Admin account created
4. Redirect to `/login` → Authenticate with new credentials
5. Redirect to `/dashboard` → See empty dashboard
6. Navigate to `/settings` → Create first organization
7. Invite initial admin members → Assign roles
8. Configure first AI provider → Test connection
9. Create first environment → Add variables/secrets

**Failure Path:**
- DB provisioning fails → Wizard shows error banner with diagnostic detail
- Admin creation fails → Inline validation errors (duplicate email, weak password)
- Provider test fails → Connection error with retry; does not block setup

### 1.2 Ownership Transfer

| Attribute | Value |
|---|---|
| **Goal** | Transfer organization ownership to another admin |
| **Entry Point** | `/settings` → Members tab |
| **Navigation** | Settings (S-009) → Members Panel (S-011) → Transfer Modal (M-004) |

**Success Path:**
1. Navigate to Members Panel → Locate target user
2. Click "Transfer Ownership" → Confirmation dialog appears
3. Confirm action → API processes transfer
4. Current user demoted to Admin → Target user promoted to Owner
5. Toast: "Ownership transferred successfully"

**Failure Path:**
- Target user not an active member → Error: "User must be an active member"
- Permission denied → 403 response → "Contact current owner"

---

## Journey 2: Admin

### 2.1 Team Onboarding

| Attribute | Value |
|---|---|
| **Goal** | Invite developers and configure their access |
| **Entry Point** | `/settings` → Members tab |
| **Navigation** | Settings (S-009) → Members (S-011) → Invite Modal (M-002) |

**Success Path:**
1. Navigate to Settings → Select Members tab
2. Click "Invite Member" → Enter email + select role
3. Submit → Backend creates membership
4. New member appears in table with "Invited" status
5. Member logs in → Status changes to "Active"

**Failure Path:**
- Duplicate email → Inline error: "User already a member of this organization"
- Invalid email format → Field-level validation error
- Rate limit exceeded → Toast: "Too many invitations. Please wait."

### 2.2 Environment & Provider Setup

| Attribute | Value |
|---|---|
| **Goal** | Configure runtime environments and AI providers |
| **Entry Point** | `/config/environments` and `/config/providers` |
| **Navigation** | Sidebar → Configuration → Environments / Providers |

**Success Path:**
1. Create environment (Dev/Staging/Prod) → Set type
2. Add variables → Key/Value pairs saved
3. Add secrets → Encrypted with AES-256, stored as references
4. Navigate to Providers → Add AI provider (Ollama, OpenAI)
5. Enter credentials → Test connection → Green "Connected"
6. Save provider → Available for workspace AI features

**Failure Path:**
- Missing required fields → Inline validation
- Provider connection timeout → "Connection failed" with endpoint detail
- Secret encryption failure → 500 error → "Contact administrator"

### 2.3 Security & Audit Review

| Attribute | Value |
|---|---|
| **Goal** | Review member activity and platform mutations |
| **Entry Point** | `/settings` (future: `/audit`) |
| **Navigation** | Settings → Audit tab (WP-8.6.13, not yet built) |

**Success Path:**
1. Navigate to Audit Viewer → Filter by date/user/action
2. Review login events, config changes, member mutations
3. Export audit log for compliance

**Failure Path:**
- Audit UI not built → ❌ Currently no UI (backend EventBus logs exist)
- Permission denied → 403 if viewer-level user attempts access

---

## Journey 3: Developer

### 3.1 Daily Development Workflow

| Attribute | Value |
|---|---|
| **Goal** | Write code, run tasks, use AI assistant, commit changes |
| **Entry Point** | `/dashboard` → Projects → Workspace |
| **Navigation** | Dashboard (S-005) → Projects (S-012) → Workspace (S-014) |

**Success Path:**
1. Login → Dashboard shows health metrics
2. Navigate to Projects → Select active project
3. Open Workspace → File tree loads, editor initializes
4. Browse files → Select file → Editor displays content
5. Switch to Tasks tab → Run build/test script
6. Switch to AI section → Ask about error context
7. AI responds with contextual solution
8. Switch to Repository tab → Review changed files
9. Commit changes via Git

**Failure Path:**
- Workspace not found → 404 page with "Return to projects"
- Indexer timeout → Loading state with retry option
- Task execution fails → Error detail in task output pane
- AI provider unavailable → "No AI provider configured" empty state

### 3.2 Project Creation

| Attribute | Value |
|---|---|
| **Goal** | Create a new project with repository binding |
| **Entry Point** | `/projects` → "New Project" |
| **Navigation** | Projects (S-012) → New Project Form → Project Detail (S-013) |

**Success Path:**
1. Navigate to Projects → Click "New Project"
2. Enter name, slug, description, repository URL
3. Configure build/start commands
4. Submit → Project created with unique slug
5. Redirect to project detail → All tabs available

**Failure Path:**
- Duplicate slug → Inline validation error
- Invalid repository URL → Warning but non-blocking
- Missing required fields → Form-level validation

---

## Journey 4: Operator

### 4.1 Infrastructure Monitoring

| Attribute | Value |
|---|---|
| **Goal** | Monitor service health and respond to incidents |
| **Entry Point** | `/dashboard` → Services section |
| **Navigation** | Dashboard (S-005) → Workspace Infrastructure (S-014e) → Services (S-014f) |

**Success Path:**
1. Login → Dashboard shows health metrics
2. Review system status cards → Identify warnings
3. Navigate to Workspace → Infrastructure tab
4. Check detected tools and versions
5. Switch to Services tab → Review running processes
6. Restart degraded service → Status returns to "Running"

**Failure Path:**
- Service restart fails → Error status with log output
- Monitoring dashboard not built → ❌ Placeholder at `/monitor`
- Cannot detect infrastructure → Empty tools list

### 4.2 Deployment Pipeline

| Attribute | Value |
|---|---|
| **Goal** | Deploy code to staging/production environments |
| **Entry Point** | `/workspace/[id]` → Deployments tab |
| **Navigation** | Workspace (S-014) → Deployments (S-014h) |

**Success Path:**
1. Navigate to Workspace → Select Deployments tab
2. Review deployment history → Branch, commit, status
3. Trigger deployment → Select target environment
4. Monitor progress → Success confirmation
5. Verify health check passes

**Failure Path:**
- Deployment UI is placeholder → ❌ No backend integration
- Build failure → Error detail with log output
- Health check fails → Automatic rollback (future scope)

---

## Journey 5: Viewer

### 5.1 Read-Only Exploration

| Attribute | Value |
|---|---|
| **Goal** | View dashboards, files, and project state without mutations |
| **Entry Point** | `/login` → `/dashboard` |
| **Navigation** | Dashboard (S-005) → Projects (S-012) → Workspace (S-014) |

**Success Path:**
1. Login with viewer credentials → Dashboard loads
2. Browse projects → Select project to view
3. Open workspace → File tree is browsable
4. View files → Editor in read-only mode
5. View environment variables → Secrets are masked
6. View Git state → Read-only commit history

**Failure Path:**
- Attempt to edit file → Editor buffer is `readOnly`, mutation blocked
- Attempt to create environment → Button disabled, API returns 403
- Attempt to invite member → Button hidden, API returns 403

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins.
