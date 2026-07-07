# Wireframes — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## Directory Structure

```
docs/design/WIREFRAMES/
├── README.md              (this file)
├── DESKTOP/
│   ├── 01-bootstrap.md
│   ├── 02-login.md
│   ├── 03-setup-wizard.md
│   ├── 04-dashboard.md
│   ├── 05-settings.md
│   ├── 06-members.md
│   ├── 07-providers.md
│   ├── 08-environments.md
│   ├── 09-projects.md
│   ├── 10-workspace.md
│   └── 11-devcenter.md
├── TABLET/
│   ├── (same page set — collapsed sidebar, stacked grid)
│   └── ...
└── MOBILE/
    ├── (same page set — single column, bottom nav)
    └── ...
```

---

## Wireframe Specification Per Page

Each wireframe document must contain:

1. **Page Title & Route**
2. **ASCII Layout Diagram** (structural wireframe)
3. **Component Mapping** (which components from COMPONENT_LIBRARY.md)
4. **Responsive Notes** (what changes at each breakpoint)
5. **Interaction Notes** (hover, click, keyboard)
6. **Pixel Tolerance** (±2px from approved mockup)

---

## Desktop Wireframes

### 01 — Bootstrap Loading (`/`)

```
┌─────────────────────────────────────────────────┐
│                                                 │
│           ┌───────────────────────┐             │
│           │  [Brand Mark]  D      │             │
│           │  DevServer             │             │
│           │                       │             │
│           │  Initializing...      │             │
│           │                       │             │
│           │  ████ ████ ████ ████  │             │
│           │  (pulse animation)    │             │
│           └───────────────────────┘             │
│                                                 │
└─────────────────────────────────────────────────┘
```

Components: `bootstrap-screen.tsx`, brand-lockup

### 02 — Login (`/login`)

```
┌─────────────────────────────────────────────────┐
│                                                 │
│           ┌───────────────────────┐             │
│           │  [Brand Mark]         │             │
│           │  Welcome Back         │             │
│           │                       │             │
│           │  ┌─────────────────┐  │             │
│           │  │ Email           │  │             │
│           │  └─────────────────┘  │             │
│           │  ┌─────────────────┐  │             │
│           │  │ Password        │  │             │
│           │  └─────────────────┘  │             │
│           │                       │             │
│           │  [    Sign In     ]   │             │
│           │                       │             │
│           │  ┌─ Knowledge ──────┐ │             │
│           │  │ Learn Card       │ │             │
│           │  │ Learn Card       │ │             │
│           │  └──────────────────┘ │             │
│           └───────────────────────┘             │
└─────────────────────────────────────────────────┘
```

Components: `Input`, `Button(primary)`, `LearnCard`

### 04 — Dashboard (`/dashboard`)

```
┌──────────┬──────────────────────────────────────┐
│ SIDEBAR  │ TOPBAR  [Search...]  [Org▾] [Badge]  │
│          ├──────────────────────────────────────┤
│ [Brand]  │ SECTION HEADER                       │
│          │ Dashboard / System Overview          │
│ Overview │                                      │
│ ● Dashbd │ ┌────────┐┌────────┐┌────────┐┌────┐│
│   Worksp │ │Metric 1││Metric 2││Metric 3││M 4 ││
│   Projct │ └────────┘└────────┘└────────┘└────┘│
│   Settin │                                      │
│          │ ┌─────────────────┐┌────────────────┐│
│ Operatns │ │  Health Card    ││ Activity Feed  ││
│   Deploy │ │  Checklist      ││ Timeline items ││
│   Server │ │                 ││                ││
│   Domain │ └─────────────────┘└────────────────┘│
│   SSL    │                                      │
│   Servcs │ ┌─────────────────┐┌────────────────┐│
│          │ │  Tasks Queue    ││ Recommendations││
│ Platform │ └─────────────────┘└────────────────┘│
│   Databs │                                      │
│   Storag │ ┌────────────────────────────────────┐│
│   Termnl │ │  Knowledge Cards                  ││
│   Monitr │ └────────────────────────────────────┘│
│   Logs   │                                      │
│   Users  │                                      │
└──────────┴──────────────────────────────────────┘
```

Components: `AppShell`, `SectionHeader`, `MetricCard`, `Card`, `Badge`

### 10 — Workspace (`/workspace/[id]`)

```
┌──────────┬───────────────┬──────────────────────┐
│ SIDEBAR  │ SECTION TABS  │ CONTENT AREA         │
│          │               │                      │
│ (same    │ ▸ Overview    │ (varies by section)  │
│  as      │ ▸ Files       │                      │
│  dashbd) │ ▸ Repository  │ Overview:            │
│          │ ▸ Environment │  Health Score + Meta  │
│          │ ▸ Infrastr.   │                      │
│          │ ▸ Services    │ Files:               │
│          │ ▸ Tasks       │  Tree + Editor split  │
│          │ ▸ Deployments │                      │
│          │ ▸ Knowledge   │ Tasks:               │
│          │ ▸ ...more     │  Queue + Output pane  │
│          │               │                      │
└──────────┴───────────────┴──────────────────────┘
```

Components: `AppShell`, `Tabs`, workspace section components

---

## Tablet Adaptations (721px–960px)

- Sidebar collapses to top → horizontal scroll or hamburger
- Metric grid: 2 → 1 columns
- Page grid: 2 → 1 columns
- Workspace: section tabs move to horizontal scrollable bar

## Mobile Adaptations (≤720px)

- Sidebar hidden → bottom navigation bar or hamburger menu
- All grids: 1 column
- Content padding reduced to 18px
- Cards stack vertically
- Module rows, config rows → column layout
- Tables → card-based list view

---

## Pixel Tolerance Rule

> Every implemented screen must match the approved mockup within **±2px** tolerance unless a justified deviation is documented and approved.

---

> [!IMPORTANT]
> This wireframe set must be approved before any WP-8.6.17 UI implementation begins.
