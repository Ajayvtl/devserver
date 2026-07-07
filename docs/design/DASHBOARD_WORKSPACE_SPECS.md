# Developer Operations Control Deck & Workspace Specifications

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Product Identity & Control Deck Layout (Gap #1 / Refinement)

The dashboard (`/dashboard` S-005) is structured as a **Developer Operations Platform Control Deck** rather than a passive charts list. It explicitly answers five core operational questions upon landing:

```
┌────────────────────────────────────────────────────────────────────────┐
│  [PLATFORM STATUS DECK]                                                │
│  ● HEALTH: Healthy (98%)  │  ● AI: Configured (Ollama)  │ ● BLOCKS: None│
├────────────────────────────────────────────────────────────────────────┤
│  [PRIORITY ACTION DECK]                                                │
│  ➔ WHAT TO DO NEXT: Deploy version v1.2.4-rc to Staging                 │
│  ⚠️ ATTENTION TODAY: Executor SSH-Remote has high memory usage (92%)     │
├────────────────────────────────────────────────────────────────────────┤
│  [W-1: Workspace Quick Cards] (8 Cols)    │ [W-2: Active Executors] (4)│
├───────────────────────────────────────────┼────────────────────────────┤
│  [W-3: Resource Telemetry] (8 Cols)       │ [W-4: Running Jobs] (4)    │
├───────────────────────────────────────────┼────────────────────────────┤
│  [W-5: Saved Layouts / Personalization] (12 Cols)                      │
└────────────────────────────────────────────────────────────────────────┘
```

### 1.1 Core Questions Answered in Landing Deck
1. **Is my platform healthy?**
   - *UI Visual*: Circular ring graph representing system health score (aggregate of executor latency, active services status, database connection, and storage availability).
2. **Is AI configured?**
   - *UI Visual*: Connection status badge next to active provider model config. Green for verified Ollama/OpenAI link, Orange for degraded, Red for unconfigured/offline with direct "Configure" path.
3. **Are deployments blocked?**
   - *UI Visual*: Release status strip showing pipeline blocks (e.g., failing integration runs, environment secret mismatch warning, or locked release gates).
4. **What needs attention today?**
   - *UI Visual*: "Needs Attention" alert list filtering warnings from all workspaces (e.g., "Disk space >90% on executor WSL", "Token budget exceeded 80% on OpenAI").
5. **What should I do next?**
   - *UI Visual*: Actionable recommended Next Step banner (e.g., "Workspace 'auth-service' has 4 uncommitted files. Open workspace to sync or commit.").

---

## 2. Dashboard Personalization Specs (Gap #4 / Refinement)

Enterprise developers can customize their Operations Control Deck layout using these options:

* **Pinned Workspaces & Favorites**:
  - Workspace cards in S-007 and S-005 contain a star icon button.
  - Starred items are pinned to the top of the Workspace Quick Cards widget (`W-1`).
* **Widget Visibility Settings**:
  - A settings button in the topbar of the Dashboard opens a dropdown checklist.
  - Users can toggle the visibility of individual widgets (e.g., hide the Resource Telemetry chart if running local-only).
* **Saved Layout Configurations**:
  - Layout is stored locally in browser `localStorage` as a JSON block mapping grid sizes and column locations (e.g., `{ widgetId: "W-3", size: "col-8", order: 2 }`).
  - Allows drag-and-drop rearrangement of dashboard panels using lightweight grid libraries.
* **Recent Items Registry**:
  - A client-side log tracks the last 5 workspace explorer routes visited, active command palette searches triggered, and logs files viewed. Stored inside the user's local session cache.

---

## 3. Workspace Layout & UX Specification

The workspace explorer (`/workspace/[id]` S-024) is designed as a split-pane IDE wrapper.

```
┌───────────┬──────────────┬───────────────────────────────┬──────────────┐
│  L1 Nav   │ L2 Sidebar   │ Editor Tabs Area              │ R1 Sidebar   │
│           │              ├───────────────────────────────┤              │
│  [Icons]  │  [Files      │                               │  [AI         │
│           │   Tree       │        Code Editor            │   Assistant  │
│           │   Browser    │                               │   Chat]      │
│           │              ├───────────────────────────────┤              │
│           │   or         │ Console / Terminal Output     │              │
│           │   Git state] │                               │              │
└───────────┴──────────────┴───────────────────────────────┴──────────────┘
```

### Panes & Resizing Rules
* **L2 Sidebar (Explorer)**: Default width `280px`. Left drag-handle allows resizing from `200px` to `450px`. Can be collapsed entirely.
* **R1 Sidebar (AI Panel)**: Default width `350px`. Right drag-handle allows resizing from `300px` to `600px`. Can be collapsed.
* **Center Editor Area**: Takes up remaining workspace width. Consists of:
  - **Tabs Area**: Row of active file tabs with scroll capability and "Close" indicators.
  - **Terminal / Console**: Docked at bottom. Default height `200px`. Top drag-handle allows resizing from `100px` to `500px`.

### Keyboard Shortcuts Reference

| Command | Action | Keyboard Shortcut (Win/Linux) | Keyboard Shortcut (Mac) |
|---|---|---|---|
| Toggle Left Sidebar | Collapse/expand file explorer | `Ctrl + \` | `Cmd + \` |
| Toggle Right Sidebar | Collapse/expand AI assistant | `Ctrl + Shift + A` | `Cmd + Shift + A` |
| Toggle Terminal Panel | Collapse/expand output terminal | `Ctrl + \`` | `Cmd + \`` |
| Search File name | Open workspace file search popup | `Ctrl + P` | `Cmd + P` |
| Open Command Palette | Open universal search modal | `Ctrl + K` | `Cmd + K` |
| Save active file | Write file buffers to workspace server | `Ctrl + S` | `Cmd + S` |
| Close active tab | Dismiss current editor file tab | `Ctrl + W` | `Cmd + Option + W` |

---

## 4. Performance Budgets (Gap #6 / Refinement)

To maintain a fast interface under heavy developer use, the frontend client and Go server must operate within these budgets:

| Page / Metric | Performance Target | API Latency Target (95th Pct) | Max Permitted Limit (Hard Stop) |
|---|---|---|---|
| **Dashboard Load** | Full control deck render | `< 150ms` | `< 1.0 second` |
| **Workspace IDE Boot** | Mount shell & file tree explorer | `< 200ms` | `< 2.0 seconds` |
| **Command Palette Search**| Query matching returns | `< 50ms` | `< 150ms` |
| **WebSocket Logs Feed** | Event push frequency | Dynamic batching (see below) | Max 100 log lines per 100ms |
| **Active Workspaces** | Max concurrent running per user| N/A | **5 active** (others put to sleep) |
| **Client memory** | Browser tab memory allocation | N/A | `< 250MB` heap usage |

### WebSocket Log Stream Batching Rule
To prevent browser UI threads from freezing during intense execution logs output (e.g., a rapid build run):
1. The Go server batches logs in a `100ms` buffer.
2. If log lines exceed 50 inside the window, they are sent as a single batched array payload.
3. The Next.js frontend uses a virtualized list wrapper to render console outputs, keeping DOM node count static at a maximum of `1000` visible rows.

---

## 5. Production Acceptance Checklist

Every page must satisfy this comprehensive checklist before receiving QA sign-off:

```
[ ] Responsive Layout: Verified on Mobile, Tablet, and Desktop screen widths.
[ ] Accessibility: Checked that contrast ratios are WCAG AA compliant.
[ ] Screen Reader Compatibility: Interactive elements labeled with ARIA attributes.
[ ] Color Themes: Checked that both Light and Dark mode styles render correctly.
[ ] UI States: Verified loading skeleton, empty dataset, and error banner components.
[ ] Error Handling: Handled 401 (redirect), 403 (action block), and 500 (API failure) states.
[ ] Resiliency: Verified "Offline View" fallback behavior.
```

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
