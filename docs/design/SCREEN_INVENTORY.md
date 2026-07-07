# Screen Inventory — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Full-Page Screens

### 1.1 Authentication & Bootstrap
| ID | Screen Name | Route | Component File | Actors |
|---|---|---|---|---|
| S-001 | Bootstrap Decision | `/` | `components/bootstrap-router.tsx` | System |
| S-002 | Bootstrap Loading | `/` | `components/bootstrap-screen.tsx` | System |
| S-003 | Setup Wizard | `/setup` | `app/setup/page.tsx` | Root Admin |
| S-004 | Login | `/login` | `app/login/page.tsx` | All Users |

### 1.2 Dashboard & Overview
| ID | Screen Name | Route | Component File | Actors |
|---|---|---|---|---|
| S-005 | Enterprise Dashboard | `/dashboard` | `dashboard/dashboard-page.tsx` | Admin, Owner |
| S-006 | Dashboard Shell | `/dashboard` | `dashboard/dashboard-shell.tsx` | All Auth |
| S-007 | Projects List | `/projects` | `app/projects/page.tsx` | All Auth |
| S-008 | Project Detail | `/projects/[slug]` | `app/projects/page.tsx` | Developer, Admin |

### 1.3 Enterprise & System Admin
| ID | Screen Name | Route | Component File | Actors |
|---|---|---|---|---|
| S-009 | Telemetry Monitor | `/monitor` | `app/monitor/page.tsx` | Operator, Admin |
| S-010 | Deployments Board | `/deploy` | `app/deploy/page.tsx` | Operator, Admin |
| S-011 | Backup & Restore Hub | `/backup` | `app/backup/page.tsx` | Admin, Owner |
| S-012 | Settings Hub | `/settings` | `settings/settings-page.tsx` | Admin, Owner |
| S-013 | User Profile Settings | `/profile` | `app/profile/page.tsx` | All Auth |

### 1.4 Configuration
| ID | Screen Name | Route | Component File | Actors |
|---|---|---|---|---|
| S-014 | AI Providers | `/config/providers` | `config/providers-panel.tsx` | Admin |
| S-015 | Environments & Secrets | `/config/environments` | `config/environments-panel.tsx` | Admin, Editor |

### 1.5 DevCenter
| ID | Screen Name | Route | Component File | Actors |
|---|---|---|---|---|
| S-016 | Architecture Docs | `/devcenter/architecture` | `app/devcenter/page.tsx` | Developer, Admin |
| S-017 | Database Schema | `/devcenter/database` | `app/devcenter/page.tsx` | Developer, Admin |
| S-018 | Tasks & Engine | `/devcenter/tasks` | `app/devcenter/page.tsx` | Developer |
| S-019 | Dependency Graph | `/devcenter/dependencies` | `app/devcenter/page.tsx` | Developer |
| S-020 | API Reference | `/devcenter/api` | `app/devcenter/page.tsx` | Developer |
| S-021 | Knowledge Base | `/devcenter/knowledge` | `app/devcenter/page.tsx` | All Auth |
| S-022 | Workspace Doctor | `/devcenter/project-doctor` | `app/devcenter/page.tsx` | Developer, Admin |
| S-023 | AI Cost & Token Usage | `/devcenter/ai-usage` | `app/devcenter/page.tsx` | Admin, Owner |

### 1.6 Workspace Runtime (S-024)
| ID | Screen Name | Route | Component File | Actors |
|---|---|---|---|---|
| S-024 | Workspace Explorer | `/workspace/[id]` | `workspace/workspace-page-client.tsx` | Developer |
| S-025 | Workspace Layout | `/workspace/[id]` | `workspace/workspace-layout.tsx` | Developer |
| S-026 | Code Editor | `/workspace/[id]` | `workspace/workspace-editor.tsx` | Developer |

#### Workspace Embedded Sections (within S-024)
| ID | Section Name | Component File | Purpose |
|---|---|---|---|
| S-024a | Overview | `workspace/overview-section.tsx` | Health summary & git meta |
| S-024b | Files | `workspace/files-section.tsx` | File Tree Explorer |
| S-024c | Services | `workspace/services-section.tsx` | Local systems status |
| S-024d | Tasks | `workspace/tasks-section.tsx` | Task Runner CLI & Logs |
| S-024e | AI Assistant | `workspace/ai-section.tsx` | Workspace prompt chat |

---

## 2. Modals, Drawers & Overlay Dialogs

| ID | Overlay Name | Location | Type | Purpose |
|---|---|---|---|---|
| M-001 | Command Palette | Global (Ctrl+K) | Modal | Universal search and jump actions |
| M-002 | Notification Inbox | Topbar icon | Drawer | Notification center, mentions, job events |
| M-003 | Help & Guide Popover | Topbar icon | Popover | Keyboard shortcuts list and tutorial lookup |
| M-004 | System Status dropdown| Topbar profile | Popover | Executor health, task queue sizes |
| M-005 | Create Organization | Settings Hub | Dialog | Input name/details for new organization |
| M-006 | Invite Member | Settings Hub | Dialog | Email + Role invite dispatch |
| M-007 | Edit Member Role | Settings Hub | Dialog | Dropdown selection to reassign roles |
| M-008 | Transfer Ownership | Settings Hub | Dialog | Destructive owner pass (Confirm block) |
| M-009 | Delete Member | Settings Hub | Dialog | Destructive member remove confirm |
| M-010 | Delete Env Variable | Config panel | Dialog | Variable delete confirmation |
| M-011 | Provider Test | Config panel | Inline | Live network endpoint test spinner |
| M-012 | Backup & Restore Wizard| Backup panel | Dialog | Step-by-step target backup or restore |
| M-013 | Clone Workspace | Projects list | Dialog | Clone workspace git repo path input |
| M-014 | Archive Workspace | Projects list | Dialog | Read-only workspace archive lock confirm |
| M-015 | Delete Workspace | Projects list | Dialog | Destructive workspace delete (Type confirmation)|
| M-016 | Duplicate Environment | Config panel | Dialog | Duplicate vars list to a new env target |
| M-017 | Rollback Release | Deployments | Dialog | Select release target rollback confirm |
| M-018 | Promote Release | Deployments | Dialog | Push staging deployment config to prod |

---

## 3. UI States (Empty, Error, Loading)

### 3.1 Empty States
| ID | Screen / Context | Trigger Condition | Success State |
|---|---|---|---|
| E-001 | Organizations List | User has no organization memberships | Add illustration + "Create Org" CTA |
| E-002 | Members Table | Zero members found in active org | Add illustration + "Invite Member" CTA |
| E-003 | Environments List | No env targets configured for organization | Add illustration + "Add Environment" CTA |
| E-004 | Variables Table | Selected environment contains zero variables | Add label + "Create Variable" CTA |
| E-005 | Secrets Table | Selected environment contains zero secrets | Add key icon + "Add Secret" CTA |
| E-006 | AI Providers | No providers active or configured | Add card grid + "Configure Provider" CTA |
| E-007 | Project List | No projects created in workspace | Onboarding card + "New Project" CTA |
| E-008 | Deployments Table | Project has never triggered a deploy run | Add timeline + "Trigger Deploy" CTA |
| E-009 | Task Runner Console | No tasks configured or triggered | Add console log + "Create Task" CTA |
| E-010 | Active Services | Workspace lists zero active daemon processes | Renders warning + check systems doc |
| E-011 | Command Palette | Keyboard query returns zero matches | "No commands or navigation destinations found" |
| E-012 | Notification Inbox | No new alert events or inbox messages | Check mark icon + "All caught up" |

### 3.2 Error States
| ID | Screen / Context | Trigger Condition | Visual Handler |
|---|---|---|---|
| ER-001 | Login Form | Invalid credentials | Inline red banner under header |
| ER-002 | Dashboard / Shell | Backend API unreachable or 500 error | Full viewport alert banner with refresh |
| ER-003 | Settings Panel | User attempts non-permitted action (403) | Grayed components + tooltips warning |
| ER-004 | AI Providers Test | Provider connection ping fails | Text label red text with response debug log |
| ER-005 | Workspace Explorer | Specified workspace identifier not found (404) | Full page "Workspace not found" return link |
| ER-006 | Workspace runtime | Adapter indexer timeout / file tree hang | Warning alert banner + "Restart Indexer" |
| ER-007 | Deployments board | Build log execution fails | Red status dot + expandable failed logs terminal |
| ER-008 | Telemetry Monitor | Metric data feeds crash / offline | Line charts draw straight red line with error label |
| ER-009 | Notification Inbox | WebSocket alert stream disconnects | Tiny red indicator dot in topbar header |
| ER-010 | Global API request | Client rate-limited by gateway (429) | Alert banner: "Too many actions. Please wait." |

### 3.3 Loading States
| ID | Screen / Context | Implementation Details |
|---|---|---|
| L-001 | Bootstrap Screen | Pulse logo + animated bar indicators |
| L-002 | Dashboard Board | Skeleton card shapes matching resource metric blocks |
| L-003 | Table lists | Shimmer overlay lines across row data items |
| L-004 | Workspace Explorer | Split-pane shimmer matching tree left, editor right |
| L-005 | Telemetry charts | Wave charts draw straight loading lines |
| L-006 | Action Buttons | Switch label to spinner, disable clicking |

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
