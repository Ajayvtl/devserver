# Screen Inventory — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Full-Page Screens

### 1.1 Authentication & Bootstrap

| ID | Screen | Route | Component | Actors |
|---|---|---|---|---|
| S-001 | Bootstrap Decision | `/` | `bootstrap-router.tsx` | System |
| S-002 | Bootstrap Loading | `/` | `bootstrap-screen.tsx` | System |
| S-003 | Setup Wizard | `/setup` | `app/setup/page.tsx` | Root Admin |
| S-004 | Login | `/login` | `app/login/page.tsx` | All Users |

### 1.2 Dashboard

| ID | Screen | Route | Component | Actors |
|---|---|---|---|---|
| S-005 | Enterprise Dashboard | `/dashboard` | `dashboard/dashboard-page.tsx` | Admin, Owner |
| S-006 | Dashboard Shell | `/dashboard` | `dashboard/dashboard-shell.tsx` | All Auth |

### 1.3 Configuration

| ID | Screen | Route | Component | Actors |
|---|---|---|---|---|
| S-007 | AI Providers | `/config/providers` | `config/providers-panel.tsx` | Admin |
| S-008 | Environments & Secrets | `/config/environments` | `config/environments-panel.tsx` | Admin, Editor |

### 1.4 Settings & Governance

| ID | Screen | Route | Component | Actors |
|---|---|---|---|---|
| S-009 | Settings Hub | `/settings` | `settings/settings-page.tsx` | Admin, Owner |
| S-010 | Organizations Panel | `/settings` (tab) | `settings/orgs-panel.tsx` | Admin, Owner |
| S-011 | Members Management | `/settings` (tab) | `settings/members-panel.tsx` | Admin, Owner |

### 1.5 Projects

| ID | Screen | Route | Component | Actors |
|---|---|---|---|---|
| S-012 | Projects List | `/projects` | `app/projects/page.tsx` | All Auth |
| S-013 | Project Detail | `/projects/[slug]` | `app/projects/page.tsx` | Developer, Admin |

### 1.6 Workspace Runtime

| ID | Screen | Route | Component | Actors |
|---|---|---|---|---|
| S-014 | Workspace Explorer | `/workspace/[id]` | `workspace/workspace-page-client.tsx` | Developer |
| S-015 | Workspace Layout | `/workspace/[id]` | `workspace/workspace-layout.tsx` | Developer |
| S-016 | Code Editor | `/workspace/[id]` | `workspace/workspace-editor.tsx` | Developer |

#### Workspace Sections (embedded in S-014)

| ID | Section | Component |
|---|---|---|
| S-014a | Overview | `workspace/overview-section.tsx` |
| S-014b | Files | `workspace/files-section.tsx` |
| S-014c | Repository | `workspace/repository-section.tsx` |
| S-014d | Environment | `workspace/environment-section.tsx` |
| S-014e | Infrastructure | `workspace/infrastructure-section.tsx` |
| S-014f | Services | `workspace/services-section.tsx` |
| S-014g | Tasks | `workspace/tasks-section.tsx` |
| S-014h | Deployments | `workspace/deployments-section.tsx` |
| S-014i | Knowledge | `workspace/knowledge-section.tsx` |
| S-014j | Remaining (DB, Domains, Logs, AI, Doctor, Settings, MCP) | `workspace/remaining-sections.tsx` |

### 1.7 DevCenter

| ID | Screen | Route |
|---|---|---|
| S-017 | Architecture | `/devcenter/architecture` |
| S-018 | Tasks | `/devcenter/tasks` |
| S-019 | Knowledge Base | `/devcenter/knowledge` |
| S-020 | Dependencies | `/devcenter/dependencies` |
| S-021 | API Reference | `/devcenter/api` |
| S-022 | Database Schema | `/devcenter/database` |
| S-023 | Project Doctor | `/devcenter/project-doctor` |

### 1.8 Placeholder Screens (require action)

| ID | Screen | Route | Action Required |
|---|---|---|---|
| S-024 | Deploy | `/deploy` | Build pipeline UI or hide from nav |
| S-025 | Monitor | `/monitor` | Build monitoring or hide from nav |
| S-026 | Doctor | `/doctor` | Redirect to `/devcenter/project-doctor` |
| S-027 | Service | `/service` | Integrate or remove |
| S-028 | Backup | `/backup` | Build or hide from nav |
| S-029 | Install | `/install` | Build or hide from nav |

---

## 2. Modals & Dialogs

| ID | Modal | Trigger | States |
|---|---|---|---|
| M-001 | Create Organization | S-010 | Default, Validating, Success, Error |
| M-002 | Invite Member | S-011 | Default, Validating, Success, Duplicate Error |
| M-003 | Edit Member Role | S-011 | Default, Saving, Success, Permission Error |
| M-004 | Transfer Ownership | S-011 | Confirmation, Processing, Success, Error |
| M-005 | Delete Member | S-011 | Confirmation, Processing, Success, Error |
| M-006 | Delete Env Variable | S-008 | Confirmation, Processing, Success, Error |
| M-007 | Provider Test | S-007 | Testing, Connected, Failed |
| M-008 | Generic Confirmation | Global | Pending, Confirmed, Cancelled |

## 3. Drawers & Side Panels

| ID | Drawer | Location |
|---|---|---|
| D-001 | Workspace Section Sidebar | S-014 |
| D-002 | App Navigation Sidebar | App Shell |
| D-003 | Provider Detail Expand | S-007 |

## 4. Wizards

| ID | Wizard | Route | Steps |
|---|---|---|---|
| W-001 | Initial Setup | `/setup` | System Check → DB Init → Admin → Complete |
| W-002 | Environment Creation | `/config/environments` | Name → Type → Save |
| W-003 | Provider Config | `/config/providers` | Type → Credentials → Test → Save |

## 5. Empty States

| ID | Screen | Condition | Status |
|---|---|---|---|
| E-001 | Organizations | Zero orgs | ✅ Implemented |
| E-002 | Members | Zero members | ✅ Implemented |
| E-003 | Variables | No env vars | ✅ Implemented |
| E-004 | Secrets | No secrets | ✅ Implemented |
| E-005 | Providers | No providers | ✅ Implemented |
| E-006 | Projects | No projects | 🟡 Partial |
| E-007 | Workspace Files | Empty dir | 🟡 Partial |
| E-008 | Deployments | No history | ❌ Placeholder data |
| E-009 | Tasks | No tasks | 🟡 Partial |
| E-010 | Audit Logs | No events | ❌ Missing UI |
| E-011 | Knowledge | No articles | 🟡 Partial |
| E-012 | Search Results | No matches | ❌ Missing |

## 6. Error States

| ID | Screen | Condition | Status |
|---|---|---|---|
| ER-001 | Login | Invalid credentials | ✅ Toast |
| ER-002 | Login | Network failure | 🟡 Generic |
| ER-003 | Dashboard | API unreachable | ❌ Blank |
| ER-004 | Settings | 403 Forbidden | ✅ Handled |
| ER-005 | Members | Duplicate email | ✅ Backend |
| ER-006 | Providers | Connection fail | ✅ Feedback |
| ER-007 | Environments | Create failure | ✅ Generic toast |
| ER-008 | Workspace | 404 Not found | 🟡 Partial |
| ER-009 | Workspace | Indexer timeout | ❌ Hangs |
| ER-010 | Global | 401 Session expired | ✅ Redirect |
| ER-011 | Global | 429 Rate limited | ❌ Missing |
| ER-012 | Global | 500 Server error | 🟡 Generic |

## 7. Loading States

| ID | Screen | Status |
|---|---|---|
| L-001 | Bootstrap | ✅ Animated pulse bars |
| L-002 | Dashboard | 🟡 Text-only |
| L-003 | Settings | ✅ Shimmer skeletons |
| L-004 | Providers | 🟡 Basic |
| L-005 | Environments | 🟡 Basic |
| L-006 | Workspace | ✅ Partial skeleton |
| L-007 | Projects | 🟡 Basic |
| L-008 | Login | ✅ Button state |

## 8. Summary

| Category | Count |
|---|---|
| Full-page screens | 29 |
| Workspace sections | 10 |
| Modals & dialogs | 8 |
| Drawers & panels | 3 |
| Wizards | 3 |
| Empty states | 12 |
| Error states | 12 |
| Loading states | 8 |
| **Total UI surfaces** | **85** |

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation code may be written.
