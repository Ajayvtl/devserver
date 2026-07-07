# Component Library — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## Component Categorization Schema

To prevent code duplication, every component is classified into one of these strict architectural tiers:
1. **Atomic**: Reusable, pure visual primitives with no direct dependencies or API calls.
2. **Composite**: Combines multiple atomic components. May contain internal state but no direct API hooks.
3. **Layout**: Structures pages, panels, grids, and shell containers.
4. **Shared**: Utility components used across multiple pages (e.g., toast alerts, modals).
5. **Page-only**: Specialized components tied directly to a single route, often executing API calls.

---

## 1. Core Component Specifications

### C-001: Button (`components/ui/button.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Atomic** (Shared Primitive) |
| **Purpose** | Primary interactive trigger for actions |
| **Variants** | `primary` (gradient accent), `secondary` (border), `ghost` (text-only), `destructive` (danger red) |
| **Props** | `variant: ButtonVariant`, `onClick: () => void`, `disabled?: boolean`, `loading?: boolean`, `children` |
| **States** | Default, Hover (translateY -1px), Focus (outline ring), Active (translateY 0px), Loading (shows spinner, disables clicks), Disabled (opacity 0.55, pointer-events none) |
| **Responsive Rules** | **Desktop**: Fixed width or inline padding. **Tablet**: Same. **Mobile**: Expand to `width: 100%` inside modal footers or forms to allow easy thumb-tap targets. |
| **Accessibility** | Focus ring on tab navigation; screen reader announces loading state as `aria-busy="true"`. |

### C-002: Badge (`components/ui/badge.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Atomic** (Shared Primitive) |
| **Purpose** | Visual tag for status, counts, or entity categories |
| **Variants** | `neutral`, `accent`, `success` (green), `warning` (yellow), `danger` (red), `info` (blue) |
| **Props** | `tone: BadgeTone`, `size?: 'sm' \| 'md'`, `children` |
| **States** | Static |
| **Responsive Rules** | Font-size shrinks by 1px on mobile. Remains inline and does not wrap. |
| **Accessibility** | Purely presentation; parent container must provide textual explanation if status is semantic. |

### C-003: Card (`components/ui/card.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Layout** (Container) |
| **Purpose** | Content grouping box with header, body, and actions |
| **Variants** | `default` (bordered), `elevated` (subtle shadow), `interactive` (hover border highlight) |
| **Props** | `eyebrow?: string`, `title?: string`, `badge?: ReactNode`, `actions?: ReactNode`, `children` |
| **States** | Default, Loading (renders C-011 skeleton inside), Empty (renders C-009 empty-state), Error (renders C-010 error-state) |
| **Responsive Rules** | **Desktop/Tablet**: Padding is 24px (`--space-10`). **Mobile**: Padding drops to 16px (`--space-7`) to maximize viewable screen estate. Interactivity border-glow is disabled on touch screens. |
| **Accessibility** | Semantically rendered as `<section>` or `<article>`. Header uses `<h2>` or `<h3>` dynamically. |

### C-004: Table (`components/ui/table.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Composite** (Shared) |
| **Purpose** | Tabular representation of structures (members, variables, logs) |
| **Props** | `headers: string[]`, `children: ReactNode` (rows) |
| **States** | Default, Loading (renders skeleton rows), Empty |
| **Responsive Rules** | **Desktop/Tablet**: Standard grid layout with horizontal rows. **Mobile**: Column-headers are hidden; each row transforms into an individual card-like panel with labels stack-aligned to the left. |
| **Accessibility** | Must include `role="table"`, `<thead scope="col">` on headers, and dynamic cell labeling for screen readers. |

### C-005: Dialog (`components/ui/dialog.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Shared** (Overlay) |
| **Purpose** | Modal dialog blocking screen actions for confirmation or wizards |
| **Props** | `isOpen: boolean`, `onClose: () => void`, `title: string`, `children` |
| **States** | Open (fadeIn & scale-up), Closed (hidden) |
| **Responsive Rules** | **Desktop**: Width is centered at `max-width: 600px`. **Tablet**: Centers at `max-width: 500px`. **Mobile**: Slide up from screen bottom to cover 100% width and 90% height, maximizing mobile keyboard screen room. |
| **Accessibility** | Implements `focus-trap` inside modal; `aria-modal="true"`; closes instantly on pressing `Escape`. |

### C-006: Tabs (`components/ui/tabs.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Composite** (Navigation) |
| **Purpose** | Local page route or category switcher |
| **Props** | `tabs: {key: string, label: string}[]`, `activeKey: string`, `onChange: (key: string) => void` |
| **States** | Active (accent border), Inactive, Disabled |
| **Responsive Rules** | **Desktop/Tablet**: Horizontal tab row. **Mobile**: Transforms into a touch-friendly horizontal swipe-carousal with overflow indicators, or switches to a dropdown menu. |
| **Accessibility** | Implements `role="tablist"`, `role="tab"`, and `aria-selected` status. Supports arrow key keyboard navigation. |

### C-007: Empty State (`components/ui/empty-state.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Composite** (Shared Visual) |
| **Purpose** | Onboarding card for sections with zero data |
| **Props** | `title: string`, `description: string`, `illustration?: string`, `actionButton?: ReactNode` |
| **States** | Static |
| **Responsive Rules** | Vertical stack alignment. Text-align is centered. Shrinks padding on mobile. |
| **Accessibility** | Utilizes `role="status"` to announce the empty screen status. |

---

## 2. Enterprise & System Components (New)

### C-008: Notification Inbox (`components/enterprise/notification-inbox.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Page-only** / **Shared Popover** |
| **Purpose** | Shows alerts, mentions, failed background tasks, and health warnings |
| **Props** | `unreadCount: number`, `notifications: Notification[]`, `onMarkAllAsRead: () => void` |
| **States** | Idle, Opening, MarkAsRead (fades out active indicator), Empty, Error |
| **Responsive Rules** | **Desktop**: Drops down from topbar. **Mobile**: Opens as full-screen modal with bottom navigation back to settings. |
| **Accessibility** | Screen reader announces new notification count dynamically with `aria-live="polite"`. |

### C-009: Command Palette (`components/enterprise/command-palette.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Shared** (Search System) |
| **Purpose** | Quick navigation (Ctrl+K) for workspaces, services, files, search commands, and settings |
| **Props** | `isOpen: boolean`, `onClose: () => void` |
| **States** | Open (scale-in overlay), Searching, ResultsFound, EmptyState |
| **Responsive Rules** | **Desktop/Tablet**: Centered modal overlay (`max-width: 650px`). **Mobile**: Overlaid full-screen menu with large search text input and instant touch-targets. |
| **Accessibility** | Full arrow keys controls; input autofocuses instantly; `role="combobox"`. |

### C-010: CPU/RAM/Network Timeline (`components/monitoring/metrics-timeline.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Composite** (Charts) |
| **Purpose** | Real-time system resource observer graphs |
| **Props** | `metricType: 'cpu' \| 'ram' \| 'disk' \| 'network'`, `dataPoints: MetricDataPoint[]` |
| **States** | Loading (shimmer charts), Connected (live updating), Degraded (orange glow), Alerting (red pulse) |
| **Responsive Rules** | **Desktop/Tablet**: Side-by-side grid panels. **Mobile**: Vertical stack layout; graph detail density is automatically downsampled for readability on narrow viewports. |
| **Accessibility** | Includes fallback tabular summary view for visual accessibility. |

### C-011: Timeline Feed (`components/settings/timeline-feed.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Composite** (Observable) |
| **Purpose** | Activity timeline tracking personal mutations, workspace tasks, and org logs |
| **Props** | `activities: ActivityItem[]`, `scope: 'personal' \| 'workspace' \| 'organization'` |
| **States** | Loading, Active, Filtered, Empty |
| **Responsive Rules** | **Desktop/Tablet**: Vertical left-aligned line with offset descriptions and dates. **Mobile**: Dates stack above content; inline icons collapse to save width. |
| **Accessibility** | Semantically structured as an ordered list `<ol>` with step status. |

### C-012: Backup & Restore Wizard (`components/operations/backup-wizard.tsx`)

| Property | Detail |
|---|---|
| **Tier** | **Composite** (Wizard) |
| **Purpose** | Direct UI flow for backing up, restoring database, exporting configuration, or importing states |
| **Props** | `type: 'backup' \| 'restore' \| 'export' \| 'import'`, `onTrigger: (config: any) => Promise<void>` |
| **States** | SelectTargets, RunningProgress (shows C-012 progress), CompleteSuccess, AlertFailed |
| **Responsive Rules** | Fits into standard dialog grid. Shrinks grid padding on mobile. |
| **Accessibility** | Stepper progress indicator has dynamic announcements for active step and tasks success. |

---

## Component Tier and Ownership Matrix

To keep dependencies clean, files must follow these import rules:
* **Atomic** components may NOT import other local components.
* **Composite** components may only import **Atomic** and utility functions.
* **Layout** and **Shared** overlays can import both **Atomic** and **Composite** components.
* **Page-only** components can import any category, but cannot be imported by other modules (except routing pages).

```
┌───────────────────────────────────────────────┐
│                   PAGE-ONLY                   │
└───────────────────────┬───────────────────────┘
                        │ imports
┌───────────────────────▼───────────────────────┐
│              LAYOUT / SHARED                  │
└───────────────────────┬───────────────────────┘
                        │ imports
┌───────────────────────▼───────────────────────┐
│                  COMPOSITE                    │
└───────────────────────┬───────────────────────┘
                        │ imports
┌───────────────────────▼───────────────────────┐
│                    ATOMIC                     │
└───────────────────────────────────────────────┘
```

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
