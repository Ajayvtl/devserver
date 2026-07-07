# Dashboard, Workspace, Performance & Acceptance Specifications

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Dashboard Widget Inventory (Gap #8)

The main dashboard (`/dashboard` S-005) uses a 12-column responsive layout grid. It contains the following widgets:

```
┌──────────────────────────────────────────────────────────────────────────┐
│  [W-1: Health Checklist] (8 Cols)          │ [W-2: Active Executors] (4) │
├────────────────────────────────────────────┼─────────────────────────────┤
│  [W-3: Resource Telemetry CPU/RAM] (8 Cols)│ [W-4: Running Jobs] (4)     │
├────────────────────────────────────────────┼─────────────────────────────┤
│  [W-5: Workspaces list] (6 Cols)           │ [W-6: Activity Timeline](6) │
├────────────────────────────────────────────┼─────────────────────────────┤
│  [W-7: Recommendations] (12 Cols)                                       │
└──────────────────────────────────────────────────────────────────────────┘
```

### Widget Details

1. **W-1: Health Checklist (8 columns)**: List of onboarding tasks (e.g., configure AI provider, define dev environment, invite team). Includes status tags.
2. **W-2: Active Executors (4 columns)**: Dynamic list of connected execution environments (Local, WSL, SSH, Kubernetes) with active task counts and latency ratings.
3. **W-3: Resource Telemetry (8 columns)**: Micro charts plotting overall CPU, memory, and disk consumption across active workspaces.
4. **W-4: Running Jobs Queue (4 columns)**: Real-time list of executing tasks with progress bars and quick-cancel buttons.
5. **W-5: Workspace Quick Link (6 columns)**: List of pinned and recently opened development workspaces, showing branch and last active timestamp.
6. **W-6: Activity Timeline (6 columns)**: Unified audit and task log events feed, detailing changes made across the organization.
7. **W-7: Recommendations (12 columns)**: Dynamic list of tasks (e.g., "Archive idle workspace to save disk", "Update Ollama endpoint to Ollama v2").

---

## 2. Workspace Layout & UX Specification (Gap #9)

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

## 3. Performance Targets (Gap #15)

To ensure high-grade enterprise responsiveness, all UI mutations and navigation renders must hit these measurable performance metrics:

| Page / Flow | Target Load Metric | API Latency Target (95th Percentile) | Max Permissible Load Time (Hard Limit) |
|---|---|---|---|
| **Dashboard** | Full render with metrics | `< 150ms` | `< 1.0 second` |
| **Workspace Shell** | Initial interface load | `< 200ms` | `< 2.0 seconds` |
| **Workspace File Open**| Load content to editor | `< 50ms` | `< 300ms` |
| **Global Search** | Retrieve dropdown matches | `< 100ms` | `< 300ms` |
| **AI Provider Test** | Connection handshake ping | `< 2.0 seconds` | `< 5.0 seconds` |
| **Deployments Board** | Render release grid | `< 120ms` | `< 2.0 seconds` |

---

## 4. Production Acceptance Checklist (Gap #16)

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
