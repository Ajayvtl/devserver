# Component Library — DevServer Platform

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## Existing Components (18 UI primitives)

### C-001: Button (`components/ui/button.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Primary interactive trigger for actions |
| **Variants** | `primary`, `secondary`, `ghost` |
| **States** | Default, Hover (translateY -1px), Disabled (opacity 0.55), Loading |
| **Props** | `variant`, `onClick`, `disabled`, `children`, `type`, `className` |
| **Accessibility** | Must have visible label or `aria-label`; supports `:focus-visible` |

### C-002: Badge (`components/ui/badge.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Visual label for status, count, or category |
| **Variants** | `neutral`, `accent`, `success`, `warning`, `danger`, `info` |
| **States** | Static only |
| **Props** | `tone: Tone`, `children` |
| **Accessibility** | Decorative only; parent must convey meaning |

### C-003: Card (`components/ui/card.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Content container with header and body |
| **Variants** | Default (single variant) |
| **States** | Default, Loading (skeleton), Empty, Error |
| **Props** | `eyebrow`, `title`, `badge`, `children` |
| **Accessibility** | Use `<section>` with `aria-labelledby` on title |

### C-004: Input (`components/ui/input.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Text input field for forms and search |
| **Variants** | Default (single variant) |
| **States** | Default, Focus, Error, Disabled |
| **Props** | Standard `<input>` props + `className` |
| **Accessibility** | Must have associated `<label>` or `aria-label` |

### C-005: Textarea (`components/ui/textarea.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Multi-line text input |
| **Variants** | Default |
| **States** | Default, Focus, Error, Disabled |
| **Props** | Standard `<textarea>` props + `className` |
| **Accessibility** | Must have associated `<label>` or `aria-label` |

### C-006: Dialog (`components/ui/dialog.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Modal dialog for confirmations and forms |
| **Variants** | Default, Destructive (red accent) |
| **States** | Open, Closed |
| **Props** | `isOpen`, `onClose`, `title`, `children` |
| **Accessibility** | `role="dialog"`, `aria-modal="true"`, focus trap, Escape to close |

### C-007: Table (`components/ui/table.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Data table for list displays |
| **Variants** | Default |
| **States** | Default, Loading (skeleton rows), Empty |
| **Props** | `headers: string[]`, `children` (rows) |
| **Accessibility** | Proper `<thead>/<tbody>`, `scope="col"` on headers |

### C-008: Tabs (`components/ui/tabs.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Tab navigation for section switching |
| **Variants** | Default |
| **States** | Active, Inactive, Disabled |
| **Props** | `tabs: {key, label}[]`, `active`, `onChange` |
| **Accessibility** | `role="tablist"`, `role="tab"`, `aria-selected`, keyboard arrow keys |

### C-009: Empty State (`components/ui/empty-state.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Placeholder for zero-data screens |
| **Variants** | Default |
| **States** | Static |
| **Props** | `title`, `description`, `action?` (CTA button) |
| **Accessibility** | Informational; use `role="status"` |

### C-010: Error State (`components/ui/error-state.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Error feedback for failed data loads |
| **Variants** | Default |
| **States** | Static |
| **Props** | `title`, `description`, `onRetry?` |
| **Accessibility** | `role="alert"`, `aria-live="assertive"` |

### C-011: Skeleton (`components/ui/skeleton.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Loading placeholder with shimmer animation |
| **Variants** | Default |
| **States** | Animating |
| **Props** | `className` (height/width via CSS) |
| **Accessibility** | `aria-busy="true"`, `aria-label="Loading"` |

### C-012: Progress (`components/ui/progress.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Progress bar for task/operation tracking |
| **Variants** | Default |
| **States** | Indeterminate, Determinate |
| **Props** | `value: number` (0-100) |
| **Accessibility** | `role="progressbar"`, `aria-valuenow`, `aria-valuemin`, `aria-valuemax` |

### C-013: Section Header (`components/ui/section-header.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Page/section title with optional action buttons |
| **Variants** | Default |
| **States** | Static |
| **Props** | `eyebrow`, `title`, `description`, `actions` |
| **Accessibility** | Title rendered as `<h1>` or `<h2>` based on hierarchy |

### C-014: Metric Card (`components/ui/metric-card.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | KPI display with label, value, and trend |
| **Variants** | Default |
| **States** | Default, Loading (skeleton) |
| **Props** | `label`, `value`, `delta`, `helper`, `tone` |
| **Accessibility** | `aria-label` with full context string |

### C-015: Code Block (`components/ui/code-block.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Formatted code display |
| **Variants** | Default |
| **States** | Static |
| **Props** | `title`, `children` (code string) |
| **Accessibility** | `<pre>` with `<code>`, avoid `aria-hidden` |

### C-016: Stepper (`components/ui/stepper.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Multi-step wizard progress indicator |
| **Variants** | Default |
| **States** | Complete, Active, Pending |
| **Props** | `steps: WizardStep[]`, `current: number` |
| **Accessibility** | `aria-current="step"`, step count announcement |

### C-017: Learn Card (`components/ui/learn-card.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Knowledge article preview card |
| **Variants** | Default |
| **States** | Static |
| **Props** | `title`, `question`, `summary`, `link` |
| **Accessibility** | Card as `<article>`, link as primary action |

### C-018: Provider Card (`components/ui/provider-card.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | AI provider configuration and management card |
| **Variants** | Default |
| **States** | Default, Testing, Connected, Error |
| **Props** | Full provider config object |
| **Accessibility** | Expandable region with `aria-expanded` |

---

## Composite Components (from feature modules)

### C-019: App Shell (`components/app-shell.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Top-level layout with sidebar nav, topbar, and content area |
| **Props** | `children`, `server`, `notifications` |

### C-020: Members Panel (`components/settings/members-panel.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Full member management interface with CRUD operations |
| **Props** | Internal — uses `useMembersManagement` hook |

### C-021: Organizations Panel (`components/settings/orgs-panel.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Organization listing and creation |
| **Props** | Internal — uses auth context |

### C-022: Environments Panel (`components/config/environments-panel.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Environment/variable/secret management |
| **Props** | Internal — uses API client |

### C-023: Providers Panel (`components/config/providers-panel.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | AI provider listing, creation, and connection testing |
| **Props** | Internal — uses API client |

### C-024: Toast (`components/toast.tsx`)

| Property | Detail |
|---|---|
| **Purpose** | Notification popup for async feedback |
| **Variants** | Success, Error, Info, Warning |
| **Props** | `message`, `type`, `duration` |
| **Accessibility** | `role="alert"`, `aria-live="polite"` |

---

## Components Needed (not yet built)

| ID | Component | Purpose | Priority |
|---|---|---|---|
| C-025 | Dropdown Menu | Context menus, action menus | High |
| C-026 | Select | Styled select replacement | High |
| C-027 | Switch/Toggle | Boolean settings | Medium |
| C-028 | Tooltip | Contextual help text | Medium |
| C-029 | Avatar | User identity display | Medium |
| C-030 | Breadcrumb | Navigation hierarchy | Low |
| C-031 | Pagination | Table/list page navigation | High |
| C-032 | Date Picker | Date range selection | Low |
| C-033 | Command Palette | Quick navigation (Ctrl+K) | Medium |

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins.
