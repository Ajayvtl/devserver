# API → UI Mapping — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## Mapping Convention

Every UI widget maps through this chain:

```
Widget → API Endpoint → State Management → Loading State → Error State → Permission Gate → Fallback
```

---

## 1. Authentication

| Widget | API | State | Loading | Error | Permission | Fallback |
|---|---|---|---|---|---|---|
| Login Form | `POST /api/v1/auth/login` | `AuthContext.tokens` | Button: "Logging in..." | Red toast with error message | Anonymous | Retry form |
| Token Refresh | `POST /api/v1/auth/refresh` | `AuthContext.tokens` | Silent (background) | Redirect to `/login` | Authenticated | Session expired interstitial |
| User Identity | `GET /api/v1/users/me` | `AuthContext.user` | Spinner in topbar | Redirect to `/login` | Authenticated | "Unknown user" |

## 2. Organizations

| Widget | API | State | Loading | Error | Permission | Fallback |
|---|---|---|---|---|---|---|
| Org List Table | `GET /api/v1/organizations` | `orgs[]` via `useAuth` | Skeleton table rows | Error banner + retry | Authenticated | Empty: "No organizations found" |
| Create Org Form | `POST /api/v1/organizations` | Optimistic add to `orgs[]` | Button: "Creating..." | Inline validation + toast | Admin+ | N/A |
| Org Switcher (topbar) | Client-side (from `orgs[]`) | `currentOrgId` | Select disabled | N/A | Authenticated | First org auto-selected |

## 3. Members

| Widget | API | State | Loading | Error | Permission | Fallback |
|---|---|---|---|---|---|---|
| Members Table | `GET /api/v1/organizations/members` | `useMembersManagement.members` | Skeleton rows | ErrorState component | Admin+ (403 → access denied) | Empty: "No members yet" |
| Invite Form | `POST /api/v1/organizations/members` | Append to `members[]` | Button: "Inviting..." | Inline: duplicate email, missing fields | Admin+ | N/A |
| Role Dropdown | `PUT /api/v1/organizations/members` | Update member in list | Spinner on row | Toast: "Failed to update role" | Admin+ (403 check) | Current role preserved |
| Status Toggle | `PUT /api/v1/organizations/members` | Update member status | Spinner on row | Toast: "Failed to update status" | Admin+ | Current status preserved |
| Remove Member | `DELETE /api/v1/organizations/members?userId=` | Remove from `members[]` | Row fade-out | Toast: "Failed to remove" | Admin+ (403 check) | Row restored |
| Transfer Ownership | `POST /api/v1/organizations/transfer-ownership` | Refresh all org data | Modal: "Transferring..." | Modal error message | Owner only | Modal dismissed |

## 4. Roles

| Widget | API | State | Loading | Error | Permission | Fallback |
|---|---|---|---|---|---|---|
| Roles Dropdown | `GET /api/v1/roles` | `roles[]` via hook | Spinner in dropdown | Fallback: `[{id: "role-1", name: "admin"}]` | Authenticated | Default admin role |

## 5. Settings

| Widget | API | State | Loading | Error | Permission | Fallback |
|---|---|---|---|---|---|---|
| Settings List | `GET /api/v1/settings?scope=&ownerId=` | `settings[]` | Skeleton rows | Error banner | Admin+ | Empty: "No settings configured" |
| Save Setting | `POST /api/v1/settings` | Update in `settings[]` | Button: "Saving..." | Toast: "Failed to save" | Admin+ | Previous value preserved |

## 6. Environments

| Widget | API | State | Loading | Error | Permission | Fallback |
|---|---|---|---|---|---|---|
| Env List | `GET /api/v1/environments?ownerId=` | `environments[]` | Skeleton cards | Error banner + retry | Admin+ | Empty: "No environments" |
| Create Env | `POST /api/v1/environments` | Add to `environments[]` | Button: "Creating..." | Inline validation + toast | Admin+ | N/A |
| Variables Table | `GET /api/v1/variables?envId=` | `variables[]` | Skeleton rows | Error banner | Admin+ | Empty: "No variables" |
| Save Variable | `POST /api/v1/variables` | Add to `variables[]` | Button: "Saving..." | Toast: "Failed to save" | Admin+ | N/A |
| Secrets Table | `GET /api/v1/secrets?envId=` | `secrets[]` (refs only) | Skeleton rows | Error banner | Admin+ | Empty: "No secrets" |
| Save Secret | `POST /api/v1/secrets` | Add ref to `secrets[]` | Button: "Saving..." | Toast: "Failed to save" | Admin+ | N/A |

## 7. Providers

| Widget | API | State | Loading | Error | Permission | Fallback |
|---|---|---|---|---|---|---|
| Provider List | `GET /api/v1/providers?scope=&ownerId=` | `providers[]` | Skeleton cards | Error banner + retry | Admin+ | Empty: "No providers configured" |
| Save Provider | `POST /api/v1/providers` | Add to `providers[]` | Button: "Saving..." | Toast: "Failed to save" | Admin+ | N/A |
| Test Connection | `POST /api/v1/providers/test` | `testResult` | Card: "Testing..." spinner | Card: "Connection failed" + error msg | Admin+ | N/A |

## 8. Workspace (WebSocket + REST)

| Widget | API | State | Loading | Error | Permission | Fallback |
|---|---|---|---|---|---|---|
| Overview | `GET /api/workspaces/[id]` | `WorkspaceOverview` | Full section skeleton | "Workspace not found" page | Developer+ | N/A |
| File Tree | `GET /api/workspaces/[id]/files` | `WorkspaceFilesResponse` | Skeleton tree nodes | Retry banner | Developer+ | Empty: "No files found" |
| Git Info | `GET /api/workspaces/[id]/git` | `WorkspaceGitInfo` | Skeleton rows | "Git not available" message | Developer+ | "Not a git repository" |
| Environment | `GET /api/workspaces/[id]/environment` | `WorkspaceEnvironmentInfo` | Skeleton rows | Error banner | Developer+ | Empty: "No env files" |
| Infrastructure | `GET /api/workspaces/[id]/infrastructure` | `InfrastructureInfo` | Skeleton tool cards | Error banner | Developer+ | Empty: "No tools detected" |
| Services | `GET /api/workspaces/[id]/services` | `ProviderInfo[]` | Skeleton service rows | Error banner | Developer+ | Empty: "No services running" |
| Tasks | `POST /api/tasks` / `WS /ws/events` | `TaskEvent` stream | Progress bar | Error detail in output | Developer+ | Empty: "No tasks queued" |
| Deployments | `GET /api/workspaces/[id]/deployments` | `DeploymentInfo` | Skeleton rows | Error banner | Operator+ | Empty: "No deployments" |
| AI Chat | `POST /ws/events` (WebSocket) | Streaming tokens | "Thinking..." indicator | "AI provider unavailable" | Developer+ | "Configure an AI provider" |

## 9. Dashboard

| Widget | API | State | Loading | Error | Permission | Fallback |
|---|---|---|---|---|---|---|
| Metrics Grid | `GET /api/dashboard` | `DashboardDataV2.metrics` | Skeleton metric cards | Error banner | Authenticated | "Unable to load metrics" |
| Activity Feed | `GET /api/dashboard` | `DashboardDataV2.activity` | Skeleton timeline | Error banner | Authenticated | Empty: "No recent activity" |
| Task Queue | `GET /api/dashboard` | `DashboardDataV2.tasks` | Skeleton list | Error banner | Authenticated | Empty: "No active tasks" |
| Recommendations | `GET /api/dashboard` | `DashboardDataV2.recommendations` | Skeleton cards | Hidden on error | Authenticated | Hidden |

---

## 10. Dead Endpoints (API exists, no UI)

| API Endpoint | Status | Required Action |
|---|---|---|
| EventBus audit log events | Backend only | Build Audit Viewer UI (WP-8.6.13) |
| `GET /api/workspaces/[id]/deployments` | Data available | Build deployment pipeline (WP-8.6.14) |
| Health telemetry via providers | Backend only | Build Monitoring Dashboard (WP-8.6.15) |

## 11. Dead UI (UI exists, no API)

| UI Screen | Status | Required Action |
|---|---|---|
| `/deploy` | Placeholder component | Wire to deployment backend |
| `/monitor` | Empty directory | Wire to health monitoring backend |
| `/doctor` | Empty directory | Redirect to `/devcenter/project-doctor` |

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins.
