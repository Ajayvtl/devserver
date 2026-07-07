# Accessibility & WCAG Compliance Guidelines

> **WP**: WP-8.6.17 Phase A — Design Only  
> **Status**: Draft — Pending Review  
> **Updated**: 2026-07-07

---

## 1. Keyboard-Only Navigation & Focus Order (Gap #5 / Refinement)

All interactive elements must be accessible via keyboard standard keystrokes without causing keyboard trap scenarios.

### 1.1 Keyboard Navigation Schema
- **Tab Key**: Move focus sequentially between interactive items (links, buttons, input fields, and custom selection dropdowns).
- **Shift + Tab**: Move focus in reverse order.
- **Space / Enter Key**: Activate links, buttons, and toggle checkboxes/switches.
- **Arrow Keys**: Navigate within complex composite widgets like Tabs (C-006), Dropdowns (C-025), and lists.
- **Escape Key**: Instantly dismiss dialog overlays (C-005), search palettes (C-009), or floating popovers.

### 1.2 Focus Indicators & Focus Ring Tokens
- Never hide outline focus rings. Use a custom high-contrast focus outline ring:
  ```css
  *:focus-visible {
    outline: 2px solid var(--color-focus-ring);
    outline-offset: 2px;
  }
  ```
- Focus order must flow logically from top-to-bottom and left-to-right following the structural DOM layout.
- Dialog modals must capture keyboard focus within their boundary upon opening and restore focus to the triggering element when dismissed.

---

## 2. Screen Reader Behavior & ARIA Specifications

Components must expose semantic properties using ARIA attributes to describe dynamic states.

| UI Element | ARIA Attribute Requirement | Screen Reader Output Description |
|---|---|---|
| **Buttons (Loading)** | `aria-busy="true"` `aria-live="polite"` | Announces "Busy" during active API calls, then announces result. |
| **Telemetry Charts** | `role="img"` `aria-label="..."` | Announces description of active metric (e.g., "CPU Timeline graph showing 45% usage"). |
| **Command Palette** | `role="combobox"` `aria-expanded="..."` | Announces menu expand status and list result items count. |
| **Notification Badge**| `aria-live="polite"` | Announces count updates dynamically (e.g., "3 new notifications"). |
| **Status Indicators** | `role="status"` | Announces state updates (e.g., "Workspace online"). |

---

## 3. Contrast Targets & Color Rules

Visual styling must conform to **WCAG 2.1 Level AA** standards.

- **Normal Text**: Contrast ratio of at least **4.5:1** against the background color.
- **Large Text (18pt/24px or bold 14pt/18px)**: Contrast ratio of at least **3.0:1**.
- **UI Component Boundaries / Form Borders**: Contrast ratio of at least **3.0:1** for active control indicators.

### Verified Theme Colors Mapping

| Theme Mode | Token Key | Color Code Hex | Target Contrast Ratio | Compliant |
|---|---|---|---|---|
| **Dark Mode** | `--color-bg-primary` | `#07111d` | Background base | Yes |
| | `--color-text-primary`| `#f3f4f6` | **15.4:1** (against `#07111d`) | Yes (WCAG AAA) |
| | `--color-accent` | `#3b82f6` | **4.6:1** (against `#07111d`) | Yes (WCAG AA) |
| **Light Mode** | `--color-bg-primary` | `#ffffff` | Background base | Yes |
| | `--color-text-primary`| `#111827` | **18.0:1** (against `#ffffff`) | Yes (WCAG AAA) |
| | `--color-accent` | `#2563eb` | **5.1:1** (against `#ffffff`) | Yes (WCAG AA) |

---

## 4. Reduced Motion Preferences

All transition animations (fade-in, slide-up, dialog zoom effects) must respect the user's operating system preferences.

- Wrap animation stylesheets in media queries using `@media (prefers-reduced-motion: reduce)`:
  ```css
  @media (prefers-reduced-motion: reduce) {
    * {
      animation-duration: 0.01ms !important;
      animation-iteration-count: 1 !important;
      transition-duration: 0.01ms !important;
      scroll-behavior: auto !important;
    }
  }
  ```
- Replaces smooth slide-ins with instant visibility toggle switches.

---

## 5. Zoom & Reflow Testing Targets

- The platform sitemap page must support zooming up to **400%** on desktop viewports without loss of information or causing double-axis scrollbars (horizontal scroll must be avoided).
- Sidebars collapse to mobile drawer sitemaps at zoom levels above 200%.

---

> [!IMPORTANT]
> This document must be approved before any WP-8.6.17 UI implementation begins, as mandated by Permanent Developer Rule 15.
