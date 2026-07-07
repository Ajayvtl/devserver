# API ➔ UI Mapping — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## Mapping Convention

Every UI widget maps through this comprehensive lifecycle chain:

```
Widget ➔ API Endpoint ➔ State Management ➔ Loading State ➔ Error State ➔ Permission Gate ➔ Fallback UI ➔ Retry & Offline Policy ➔ Cache Policy
```

---

## 1. Authentication & Session

| Widget | API Endpoint | State Management | Loading State | Error State | Permission | Fallback UI | Retry & Offline Policy | Cache Policy |
|---|---|---|---|---|---|---|---|---|
| **Login Form** | `POST /api/v1/auth/login` | `AuthContext.tokens` | Disable submit button; show spinner | Alert banner inside card | Anonymous | Remain on form; enable inputs | Manual retry only; 5s lockout after 5 failures | None (always live) |
| **Token Refresh** | `POST /api/v1/auth/refresh` | `AuthContext.tokens` | Silent background update | Redirect to `/login` | Authenticated | Interstitial modal with login redirect | Auto-retry twice with 2s backoff; fallback to offline cache if no network | None |
| **User Identity** | `GET /api/v1/users/me` | `AuthContext.user` | Topbar profile skeleton | Avatar placeholder with indicator | Authenticated | Render anonymous icon | Auto-retry 3x; then disable features needing auth | SWR; cache in local storage (24h TTL) |
| **Profile Panel** | `PUT /api/v1/users/me` | Local state ➔ `AuthContext` sync | Button: "Saving..." spinner | Inline field-level error messages | Authenticated | Keep editing state with warning | Manual retry | None |
| **MFA Toggle** | `POST /api/v1/auth/mfa` | Local boolean state | Dialog: "Verifying..." | Dialog banner | Authenticated | Keep switch toggled off | Manual retry | None |
| **Sessions List** | `GET /api/v1/users/sessions` | `sessions[]` state | Skeleton rows | Banner: "Failed to load sessions" | Authenticated | "No other sessions" state | Manual retry button | SWR; refresh on panel mount |
| **API Keys Card** | `GET /api/v1/users/keys` | `apiKeys[]` state | Skeleton cards | Banner: "Failed to load keys" | Admin+ | Empty card; "Generate Key" CTA | Manual retry button | Cached (session limit) |

---

## 2. Organization & Members (Administration)

| Widget | API Endpoint | State Management | Loading State | Error State | Permission | Fallback UI | Retry & Offline Policy | Cache Policy |
|---|---|---|---|---|---|---|---|---|
| **Org List Table** | `GET /api/v1/organizations` | `orgs[]` via `useAuth` | Skeleton rows | Full table alert banner | Authenticated | EmptyState: "No orgs found" + CTA | Auto-retry 3x (exponential backoff); Offline view of cached orgs | SWR; polling every 60s |
| **Create Org Form** | `POST /api/v1/organizations` | Optimistic update ➔ Context | Button: "Creating..." spinner | Inline field validation | Admin+ | Keep dialog open; enable inputs | Manual retry | None |
| **Org Switcher** | Client-side switch | `currentOrgId` in Context | Select disabled | Toast error on config load | Authenticated | Disable switching; lock to active org | Manual reload on fail | Live context variable |
| **Members Table** | `GET /api/v1/organizations/members` | `useMembers.members` | Skeleton rows | Banner: "Unable to load members" | Admin+ | EmptyState: "No members" | Manual retry button | SWR; polling every 30s |
| **Invite Form** | `POST /api/v1/organizations/members` | Optimistic insert | Button: "Inviting..." | Inline validation / toast | Admin+ | Remain on form; preserve email | Manual retry | None |
| **Role Dropdown** | `PUT /api/v1/organizations/members` | Update list member | Row-level spinner | Toast: "Failed to update role" | Admin+ | Revert dropdown to previous | Auto-retry 1x; show offline toast | None |
| **Status Toggle** | `PUT /api/v1/organizations/members` | Update list member | Row-level spinner | Toast: "Failed to update status" | Admin+ | Revert toggle to previous | Auto-retry 1x; fallback to offline | None |
| **Remove Member** | `DELETE /api/v1/organizations/members` | Remove from list | Row opacity fade | Toast: "Failed to remove member" | Admin+ | Restore row in list | Manual retry | None |
| **Ownership Transfer** | `POST /api/v1/organizations/transfer` | Reset Org owner state | Dialog blocker with spinner | Dialog-level critical error | Owner only | Dismiss dialog; retain owner | Manual retry; confirmation prompt | None |

---

## 3. Configuration & Environments

| Widget | API Endpoint | State Management | Loading State | Error State | Permission | Fallback UI | Retry & Offline Policy | Cache Policy |
|---|---|---|---|---|---|---|---|---|
| **AI Providers Grid** | `GET /api/v1/providers` | `providers[]` state | Grid skeleton | Alert banner | Admin+ | EmptyState: "AI not configured" | Auto-retry 3x; then show offline message | SWR; manual refresh only |
| **Provider Form** | `POST /api/v1/providers` | Optimistic card append | Button: "Saving..." spinner | Banner: "Save failed" | Admin+ | Keep editor drawer open | Manual retry | None |
| **Provider Test Button**| `POST /api/v1/providers/test`| `provider.status` | Spinner + "Testing..." label | Inline label: "Connection failed" | Admin+ | Icon: Connection Warning | Manual retry; 10s timeout | None |
| **Environments List** | `GET /api/v1/environments` | `environments[]` state | Card skeleton | Full list error banner | Admin+ | EmptyState: "No environments" | Auto-retry 3x; support cached offline view | SWR; manual refresh |
| **Variables Table** | `GET /api/v1/variables` | `variables[]` state | Table skeleton | Error banner | Admin+ | EmptyState: "No variables set" | Auto-retry 3x; cache in local db | SWR |
| **Secrets Table** | `GET /api/v1/secrets` | `secrets[]` (metadata refs) | Table skeleton | Error banner | Admin+ | EmptyState: "No secrets set" | Auto-retry 3x; cache refs in local db | SWR |
| **Duplicate Env Flow** | `POST /api/v1/environments/duplicate`| Client redirect | Dialog spinner | Dialog banner | Admin+ | Keep dialog open | Manual retry | None |

---

## 4. Workspace & Development (Developer Experience)

| Widget | API Endpoint | State Management | Loading State | Error State | Permission | Fallback UI | Retry & Offline Policy | Cache Policy |
|---|---|---|---|---|---|---|---|---|
| **Workspace Shell** | `GET /api/workspaces/[id]` | `WorkspaceContext` | Full-screen shimmer | "Workspace not found" page | Developer+ | Redirect to Projects list | Auto-retry 3x; fallback to offline cache | SWR |
| **File Tree** | `GET /api/workspaces/[id]/files` | `FileTreeContext` | Tree skeleton nodes | Error label in sidebar | Developer+ | "Empty Workspace" CTA | Auto-retry 3x; offline tree cache | Cached; SWR; manual reload |
| **Git Status** | `GET /api/workspaces/[id]/git` | `GitContext` | File status spinner | "Git unavailable" tag | Developer+ | Disabled panel; warning banner | Auto-retry 3x; offline git snapshot | Live SWR; polling 15s |
| **Active Services** | `GET /api/workspaces/[id]/services`| `services[]` state | Service rows skeleton | Alert banner | Developer+ | "No active services" card | Auto-retry 5x; offline fallback | Live SWR; polling 10s |
| **Task Queue List** | `GET /api/tasks` | `useTasks.tasks` | Skeleton items | Error card | Developer+ | "No tasks in queue" card | Re-establish WebSocket; fallback poll | Live WebSocket stream; poll 5s |
| **AI Chat Feed** | `WS /ws/events` | `useAI.messages` | Shimmer on last assistant card | Toast: "Message failed" | Developer+ | Render chat panel offline info | Auto-reconnect WebSocket; 30s timeout | Session memory; no local storage cache |
| **Clone Workspace** | `POST /api/workspaces/clone` | Redirect | Dialog progress bar | Dialog alert | Developer+ | Dismiss dialog | Manual retry | None |
| **Archive Workspace** | `POST /api/workspaces/archive` | List refresh | Spinner on row | Toast error | Developer+ | Revert status in list | Manual retry | None |

---

## 5. Operations & Monitoring

| Widget | API Endpoint | State Management | Loading State | Error State | Permission | Fallback UI | Retry & Offline Policy | Cache Policy |
|---|---|---|---|---|---|---|---|---|
| **CPU/RAM Timelines** | `GET /api/v1/monitor/metrics`| `metrics[]` timeline state | Chart shimmer placeholder | "Metrics unavailable" banner | Operator+ | Static "Telemetry Offline" card | Auto-retry every 10s; offline cache data | Polling every 5s |
| **Process List** | `GET /api/v1/monitor/processes`| `processes[]` state | Table skeleton | Error banner | Operator+ | EmptyState: "No processes found" | Auto-retry 5s; cache last status | Live SWR; polling 5s |
| **Incidents Alert Bar** | `GET /api/v1/monitor/incidents`| `incidents[]` active list | Silent background load | Hidden | Operator+ | Render warning banner in topbar | Auto-retry 30s; local DB cache | Live WebSocket + polling 30s |
| **Deployments Table** | `GET /api/workspaces/[id]/deployments`| `deployments[]` state | Table skeleton | Error banner | Operator+ | EmptyState: "No deployment history"| Auto-retry 3x; SWR local cache | SWR |
| **Rollback Form** | `POST /api/workspaces/rollback`| Redirect/Refresh | Screen block dialog | Toast: "Rollback failed" | Operator+ | Dismiss dialog; keep details | Manual retry; strict confirmation | None |
| **Promote Env Form** | `POST /api/workspaces/promote` | Refresh target env | Screen block dialog | Toast: "Promotion failed" | Operator+ | Dismiss dialog | Manual retry; strict confirmation | None |
| **Backup Table** | `GET /api/v1/backup/jobs` | `backupJobs[]` state | Table skeleton | Error banner | Admin+ | EmptyState: "No backup schedule" | Auto-retry 3x | SWR |
| **Trigger Restore** | `POST /api/v1/backup/restore` | Restore progress task | Progress wizard overlay | Wizard alert banner | Admin+ | Keep wizard open; block mutations | Manual retry; strict confirm | None |

---

## 6. Global Systems (Search, Notifications, Help)

| Widget | API Endpoint | State Management | Loading State | Error State | Permission | Fallback UI | Retry & Offline Policy | Cache Policy |
|---|---|---|---|---|---|---|---|---|
| **Global Search** | `GET /api/v1/search?q=` | Local query result state | Inline search spinner | "Search failed" message | Authenticated | "No results found" card | Manual retry; offline file paths only | Client-side query cache (5m TTL) |
| **Command Palette** | `GET /api/v1/commands` | `commands[]` registry | Silent fetch on mount | Fallback commands | Authenticated | Render core system commands | Local fallback list | Cached (session limit) |
| **Notification Inbox**| `GET /api/v1/notifications`| `unreadCount`, `items[]` | Popover skeleton | Badge turns gray | Authenticated | Icon badge hidden | Auto-reconnect WebSocket; poll 10s | Live WebSocket + SWR |
| **Timeline Feed** | `GET /api/v1/activity` | `activities[]` state | Skeleton timeline | "Timeline offline" panel | Authenticated | Static cached history | Auto-retry 30s; offline cache | SWR; poll 30s |
| **Help center docs** | `GET /api/v1/help/articles`| `articles[]` state | Popover skeleton | "Help offline" message | Authenticated | Static keyboard shortcuts card | Fallback to offline stored markdown | Local storage cache (30 days TTL) |

---

## Cache & Retry Policies (Standard Definitions)

### Cache Policies
- **None**: Always request live data, never cache (typically mutable POST/PUT/DELETE forms).
- **SWR (Stale-While-Revalidate)**: Render cached content instantly while triggering a background fetch to update the cache and re-render.
- **Polling [time]**: Automatically trigger a background GET request at the specified interval to update the UI state.
- **WebSocket**: Maintain persistent connection for real-time events. Fallback to SWR/Polling if connection drops.
- **Cached**: Store in local storage/IndexedDB with specified TTL (Time-To-Live).

### Retry & Offline Policies
- **Manual Retry**: Show a "Retry" button. Do not trigger network requests automatically.
- **Auto-retry [count] with [backoff]**: Automatically retry the API call up to `count` times, doubling the delay between each attempt (exponential backoff).
- **Offline Mode**: If the user is offline or the server is unreachable:
  - Mask action buttons (disable mutations).
  - Serve stale data from the local cache with a visual "Offline View" indicator.
  - Queue mutations in IndexedDB for synchronization when online (optional per flow).

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
